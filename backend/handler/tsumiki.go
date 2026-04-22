package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"net/http"
	"strconv"
	"tsumiki/helper"
	"tsumiki/media"
	"tsumiki/middleware"
	"tsumiki/repository"
	"tsumiki/schema"

	"github.com/go-chi/chi/v5"
)

const (
	mediaMaxBytes5MB   int64 = 5 << 20
	mediaMaxBytes100MB int64 = 100 << 20
)

type TsumikiHandler interface {
	GetMyTsumikis(w http.ResponseWriter, r *http.Request)
	GetUserTsumikis(w http.ResponseWriter, r *http.Request)
	GetTsumikis(w http.ResponseWriter, r *http.Request)
	GetSpecifiedTsumiki(w http.ResponseWriter, r *http.Request)
	GetBlocks(w http.ResponseWriter, r *http.Request)
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

type writeBlockRequest struct {
	Message    *string `json:"message"`
	Percentage int     `json:"percentage"`
	Condition  int     `json:"condition"`
	MediaIDs   []int   `json:"media_ids"`
}

type addBlockRequest struct {
	writeBlockRequest
	LatestBlockID *int `json:"latest_block_id"`
}

type createTsumikiRequest struct {
	Title      string `json:"title"`
	Visibility string `json:"visibility"`
	WorkID     *int   `json:"work_id"`
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

func (th *tsumikiHandlerImpl) GetBlocks(w http.ResponseWriter, r *http.Request) {
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

	pageSize, page, _ := parsePaginationQuery(r)
	blocks, err := th.repositories.Tsumiki.GetTsumikiBlocks(tsumikiID, pageSize, page)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

	helper.ResponseOk(w, blocks)
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
	err := th.repositories.RunInTx(func(txRepos *repository.Repositories) error {
		var err error
		tsumiki, err = txRepos.Tsumiki.CreateTsumiki(userID, req.Title, req.Visibility, req.WorkID)
		return err
	})
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}

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

	// 一番大きい上限で一度制限をかける
	if err := r.ParseMultipartForm(mediaMaxBytes100MB); err != nil {
		helper.ResponseBadRequest(w, "マルチパートフォームの解析に失敗しました")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		helper.ResponseBadRequest(w, "ファイルが見つかりません")
		return
	}
	defer file.Close()

	rawContentType := header.Header.Get("Content-Type")
	mediaType, maxBytes, ext, err := resolveMediaType(rawContentType)
	if err != nil {
		helper.ResponseBadRequest(w, err.Error())
		return
	}

	data, err := io.ReadAll(io.LimitReader(file, maxBytes+1))
	if err != nil {
		helper.ResponseInternalServerError(w, "ファイルの読み込みに失敗しました")
		return
	}
	if int64(len(data)) > maxBytes {
		helper.ResponseBadRequest(w, "ファイルサイズが上限を超えています")
		return
	}

	if mediaType == "image" {
		cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
		if err != nil {
			helper.ResponseBadRequest(w, "画像の解析に失敗しました")
			return
		}
		if cfg.Width > 4096 || cfg.Height > 4096 {
			helper.ResponseBadRequest(w, "画像サイズは4096x4096以内にしてください")
			return
		}
	}

	// TODO: audio/videoの制限もかける
	// video 3分以内
	// audio 5分以内

	storedPath, err := th.media.UploadTsumikiMedia(r.Context(), tsumikiID, bytes.NewReader(data), rawContentType, ext)
	if err != nil {
		fmt.Println("S3アップロードエラー: ", err)
		helper.ResponseInternalServerError(w, "ファイルのアップロードに失敗しました")
		return
	}

	createdMedia, err := th.repositories.TsumikiBlockMedia.CreateMedia(storedPath, mediaType)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	createdMedia.Url = th.media.ResolveURL(createdMedia.Url)

	helper.ResponseOk(w, createdMedia)
}

var mediaContentTypes = map[string]struct {
	mediaType string
	maxBytes  int64
	ext       string
}{
	"image/jpeg":      {"image", mediaMaxBytes5MB, ".jpg"},
	"image/png":       {"image", mediaMaxBytes5MB, ".png"},
	"image/gif":       {"image", mediaMaxBytes5MB, ".gif"},
	"audio/mpeg":      {"audio", mediaMaxBytes5MB, ".mp3"},
	"audio/wav":       {"audio", mediaMaxBytes5MB, ".wav"},
	"audio/ogg":       {"audio", mediaMaxBytes5MB, ".ogg"},
	"audio/aac":       {"audio", mediaMaxBytes5MB, ".aac"},
	"video/mp4":       {"video", mediaMaxBytes100MB, ".mp4"},
	"video/webm":      {"video", mediaMaxBytes100MB, ".webm"},
	"video/quicktime": {"video", mediaMaxBytes100MB, ".mov"},
}

func resolveMediaType(contentType string) (mediaType string, maxBytes int64, ext string, err error) {
	mimeType, _, parseErr := mime.ParseMediaType(contentType)
	if parseErr != nil {
		return "", 0, "", fmt.Errorf("Content-Typeが不正です")
	}
	info, ok := mediaContentTypes[mimeType]
	if !ok {
		return "", 0, "", fmt.Errorf("対応していないファイル形式です")
	}
	return info.mediaType, info.maxBytes, info.ext, nil
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
	var req addBlockRequest
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

	latestBlockID, err := th.repositories.TsumikiBlock.GetLatestBlockID(tsumikiID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	neitherSet := latestBlockID == nil && req.LatestBlockID == nil
	bothMatch := latestBlockID != nil && req.LatestBlockID != nil && *latestBlockID == *req.LatestBlockID
	if !neitherSet && !bothMatch {
		helper.ResponseConflict(w, "最新ブロックIDが一致しません")
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
	var req writeBlockRequest
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
		block, err = txRepos.TsumikiBlock.UpdateBlock(blockID, req.Message, req.Percentage, req.Condition)
		if err != nil {
			return err
		}
		if block == nil {
			return nil
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
