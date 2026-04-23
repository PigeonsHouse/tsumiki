package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tsumiki/helper"
	"tsumiki/media"
	"tsumiki/middleware"
	"tsumiki/repository"
	"tsumiki/schema"

	"github.com/go-chi/chi/v5"
)

type WorkHandler interface {
	GetWorks(w http.ResponseWriter, r *http.Request)
	GetSpecifiedWork(w http.ResponseWriter, r *http.Request)
	GetWorkTsumiki(w http.ResponseWriter, r *http.Request)
	CreateWork(w http.ResponseWriter, r *http.Request)
	EditWork(w http.ResponseWriter, r *http.Request)
	UpdateWorkThumbnail(w http.ResponseWriter, r *http.Request)
	DeleteWork(w http.ResponseWriter, r *http.Request)
}

type workHandlerImpl struct {
	repositories *repository.Repositories
	media        media.MediaService
}

func NewWorkHandler(repos *repository.Repositories, mediaSvc media.MediaService) WorkHandler {
	return &workHandlerImpl{
		repositories: repos,
		media:        mediaSvc,
	}
}

type createWorkRequest struct {
	Title       string `json:"title"`
	Visibility  string `json:"visibility"`
	Description string `json:"description"`
	ThumbnailID *int   `json:"thumbnail_id"`
}

type editWorkRequest struct {
	Title       string `json:"title"`
	Visibility  string `json:"visibility"`
	Description string `json:"description"`
}

func (wh *workHandlerImpl) GetWorks(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetOptionalUserIDFromContext(r.Context())
	pageSize, page, _ := parsePaginationQuery(r)

	works, err := wh.repositories.Work.GetWorks(userID, pageSize, page)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	for i := range works {
		helper.ResolveWorkURLs(&works[i], wh.media)
	}
	helper.ResponseOk(w, works)
}

func (wh *workHandlerImpl) GetSpecifiedWork(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetOptionalUserIDFromContext(r.Context())
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
	if !wh.canAccessWork(w, work, userID) {
		return
	}

	helper.ResolveWorkURLs(work, wh.media)
	helper.ResponseOk(w, work)
}

func (wh *workHandlerImpl) GetWorkTsumiki(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetOptionalUserIDFromContext(r.Context())
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
	if !wh.canAccessWork(w, work, userID) {
		return
	}

	pageSize, page, _ := parsePaginationQuery(r)
	tsumikis, err := wh.repositories.Tsumiki.GetTsumikis(userID, pageSize, page, nil, &workID, "")
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	for i := range tsumikis {
		helper.ResolveTsumikiURLs(&tsumikis[i], wh.media)
	}
	helper.ResponseOk(w, tsumikis)
}

// canAccessWork は作品のvisibilityを確認し、アクセス不可なら403を返してfalseを返す
func (wh *workHandlerImpl) canAccessWork(w http.ResponseWriter, work *schema.Work, userID *int) bool {
	if work.Visibility == "public" {
		return true
	}
	if userID == nil {
		helper.ResponseForbidden(w, "この作品は限定公開です")
		return false
	}
	viewer, err := wh.repositories.User.FindByID(*userID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return false
	}
	if viewer == nil || viewer.GuildID == nil || work.Owner.GuildID == nil || *viewer.GuildID != *work.Owner.GuildID {
		helper.ResponseForbidden(w, "この作品は限定公開です")
		return false
	}
	return true
}

func (wh *workHandlerImpl) CreateWork(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.ResponseUnauthorized(w, "認証情報が見つかりません")
		return
	}

	var req createWorkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.ResponseBadRequest(w, "リクエストボディが不正です")
		return
	}
	if len(req.Title) > maxTitleLength {
		helper.ResponseBadRequest(w, "タイトルは200文字以内にしてください")
		return
	}
	if len(req.Description) > maxDescriptionLength {
		helper.ResponseBadRequest(w, "説明は4000文字以内にしてください")
		return
	}

	if req.ThumbnailID != nil {
		if err := validateThumbnailAvailable(wh.repositories, *req.ThumbnailID, w); err != nil {
			return
		}
	}

	work, err := wh.repositories.Work.CreateWork(userID, req.Title, req.Visibility, req.Description, req.ThumbnailID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	helper.ResolveWorkURLs(work, wh.media)
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
	var req editWorkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.ResponseBadRequest(w, "リクエストボディが不正です")
		return
	}
	if len(req.Title) > maxTitleLength {
		helper.ResponseBadRequest(w, "タイトルは200文字以内にしてください")
		return
	}
	if len(req.Description) > maxDescriptionLength {
		helper.ResponseBadRequest(w, "説明は4000文字以内にしてください")
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
	updatedWork, err := wh.repositories.Work.UpdateWork(workID, req.Title, req.Visibility, req.Description)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	helper.ResolveWorkURLs(updatedWork, wh.media)
	helper.ResponseOk(w, updatedWork)
}

func (wh *workHandlerImpl) UpdateWorkThumbnail(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		ThumbnailID int `json:"thumbnail_id"`
	}
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

	if err := validateThumbnailAvailable(wh.repositories, req.ThumbnailID, w); err != nil {
		return
	}

	if err := wh.repositories.Work.UpdateWorkThumbnail(workID, req.ThumbnailID); err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	helper.ResponseOk(w, nil)
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
	return strconv.Atoi(chi.URLParam(r, "workID"))
}
