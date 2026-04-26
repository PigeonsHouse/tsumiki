package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"tsumiki/handler"
	mediamock "tsumiki/media/mock"
	"tsumiki/repository"
	repomock "tsumiki/repository/mock"
	"tsumiki/schema"

	"go.uber.org/mock/gomock"
)

func TestPostThumbnail(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := sampleUser().ID
	thumbnail := sampleThumbnail()
	storedPath := "thumbnails/1/abc.png"

	mockThumbnail := repomock.NewMockThumbnailRepository(ctrl)
	mockThumbnail.EXPECT().Create(userID, storedPath).Return(thumbnail, nil)

	mockMedia := mediamock.NewMockMediaService(ctrl)
	mockMedia.EXPECT().
		UploadThumbnail(gomock.Any(), userID, gomock.Any(), "image/png", ".png").
		Return(storedPath, nil)
	mockMedia.EXPECT().ResolveURL(thumbnail.Url).Return("https://cdn.example.com/" + thumbnail.Url)

	repos := &repository.Repositories{Thumbnail: mockThumbnail}
	h := handler.NewThumbnailHandler(repos, mockMedia)

	req := newMultipartRequest(t, "thumbnail", "test.png", "image/png", createMinimalPNG())
	req = withUserID(req, userID)
	w := httptest.NewRecorder()

	h.PostThumbnail(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d (body: %s)", http.StatusOK, w.Code, w.Body.String())
	}

	var resp schema.ThumbnailUpload
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID != thumbnail.ID {
		t.Errorf("ID: want %d, got %d", thumbnail.ID, resp.ID)
	}
}

func TestPostThumbnail_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)

	h := handler.NewThumbnailHandler(&repository.Repositories{}, mediamock.NewMockMediaService(ctrl))

	req := newMultipartRequest(t, "thumbnail", "test.png", "image/png", createMinimalPNG())
	w := httptest.NewRecorder()

	h.PostThumbnail(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestPostThumbnail_InvalidContentType(t *testing.T) {
	ctrl := gomock.NewController(t)

	h := handler.NewThumbnailHandler(&repository.Repositories{}, mediamock.NewMockMediaService(ctrl))

	req := newMultipartRequest(t, "thumbnail", "test.txt", "text/plain", []byte("hello"))
	req = withUserID(req, sampleUser().ID)
	w := httptest.NewRecorder()

	h.PostThumbnail(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: want %d, got %d", http.StatusBadRequest, w.Code)
	}
}
