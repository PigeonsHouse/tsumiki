package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tsumiki/helper"
	"tsumiki/middleware"
	"tsumiki/repository"

	"github.com/go-chi/chi/v5"
)

type WorkHandler interface {
	GetWorks(w http.ResponseWriter, r *http.Request)
	GetSpecifiedWork(w http.ResponseWriter, r *http.Request)
	GetWorkTsumiki(w http.ResponseWriter, r *http.Request)
	CreateWork(w http.ResponseWriter, r *http.Request)
	EditWork(w http.ResponseWriter, r *http.Request)
	DeleteWork(w http.ResponseWriter, r *http.Request)
}

type workHandlerImpl struct {
	repositories *repository.Repositories
}

func NewWorkHandler(repos *repository.Repositories) WorkHandler {
	return &workHandlerImpl{
		repositories: repos,
	}
}

type workRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (wh *workHandlerImpl) GetWorks(w http.ResponseWriter, r *http.Request) {
	pageSize, page, _ := parsePaginationQuery(r)

	works, err := wh.repositories.Work.GetWorks(pageSize, page)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	helper.ResponseOk(w, works)
}

func (wh *workHandlerImpl) GetSpecifiedWork(w http.ResponseWriter, r *http.Request) {
	workID, err := parseWorkID(r)
	if err != nil {
		helper.ResponseBadRequest(w, "作品IDが不正です")
		return
	}

	work, err := wh.repositories.Work.GetWork(workID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if work == nil {
		helper.ResponseNotFound(w, "作品が見つかりません")
		return
	}

	helper.ResponseOk(w, work)
}

func (wh *workHandlerImpl) GetWorkTsumiki(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetOptionalUserIDFromContext(r.Context())
	workID, err := parseWorkID(r)
	if err != nil {
		helper.ResponseBadRequest(w, "作品IDが不正です")
		return
	}

	pageSize, page, _ := parsePaginationQuery(r)

	tsumikis, err := wh.repositories.Tsumiki.GetTsumikis(userID, pageSize, page, nil, &workID, "")
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	helper.ResponseOk(w, tsumikis)
}

func (wh *workHandlerImpl) CreateWork(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.ResponseUnauthorized(w, "認証情報が見つかりません")
		return
	}

	var req workRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.ResponseBadRequest(w, "リクエストボディが不正です")
		return
	}

	work, err := wh.repositories.Work.CreateWork(userID, req.Title, req.Description)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	helper.ResponseOk(w, work)
}

func (wh *workHandlerImpl) EditWork(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.ResponseUnauthorized(w, "認証情報が見つかりません")
		return
	}
	workID, err := parseWorkID(r)
	if err != nil {
		helper.ResponseBadRequest(w, "作品IDが不正です")
		return
	}
	var req workRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.ResponseBadRequest(w, "リクエストボディが不正です")
		return
	}

	work, err := wh.repositories.Work.GetWork(workID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if work == nil {
		helper.ResponseNotFound(w, "作品が見つかりません")
		return
	}
	if work.Owner.ID != userID {
		helper.ResponseForbidden(w, "この作品の作成者ではありません")
		return
	}
	updatedWork, err := wh.repositories.Work.UpdateWork(workID, req.Title, req.Description)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	helper.ResponseOk(w, updatedWork)
}

func (wh *workHandlerImpl) DeleteWork(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.ResponseUnauthorized(w, "認証情報が見つかりません")
		return
	}
	workID, err := parseWorkID(r)
	if err != nil {
		helper.ResponseBadRequest(w, "作品IDが不正です")
		return
	}

	work, err := wh.repositories.Work.GetWork(workID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if work == nil {
		helper.ResponseNotFound(w, "作品が見つかりません")
		return
	}
	if work.Owner.ID != userID {
		helper.ResponseForbidden(w, "この作品の作成者ではありません")
		return
	}
	if err := wh.repositories.Work.DeleteWork(workID); err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	helper.ResponseOk(w, nil)
}

func parseWorkID(r *http.Request) (int, error) {
	return strconv.Atoi(chi.URLParam(r, "workId"))
}
