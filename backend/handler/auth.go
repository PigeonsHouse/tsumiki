package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"tsumiki/auth"
	"tsumiki/env"
	"tsumiki/external"
	"tsumiki/helper"
	"tsumiki/media"
	"tsumiki/middleware"
	"tsumiki/repository"
	"tsumiki/store"
)

type OAuthState struct {
	Mode string `json:"mode"`
}

func encodeOAuthState(state OAuthState) (string, error) {
	b, err := json.Marshal(state)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func decodeOAuthState(encoded string) (OAuthState, error) {
	b, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return OAuthState{}, err
	}
	var state OAuthState
	if err := json.Unmarshal(b, &state); err != nil {
		return OAuthState{}, err
	}
	return state, nil
}

type AuthHandler interface {
	RedirectDiscord(w http.ResponseWriter, r *http.Request)
	CallbackDiscord(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type authHandlerImpl struct {
	repositories *repository.Repositories
	store        store.AuthStore
	media        media.MediaService
	discord      external.DiscordService
}

func NewAuthHandler(repos *repository.Repositories, authStore store.AuthStore, mediaSvc media.MediaService, discordSvc external.DiscordService) AuthHandler {
	return &authHandlerImpl{
		repositories: repos,
		store:        authStore,
		media:        mediaSvc,
		discord:      discordSvc,
	}
}

func (ah *authHandlerImpl) RedirectDiscord(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "json"
	}
	encodedState, err := encodeOAuthState(OAuthState{Mode: mode})
	if err != nil {
		helper.ResponseInternalServerError(w, "stateの生成に失敗しました")
		return
	}
	redirectUrl := ah.discord.GetRedirectUrl(encodedState)
	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
}

func (ah *authHandlerImpl) CallbackDiscord(w http.ResponseWriter, r *http.Request) {
	oauthState, err := decodeOAuthState(r.URL.Query().Get("state"))
	if err != nil {
		helper.ResponseBadRequest(w, "stateが無効です")
		return
	}

	// エラーパラメーターのチェック（ユーザーがキャンセルした場合など）
	if errDesc := r.URL.Query().Get("error_description"); errDesc != "" {
		fmt.Println("認証に失敗しました")
		if oauthState.Mode == "redirect" {
			http.Redirect(w, r, env.FrontendUrl+"/login/callback?error=auth_failed", http.StatusFound)
			return
		}
		helper.ResponseBadRequest(w, "認証に失敗しました")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		fmt.Println("認可コードが見つかりません")
		helper.ResponseBadRequest(w, "認可コードが見つかりません")
		return
	}

	tokenRes, err := ah.discord.ValidateRedirectedCode(code)
	if err != nil {
		fmt.Println("認可コードのバリデーションに失敗しました:", err)
		helper.ResponseBadRequest(w, "認可コードのバリデーションに失敗しました")
		return
	}

	userInfo, err := ah.discord.GetUserInfo(tokenRes)
	if err != nil {
		fmt.Println("get discord user: ", err)
		helper.ResponseBadRequest(w, "ユーザ情報の解決に失敗しました")
		return
	}

	guildsInfo, err := ah.discord.GetUserGuildsInfo(tokenRes)
	if err != nil {
		fmt.Println("get guild: ", err)
		helper.ResponseBadRequest(w, "ギルド情報の解決に失敗しました")
		return
	}

	guildID := ""
	for _, allowedID := range env.AllowGuildIds {
		for _, guildInfo := range guildsInfo {
			if guildInfo.ID == allowedID {
				guildID = allowedID
				break
			}
		}
		if guildID != "" {
			break
		}
	}
	if guildID == "" {
		fmt.Println("このDiscordユーザのログインは許容されていません")
		helper.ResponseForbidden(w, "このDiscordユーザのログインは許容されていません")
		return
	}

	user, err := ah.repositories.User.FindByDiscordUserId(userInfo.ID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if user == nil {
		err := ah.repositories.RunInTx(func(txRepos *repository.Repositories) error {
			user, err = txRepos.User.CreateUserByDiscord(userInfo.UserName, "", userInfo.ID, guildID)
			if err != nil {
				return err
			}
			avatarBody, avatarContentType, err := ah.discord.FetchAvatar(userInfo)
			if err != nil {
				return err
			}
			defer avatarBody.Close()
			avatarPath, err := ah.media.UploadAvatar(r.Context(), user.ID, avatarBody, avatarContentType)
			if err != nil {
				return err
			}
			if err := txRepos.User.UpdateAvatarUrl(user.ID, avatarPath); err != nil {
				return err
			}
			user.AvatarUrl = avatarPath
			return nil
		})
		if err != nil {
			fmt.Println("DBエラー:", err)
			helper.ResponseInternalServerError(w, "ユーザの新規登録時にエラーが発生しました")
			return
		}
	}

	tokenPair, err := auth.GenerateTokenPair(user.ID)
	if err != nil {
		fmt.Println("トークン生成エラー: ", err)
		helper.ResponseInternalServerError(w, "トークン生成エラー")
		return
	}

	if err := ah.store.SetRefreshToken(r.Context(), tokenPair.UserID, tokenPair.SessionID); err != nil {
		fmt.Println("セッション保存エラー: ", err)
		helper.ResponseInternalServerError(w, "セッション保存エラー")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenPair.AccessToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(auth.AccessTokenLiveTime / time.Second),
		Secure:   true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(auth.RefreshTokenLiveTime / time.Second),
		Secure:   true,
	})

	if oauthState.Mode == "redirect" {
		http.Redirect(w, r, env.FrontendUrl+"/login/callback", http.StatusFound)
		return
	}

	userResponse := *user
	userResponse.AvatarUrl = ah.media.ResolveURL(user.AvatarUrl)
	helper.ResponseOk(w, userResponse)
}

func (ah *authHandlerImpl) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.ResponseUnauthorized(w, "認証情報が見つかりません")
		return
	}

	if err := ah.store.DeleteAllRefreshTokens(r.Context(), userID); err != nil {
		fmt.Println("セッション削除エラー: ", err)
		helper.ResponseInternalServerError(w, "セッション削除エラー")
		return
	}

	// Cookieを無効化
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Secure:   true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Secure:   true,
	})

	helper.ResponseOk(w, nil)
}

