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

	"github.com/go-chi/chi/v5"
)

type TsumikiHandler interface {
	GetMyTsumikis(w http.ResponseWriter, r *http.Request)
	GetUserTsumikis(w http.ResponseWriter, r *http.Request)
	GetTsumikis(w http.ResponseWriter, r *http.Request)
	GetSpecifiedTsumiki(w http.ResponseWriter, r *http.Request)
	CreateTsumiki(w http.ResponseWriter, r *http.Request)
	EditTsumiki(w http.ResponseWriter, r *http.Request)
	DeleteTsumiki(w http.ResponseWriter, r *http.Request)
	PostMedia(w http.ResponseWriter, r *http.Request)
	AddBlock(w http.ResponseWriter, r *http.Request)
	EditBlock(w http.ResponseWriter, r *http.Request)
	OmitBlock(w http.ResponseWriter, r *http.Request)
}

type tsumikiHandlerImpl struct {
	tsumikiRepo repository.TsumikiRepository
	blockRepo   repository.TsumikiBlockRepository
	mediaRepo   repository.TsumikiBlockMediaRepository
	media       media.MediaService
}

func NewTsumikiHandler(
	tsumikiRepo repository.TsumikiRepository,
	blockRepo repository.TsumikiBlockRepository,
	mediaRepo repository.TsumikiBlockMediaRepository,
	mediaSvc media.MediaService,
) TsumikiHandler {
	return &tsumikiHandlerImpl{
		tsumikiRepo: tsumikiRepo,
		blockRepo:   blockRepo,
		mediaRepo:   mediaRepo,
		media:       mediaSvc,
	}
}

type blockRequest struct {
	Message    *string `json:"message"`
	Percentage int     `json:"percentage"`
	Condition  int     `json:"condition"`
	MediaIDs   []int   `json:"media_ids"`
}

type createTsumikiRequest struct {
	Title      string       `json:"title"`
	Visibility string       `json:"visibility"`
	WorkID     *int         `json:"work_id"`
	Block      blockRequest `json:"block"`
}

type editTsumikiRequest struct {
	Title      string `json:"title"`
	Visibility string `json:"visibility"`
	WorkID     *int   `json:"work_id"`
}

func (th *tsumikiHandlerImpl) GetMyTsumikis(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.ResponseUnauthorized(w, "認証情報が見つかりません")
		return
	}

	pageSize, page, keyword := parsePaginationQuery(r)
	tsumikis, err := th.tsumikiRepo.GetTsumikis(&userID, pageSize, page, &userID, nil, keyword)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	helper.ResponseOk(w, tsumikis)
}

func (th *tsumikiHandlerImpl) GetUserTsumikis(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetOptionalUserIDFromContext(r.Context())
	idStr := chi.URLParam(r, "userID")
	authorID, err := strconv.Atoi(idStr)
	if err != nil {
		helper.ResponseBadRequest(w, "ユーザIDが不正です")
		return
	}

	pageSize, page, keyword := parsePaginationQuery(r)

	tsumikis, err := th.tsumikiRepo.GetTsumikis(userID, pageSize, page, &authorID, nil, keyword)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	helper.ResponseOk(w, tsumikis)
}

func (th *tsumikiHandlerImpl) GetTsumikis(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetOptionalUserIDFromContext(r.Context())
	pageSize, page, keyword := parsePaginationQuery(r)

	tsumikis, err := th.tsumikiRepo.GetTsumikis(userID, pageSize, page, nil, nil, keyword)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	helper.ResponseOk(w, tsumikis)
}

func (th *tsumikiHandlerImpl) GetSpecifiedTsumiki(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetOptionalUserIDFromContext(r.Context())
	tsumikiID, err := parseTsumikiID(r)
	if err != nil {
		helper.ResponseBadRequest(w, "積み木IDが不正です")
		return
	}

	tsumiki, err := th.tsumikiRepo.GetTsumiki(userID, tsumikiID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if tsumiki == nil {
		helper.ResponseNotFound(w, "積み木が見つかりません")
		return
	}

	helper.ResponseOk(w, tsumiki)
}

func (th *tsumikiHandlerImpl) CreateTsumiki(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.ResponseUnauthorized(w, "認証情報が見つかりません")
		return
	}

	var req createTsumikiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.ResponseBadRequest(w, "リクエストボディが不正です")
		return
	}

	// todo: トランザクションを貼る
	tsumiki, err := th.tsumikiRepo.CreateTsumiki(userID, req.Title, req.Visibility, req.WorkID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	_, err = th.blockRepo.CreateBlock(tsumiki.ID, req.Block.Message, req.Block.Percentage, req.Block.Condition, req.Block.MediaIDs)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	// メディアとブロックを紐付ける

	helper.ResponseOk(w, tsumiki)
}

