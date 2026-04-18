package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"tsumiki/helper"
	"tsumiki/media"
	"tsumiki/middleware"
	"tsumiki/repository"

	"github.com/go-chi/chi/v5"
)

type UserHandler interface {
	GetMyInfo(w http.ResponseWriter, r *http.Request)
	GetUserInfo(w http.ResponseWriter, r *http.Request)
}

type userHandlerImpl struct {
	repository repository.UserRepository
	media      media.MediaService
}

func NewUserHandler(userRepo repository.UserRepository, mediaSvc media.MediaService) UserHandler {
	return &userHandlerImpl{
		repository: userRepo,
		media:      mediaSvc,
	}
}

func (uh *userHandlerImpl) GetMyInfo(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.ResponseUnauthorized(w, "認証情報が見つかりません")
		return
	}

	user, err := uh.repository.FindByID(userID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if user == nil {
		helper.ResponseNotFound(w, "ユーザが見つかりません")
		return
	}

	userResponse := *user
	userResponse.AvatarUrl = uh.media.ResolveURL(user.AvatarUrl)
	helper.ResponseOk(w, userResponse)
}

func (uh *userHandlerImpl) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		helper.ResponseBadRequest(w, "ユーザIDが不正です")
		return
	}

	user, err := uh.repository.FindByID(userID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if user == nil {
		helper.ResponseNotFound(w, "ユーザが見つかりません")
		return
	}

	userResponse := *user
	userResponse.AvatarUrl = uh.media.ResolveURL(user.AvatarUrl)
	helper.ResponseOk(w, userResponse)
}
