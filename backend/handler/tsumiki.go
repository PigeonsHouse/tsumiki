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
	repositories *repository.Repositories
	media        media.MediaService
}

func NewTsumikiHandler(repos *repository.Repositories, mediaSvc media.MediaService) TsumikiHandler {
	return &tsumikiHandlerImpl{
		repositories: repos,
		media:        mediaSvc,
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

type createTsumikiResponse struct {
	Tsumiki  *schema.Tsumiki      `json:"tsumiki"`
	NewBlock *schema.TsumikiBlock `json:"new_block"`
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
	tsumikis, err := th.repositories.Tsumiki.GetTsumikis(&userID, pageSize, page, &userID, nil, keyword)
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

	tsumikis, err := th.repositories.Tsumiki.GetTsumikis(userID, pageSize, page, &authorID, nil, keyword)
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

	tsumikis, err := th.repositories.Tsumiki.GetTsumikis(userID, pageSize, page, nil, nil, keyword)
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

	tsumiki, err := th.repositories.Tsumiki.GetTsumiki(userID, tsumikiID)
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

	var tsumiki *schema.Tsumiki
	var newBlock *schema.TsumikiBlock
	err := th.repositories.RunInTx(func(txRepos *repository.Repositories) error {
		var err error
		tsumiki, err = th.repositories.Tsumiki.CreateTsumiki(userID, req.Title, req.Visibility, req.WorkID)
		if err != nil {
			return err
		}
		newBlock, err = th.repositories.TsumikiBlock.CreateBlock(tsumiki.ID, req.Block.Message, req.Block.Percentage, req.Block.Condition)
		if err != nil {
			return err
		}
		medias, err := th.repositories.TsumikiBlockMedia.SetMediaRelation(newBlock.ID, req.Block.MediaIDs)
		if err != nil {
			return err
		}
		newBlock.Medias = medias
		return nil
	})
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	helper.ResponseOk(w, createTsumikiResponse{
		Tsumiki:  tsumiki,
		NewBlock: newBlock,
	})
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

	tsumiki, err := th.repositories.Tsumiki.GetTsumiki(&userID, tsumikiID)
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

	updatedTsumiki, err := th.repositories.Tsumiki.UpdateTsumiki(tsumikiID, req.Title, req.Visibility, req.WorkID)
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

	tsumiki, err := th.repositories.Tsumiki.GetTsumiki(&userID, tsumikiID)
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

	if err := th.repositories.Tsumiki.DeleteTsumiki(tsumikiID); err != nil {
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

	tsumiki, err := th.repositories.Tsumiki.GetTsumiki(&userID, tsumikiID)
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

	tsumiki, err := th.repositories.Tsumiki.GetTsumiki(&userID, tsumikiID)
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
	var block *schema.TsumikiBlock
	err = th.repositories.RunInTx(func(txRepos *repository.Repositories) error {
		block, err = txRepos.TsumikiBlock.CreateBlock(tsumikiID, req.Message, req.Percentage, req.Condition)
		if err != nil {
			return err
		}
		medias, err := txRepos.TsumikiBlockMedia.SetMediaRelation(block.ID, req.MediaIDs)
		if err != nil {
			return err
		}
		block.Medias = medias
		return nil
	})
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

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

	tsumiki, err := th.repositories.Tsumiki.GetTsumiki(&userID, tsumikiID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if tsumiki == nil {
		helper.ResponseNotFound(w, "積み木が見つかりません")
		return
	}
	isBelong, err := th.repositories.TsumikiBlock.IsBelongToTsumiki(tsumikiID, blockID)
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

	var block *schema.TsumikiBlock
	err = th.repositories.RunInTx(func(txRepos *repository.Repositories) error {
		block, err = th.repositories.TsumikiBlock.UpdateBlock(blockID, req.Message, req.Percentage, req.Condition)
		if err != nil {
			return err
		}
		if block == nil {
			return nil
		}
		medias, err := th.repositories.TsumikiBlockMedia.SetMediaRelation(block.ID, req.MediaIDs)
		if err != nil {
			return err
		}
		block.Medias = medias
		return nil
	})
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if block == nil {
		helper.ResponseNotFound(w, "ブロックが見つかりません")
		return
	}

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

	tsumiki, err := th.repositories.Tsumiki.GetTsumiki(&userID, tsumikiID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	if tsumiki == nil {
		helper.ResponseNotFound(w, "積み木が見つかりません")
		return
	}
	isBelong, err := th.repositories.TsumikiBlock.IsBelongToTsumiki(tsumikiID, blockID)
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

	if err := th.repositories.TsumikiBlock.SoftDeleteBlock(blockID); err != nil {
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