func (th *tsumikiHandlerImpl) EditTsumiki(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.ResponseUnauthorized(w, "認証情報が見つかりません")
		return
	}
	tsumikiID, err := parseTsumikiID(r)
	if err != nil {
		helper.ResponseBadRequest(w, "積み木IDが不正です")
		return
	}
	var req editTsumikiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.ResponseBadRequest(w, "リクエストボディが不正です")
		return
	}

	tsumiki, err := th.tsumikiRepo.GetTsumiki(&userID, tsumikiID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if tsumiki == nil {
		helper.ResponseNotFound(w, "積み木が見つかりません")
		return
	}
	if tsumiki.User.ID != userID {
		helper.ResponseForbidden(w, "この積み木の作成者ではありません")
		return
	}

	updatedTsumiki, err := th.tsumikiRepo.UpdateTsumiki(tsumikiID, req.Title, req.Visibility, req.WorkID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	helper.ResponseOk(w, updatedTsumiki)
}

func (th *tsumikiHandlerImpl) DeleteTsumiki(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.ResponseUnauthorized(w, "認証情報が見つかりません")
		return
	}
	_ = userID

	tsumikiID, err := parseTsumikiID(r)
	if err != nil {
		helper.ResponseBadRequest(w, "積み木IDが不正です")
		return
	}

	tsumiki, err := th.tsumikiRepo.GetTsumiki(&userID, tsumikiID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if tsumiki == nil {
		helper.ResponseNotFound(w, "積み木が見つかりません")
		return
	}
	if tsumiki.User.ID != userID {
		helper.ResponseForbidden(w, "この積み木の作成者ではありません")
		return
	}

	if err := th.tsumikiRepo.DeleteTsumiki(tsumikiID); err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	helper.ResponseOk(w, nil)
}

func (th *tsumikiHandlerImpl) PostMedia(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.ResponseUnauthorized(w, "認証情報が見つかりません")
		return
	}
	tsumikiID, err := parseTsumikiID(r)
	if err != nil {
		helper.ResponseBadRequest(w, "積み木IDが不正です")
		return
	}

	tsumiki, err := th.tsumikiRepo.GetTsumiki(&userID, tsumikiID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if tsumiki == nil {
		helper.ResponseNotFound(w, "積み木が見つかりません")
		return
	}
	if tsumiki.User.ID != userID {
		helper.ResponseForbidden(w, "この積み木の作成者ではありません")
		return
	}

	// TODO: ファイルタイプ・サイズ検証、S3アップロード、DBレコード作成
	helper.ResponseOk(w, nil)
}

func (th *tsumikiHandlerImpl) AddBlock(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.ResponseUnauthorized(w, "認証情報が見つかりません")
		return
	}
	_ = userID
	tsumikiID, err := parseTsumikiID(r)
	if err != nil {
		helper.ResponseBadRequest(w, "積み木IDが不正です")
		return
	}
	var req blockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.ResponseBadRequest(w, "リクエストボディが不正です")
		return
	}

	tsumiki, err := th.tsumikiRepo.GetTsumiki(&userID, tsumikiID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if tsumiki == nil {
		helper.ResponseNotFound(w, "積み木が見つかりません")
		return
	}
	if tsumiki.User.ID != userID {
		helper.ResponseForbidden(w, "この積み木の作成者ではありません")
		return
	}
	// メディア操作までトランザクションを貼る
	block, err := th.blockRepo.CreateBlock(tsumikiID, req.Message, req.Percentage, req.Condition, req.MediaIDs)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	// メディアの紐付け

	helper.ResponseOk(w, block)
}

func (th *tsumikiHandlerImpl) EditBlock(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.ResponseUnauthorized(w, "認証情報が見つかりません")
		return
	}
	tsumikiID, err := parseTsumikiID(r)
	if err != nil {
		helper.ResponseBadRequest(w, "積み木IDが不正です")
		return
	}
	blockID, err := strconv.Atoi(chi.URLParam(r, "blockID"))
	if err != nil {
		helper.ResponseBadRequest(w, "ブロックIDが不正です")
		return
	}
	var req blockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.ResponseBadRequest(w, "リクエストボディが不正です")
		return
	}

	tsumiki, err := th.tsumikiRepo.GetTsumiki(&userID, tsumikiID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if tsumiki == nil {
		helper.ResponseNotFound(w, "積み木が見つかりません")
		return
	}
	isBelong, err := th.blockRepo.IsBelongToTsumiki(tsumikiID, blockID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if !isBelong {
		helper.ResponseBadRequest(w, "このブロックはこの積み木に含まれていません")
		return
	}
	if tsumiki.User.ID != userID {
		helper.ResponseForbidden(w, "この積み木の作成者ではありません")
		return
	}
	// メディア操作までトランザクションを貼る
	block, err := th.blockRepo.UpdateBlock(blockID, req.Message, req.Percentage, req.Condition, req.MediaIDs)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if block == nil {
		helper.ResponseNotFound(w, "ブロックが見つかりません")
		return
	}
	// メディアの紐付け切り離しの整理

	helper.ResponseOk(w, block)
}

func (th *tsumikiHandlerImpl) OmitBlock(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.ResponseUnauthorized(w, "認証情報が見つかりません")
		return
	}
	tsumikiID, err := parseTsumikiID(r)
	if err != nil {
		helper.ResponseBadRequest(w, "積み木IDが不正です")
		return
	}
	blockID, err := strconv.Atoi(chi.URLParam(r, "blockID"))
	if err != nil {
		helper.ResponseBadRequest(w, "ブロックIDが不正です")
		return
	}

	tsumiki, err := th.tsumikiRepo.GetTsumiki(&userID, tsumikiID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if tsumiki == nil {
		helper.ResponseNotFound(w, "積み木が見つかりません")
		return
	}
	isBelong, err := th.blockRepo.IsBelongToTsumiki(tsumikiID, blockID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if !isBelong {
		helper.ResponseBadRequest(w, "このブロックはこの積み木に含まれていません")
		return
	}
	if tsumiki.User.ID != userID {
		helper.ResponseForbidden(w, "この積み木の作成者ではありません")
		return
	}

	if err := th.blockRepo.SoftDeleteBlock(blockID); err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	helper.ResponseOk(w, nil)
}

func parseTsumikiID(r *http.Request) (int, error) {
	return strconv.Atoi(chi.URLParam(r, "tsumikiID"))
}

func parsePaginationQuery(r *http.Request) (pageSize, page int, keyword string) {
	pageSize, _ = strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 50 {
		pageSize = 50
	}
	page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	keyword = r.URL.Query().Get("keyword")
	return
}
