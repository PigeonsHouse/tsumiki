package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"tsumiki/env"
)

const (
	scope       = "identify email guilds"
	apiEndpoint = "https://discord.com/api/v10"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type UserInfoResponse struct {
	ID         string  `json:"id"`
	UserName   string  `json:"username"`
	GlobalName *string `json:"global_name"`
	Avatar     string  `json:"avatar"`
}

type GuildInfoResponse struct {
	ID string `json:"id"`
}

func callbackUrl() string {
	return fmt.Sprintf("%s/api/v1/auth/discord/callback", env.BackendUrl)
}

func getAvatarUrl(userInfo UserInfoResponse) string {
	return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", userInfo.ID, userInfo.Avatar)
}

func FetchAvatar(userInfo UserInfoResponse) (io.ReadCloser, string, error) {
	resp, err := http.Get(getAvatarUrl(userInfo))
	if err != nil {
		return nil, "", err
	}
	return resp.Body, resp.Header.Get("Content-Type"), nil
}

func GetRedirectUrl() string {
	return fmt.Sprintf("%s/oauth2/authorize?client_id=%s&redirect_uri=%s&scope=%s&response_type=code",
		apiEndpoint, env.DiscordClientID, url.QueryEscape(callbackUrl()), scope)
}

func ValidateRedirectedCode(code string) (TokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", env.DiscordClientID)
	data.Set("client_secret", env.DiscordSecretID)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", callbackUrl())

	req, err := http.NewRequest("POST", apiEndpoint+"/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return TokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TokenResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TokenResponse{}, fmt.Errorf("todo 500")
	}

	var tokenRes TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		return TokenResponse{}, err
	}
	return tokenRes, nil
}

func GetUserInfo(tokenRes TokenResponse) (UserInfoResponse, error) {
	userReq, _ := http.NewRequest("GET", apiEndpoint+"/users/@me", nil)
	userReq.Header.Set("Authorization", "Bearer "+tokenRes.AccessToken)

	client := &http.Client{}
	userResp, err := client.Do(userReq)
	if err != nil {
		return UserInfoResponse{}, err
	}
	defer userResp.Body.Close()

	var userInfo UserInfoResponse
	if err := json.NewDecoder(userResp.Body).Decode(&userInfo); err != nil {
		return UserInfoResponse{}, err
	}
	return userInfo, nil
}

func GetUserGuildsInfo(tokenRes TokenResponse) ([]GuildInfoResponse, error) {
	userReq, _ := http.NewRequest("GET", apiEndpoint+"/users/@me/guilds", nil)
	userReq.Header.Set("Authorization", "Bearer "+tokenRes.AccessToken)

	client := &http.Client{}
	userResp, err := client.Do(userReq)
	if err != nil {
		return []GuildInfoResponse{}, err
	}
	defer userResp.Body.Close()

	var guildsInfo []GuildInfoResponse
	if err := json.NewDecoder(userResp.Body).Decode(&guildsInfo); err != nil {
		return []GuildInfoResponse{}, err
	}
	return guildsInfo, nil
}
