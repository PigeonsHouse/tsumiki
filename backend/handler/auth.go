package handler

import (
	"fmt"
	"net/http"
	"time"
	"tsumiki/env"
	"tsumiki/external"
	"tsumiki/helper"
	"tsumiki/media"
	"tsumiki/middleware"
	"tsumiki/repository"
	"tsumiki/store"
)

type AuthHandler interface {
	RedirectDiscord(w http.ResponseWriter, r *http.Request)
	CallbackDiscord(w http.ResponseWriter, r *http.Request)
}

type authHandlerImpl struct {
	repository repository.AuthRepository
	store      store.AuthStore
	media      media.MediaService
	discord    external.DiscordService
}

func NewAuthHandler(authRepo repository.AuthRepository, authStore store.AuthStore, mediaSvc media.MediaService, discordSvc external.DiscordService) AuthHandler {
	return &authHandlerImpl{
		repository: authRepo,
		store:      authStore,
		media:      mediaSvc,
		discord:    discordSvc,
	}
}

func (ah *authHandlerImpl) RedirectDiscord(w http.ResponseWriter, r *http.Request) {
	redirectUrl := ah.discord.GetRedirectUrl()
	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
}

func (ah *authHandlerImpl) CallbackDiscord(w http.ResponseWriter, r *http.Request) {
	// エラーパラメーターのチェック（ユーザーがキャンセルした場合など）
	if errDesc := r.URL.Query().Get("error_description"); errDesc != "" {
		fmt.Println("認証に失敗しました")
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
		fmt.Println("validate code: ", err)
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

	user, err := ah.repository.FindByDiscordUserId(userInfo.ID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if user == nil {
		user, err = ah.repository.CreateUserByDiscord(userInfo.UserName, "", userInfo.ID, guildID)
		if err != nil {
			fmt.Println("DBエラー: ", err)
			helper.ResponseInternalServerError(w, "DBエラー")
			return
		}
		avatarBody, avatarContentType, err := ah.discord.FetchAvatar(userInfo)
		if err != nil {
			fmt.Println("アバター取得エラー: ", err)
			helper.ResponseInternalServerError(w, "アバター取得エラー")
			return
		}
		defer avatarBody.Close()
		avatarPath, err := ah.media.UploadAvatar(r.Context(), user.ID, avatarBody, avatarContentType)
		if err != nil {
			fmt.Println("アバターアップロードエラー: ", err)
			helper.ResponseInternalServerError(w, "アバターアップロードエラー")
			return
		}
		if err := ah.repository.UpdateAvatarUrl(user.ID, avatarPath); err != nil {
			fmt.Println("DBエラー: ", err)
			helper.ResponseInternalServerError(w, "DBエラー")
			return
		}
		user.AvatarUrl = avatarPath
	}

	tokenPair, err := middleware.GenerateTokenPair(user.ID)
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
		MaxAge:   int(middleware.AccessTokenLiveTime / time.Second),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(middleware.RefreshTokenLiveTime / time.Second),
	})

	userResponse := *user
	userResponse.AvatarUrl = ah.media.ResolveURL(user.AvatarUrl)
	helper.ResponseOk(w, userResponse)
}
