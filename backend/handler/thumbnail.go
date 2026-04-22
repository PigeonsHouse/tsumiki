package handler

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"tsumiki/helper"
	"tsumiki/media"
	"tsumiki/middleware"
	"tsumiki/repository"
)

type ThumbnailHandler interface {
	PostThumbnail(w http.ResponseWriter, r *http.Request)
}

type thumbnailHandlerImpl struct {
	repositories *repository.Repositories
	media        media.MediaService
}

func NewThumbnailHandler(repos *repository.Repositories, mediaSvc media.MediaService) ThumbnailHandler {
	return &thumbnailHandlerImpl{repositories: repos, media: mediaSvc}
}

func (th *thumbnailHandlerImpl) PostThumbnail(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.ResponseUnauthorized(w, "認証情報が見つかりません")
		return
	}

	if err := r.ParseMultipartForm(mediaMaxBytes5MB); err != nil {
		helper.ResponseBadRequest(w, "マルチパートフォームの解析に失敗しました")
		return
	}

	data, rawContentType, ext, err := parseThumbnailFile(r)
	if err != nil {
		helper.ResponseBadRequest(w, err.Error())
		return
	}

	path, err := th.media.UploadThumbnail(r.Context(), userID, bytes.NewReader(data), rawContentType, ext)
	if err != nil {
		fmt.Println("S3アップロードエラー: ", err)
		helper.ResponseInternalServerError(w, "ファイルのアップロードに失敗しました")
		return
	}

	thumbnail, err := th.repositories.Thumbnail.Create(userID, path)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return
	}
	thumbnail.Url = th.media.ResolveURL(thumbnail.Url)

	helper.ResponseOk(w, thumbnail)
}

func validateThumbnailAvailable(repos *repository.Repositories, thumbnailID int, w http.ResponseWriter) error {
	thumbnail, err := repos.Thumbnail.Get(thumbnailID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return err
	}
	if thumbnail == nil {
		helper.ResponseBadRequest(w, "サムネイルが見つかりません")
		return fmt.Errorf("not found")
	}
	inUse, err := repos.Thumbnail.IsInUse(thumbnailID)
	if err != nil {
		fmt.Println("DBエラー: ", err)
		helper.ResponseInternalServerError(w, "DBエラー")
		return err
	}
	if inUse {
		helper.ResponseConflict(w, "このサムネイルは既に使用されています")
		return fmt.Errorf("in use")
	}
	return nil
}

func parseThumbnailFile(r *http.Request) (data []byte, rawContentType, ext string, err error) {
	file, header, formErr := r.FormFile("thumbnail")
	if formErr != nil {
		return nil, "", "", fmt.Errorf("サムネイルファイルが見つかりません")
	}
	defer file.Close()

	rawContentType = header.Header.Get("Content-Type")
	mediaType, maxBytes, ext, resolveErr := resolveMediaType(rawContentType)
	if resolveErr != nil {
		return nil, "", "", resolveErr
	}
	if mediaType != "image" {
		return nil, "", "", fmt.Errorf("サムネイルは画像ファイルのみ対応しています")
	}

	data, err = io.ReadAll(io.LimitReader(file, maxBytes+1))
	if err != nil {
		return nil, "", "", fmt.Errorf("ファイルの読み込みに失敗しました")
	}
	if int64(len(data)) > maxBytes {
		return nil, "", "", fmt.Errorf("ファイルサイズが上限を超えています")
	}

	cfg, _, imgErr := image.DecodeConfig(bytes.NewReader(data))
	if imgErr != nil {
		return nil, "", "", fmt.Errorf("画像の解析に失敗しました")
	}
	if cfg.Width > 4096 || cfg.Height > 4096 {
		return nil, "", "", fmt.Errorf("画像サイズは4096x4096以内にしてください")
	}

	return data, rawContentType, ext, nil
}
