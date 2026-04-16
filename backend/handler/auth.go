package handler

import (
	"net/http"
	"slices"
	"tsumiki/env"
	"tsumiki/external"
	"tsumiki/helper"
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
}

func NewAuthHandler(authRepo repository.AuthRepository, authStore store.AuthStore) AuthHandler {
	return &authHandlerImpl{
		repository: authRepo,
		store:      authStore,
	}
}

func (ah *authHandlerImpl) RedirectDiscord(w http.ResponseWriter, r *http.Request) {
	redirectUrl := external.GetRedirectUrl()
	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
}

func (ah *authHandlerImpl) CallbackDiscord(w http.ResponseWriter, r *http.Request) {
	// エラーパラメーターのチェック（ユーザーがキャンセルした場合など）
	if errDesc := r.URL.Query().Get("error_description"); errDesc != "" {
		helper.ResponseBadRequest(w, "認証に失敗しました")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		helper.ResponseBadRequest(w, "認可コードが見つかりません")
		return
	}

	tokenRes, err := external.ValidateRedirectedCode(code)
	if err != nil {
		helper.ResponseBadRequest(w, "認可コードのバリデーションに失敗しました")
		return
	}

	userInfo, err := external.GetUserInfo(tokenRes)
	if err != nil {
		helper.ResponseBadRequest(w, "ユーザ情報の解決に失敗しました")
		return
	}

	guildsInfo, err := external.GetUserGuildsInfo(tokenRes)
	if err != nil {
		helper.ResponseBadRequest(w, "ギルド情報の解決に失敗しました")
		return
	}

	guildID := ""
	for _, guildInfo := range guildsInfo {
		if slices.Contains(env.AllowGuildIds, guildInfo.ID) {
			guildID = guildInfo.ID
		}
	}
	if guildID == "" {
		helper.ResponseForbidden(w, "このDiscordユーザのログインは許容されていません")
		return
	}

	user, err := ah.repository.FindByDiscordUserId(userInfo.ID)
	if err != nil {
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if user == nil {
		user, err = ah.repository.CreateUserByDiscord(userInfo.UserName, userInfo.Avatar, userInfo.ID, guildID)
		if err != nil {
			helper.ResponseInternalServerError(w, "DBエラー")
			return
		}
	}

	// tokenPair, err := middleware.GenerateTokenPair(user.ID)
	_, err = middleware.GenerateTokenPair(user.ID)
	if err != nil {
		helper.ResponseInternalServerError(w, "トークン生成エラー")
		return
	}

	// TODO: redisとかもいい感じに
	// ah.store.SetRefreshToken(tokenPair)

}
