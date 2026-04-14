package handler

import (
	"fmt"
	"net/http"
	"tsumiki/env"
	"tsumiki/external"
	"tsumiki/helper"
)

func RedirectDiscord(w http.ResponseWriter, r *http.Request) {
	redirectUrl := external.GetRedirectUrl()
	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
}

func CallbackDiscord(w http.ResponseWriter, r *http.Request) {
	// エラーパラメーターのチェック（ユーザーがキャンセルした場合など）
	if errDesc := r.URL.Query().Get("error_description"); errDesc != "" {
		http.Redirect(w, r, env.FrontendUrl+"?error=access_denied", http.StatusTemporaryRedirect)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		helper.ResponseBadRequest(w, "認可コードが見つかりません")
		return
	}

	tokenRes, err := external.ValidateRedirectedCode(code)
	if err != nil {
		helper.ResponseBadRequest(w, "バリデーションに失敗しました")
		return
	}

	userInfo, err := external.GetUserInfo(tokenRes)
	if err != nil {
		helper.ResponseBadRequest(w, "ユーザ情報の解決に失敗しました")
		return
	}

	// ==============================================================
	// 【APIサーバーとしての重要ポイント】
	// ここで userInfo (DiscordのIDやユーザー名) を自社のデータベースと照合します。
	// 新規ユーザーならDBに登録、既存ならログイン処理を行います。
	//
	// その後、フロントエンド用の「独自のセッショントークン（JWTなど）」を生成します。
	// 今回はモックとして単純な文字列を生成したと仮定します。
	// ==============================================================
	fmt.Printf("取得したDiscordユーザー情報: %s\n", string(userInfo))

	myAppSessionToken := "generated_dummy_jwt_or_session_id"

	// HttpOnly Cookieにトークンをセットしてリダイレクトする
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    myAppSessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	redirectURL := env.FrontendUrl

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