func (ah *authHandlerImpl) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		helper.ResponseUnauthorized(w, "リフレッシュトークンが見つかりません")
		return
	}

	claims, err := auth.ValidateRefreshToken(cookie.Value)
	if err != nil {
		helper.ResponseUnauthorized(w, "リフレッシュトークンが無効です")
		return
	}

	valid, err := ah.store.ValidateAndDeleteRefreshToken(r.Context(), claims.UserID, claims.SessionID)
	if err != nil {
		fmt.Println("セッション確認エラー: ", err)
		helper.ResponseInternalServerError(w, "セッション確認エラー")
		return
	}
	if !valid {
		helper.ResponseUnauthorized(w, "セッションが無効です")
		return
	}

	tokenPair, err := auth.GenerateTokenPair(claims.UserID)
	if err != nil {
		fmt.Println("トークン生成エラー: ", err)
		helper.ResponseInternalServerError(w, "トークン生成エラー")
		return
	}

	if err := ah.store.SetRefreshToken(r.Context(), tokenPair.UserID, tokenPair.SessionID); err != nil {
		fmt.Println("セッション保存エラー: ", err)
		helper.ResponseInternalServerError(w, "セッション保存エラー")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenPair.AccessToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(auth.AccessTokenLiveTime / time.Second),
		Secure:   true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(auth.RefreshTokenLiveTime / time.Second),
		Secure:   true,
	})

	helper.ResponseOk(w, nil)
}
