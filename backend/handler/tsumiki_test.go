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

func TestGetTsumikis(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleTsumiki()

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().
		GetTsumikis(gomock.Nil(), gomock.Any(), gomock.Any(), gomock.Nil(), gomock.Nil(), gomock.Any()).
		Return([]schema.Tsumiki{*expected}, nil)

	mockMedia := mediamock.NewMockMediaService(ctrl)
	mockMedia.EXPECT().ResolveURL(gomock.Any()).Return("https://cdn.example.com/test").AnyTimes()

	h := handler.NewTsumikiHandler(&repository.Repositories{Tsumiki: mockTsumiki}, mockMedia)

	req := httptest.NewRequest(http.MethodGet, "/tsumikis", nil)
	w := httptest.NewRecorder()

	h.GetTsumikis(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}

	var resp []schema.Tsumiki
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("want 1 tsumiki, got %d", len(resp))
	}
	if resp[0].ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, resp[0].ID)
	}
}

func TestGetMyTsumikis(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := sampleUser().ID
	expected := sampleTsumiki()

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().
		GetTsumikis(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Nil(), gomock.Any()).
		Return([]schema.Tsumiki{*expected}, nil)

	mockMedia := mediamock.NewMockMediaService(ctrl)
	mockMedia.EXPECT().ResolveURL(gomock.Any()).Return("https://cdn.example.com/test").AnyTimes()

	h := handler.NewTsumikiHandler(&repository.Repositories{Tsumiki: mockTsumiki}, mockMedia)

	req := withUserID(httptest.NewRequest(http.MethodGet, "/tsumikis/me", nil), userID)
	w := httptest.NewRecorder()

	h.GetMyTsumikis(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetMyTsumikis_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)

	h := handler.NewTsumikiHandler(&repository.Repositories{}, mediamock.NewMockMediaService(ctrl))

	req := httptest.NewRequest(http.MethodGet, "/tsumikis/me", nil)
	w := httptest.NewRecorder()

	h.GetMyTsumikis(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestGetUserTsumikis(t *testing.T) {
	ctrl := gomock.NewController(t)
	authorID := sampleUser().ID
	expected := sampleTsumiki()

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().
		GetTsumikis(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Nil(), gomock.Any()).
		Return([]schema.Tsumiki{*expected}, nil)

	mockMedia := mediamock.NewMockMediaService(ctrl)
	mockMedia.EXPECT().ResolveURL(gomock.Any()).Return("https://cdn.example.com/test").AnyTimes()

	h := handler.NewTsumikiHandler(&repository.Repositories{Tsumiki: mockTsumiki}, mockMedia)

	req := withChiParam(httptest.NewRequest(http.MethodGet, "/users/1/tsumikis", nil), "userID", "1")
	w := httptest.NewRecorder()

	h.GetUserTsumikis(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}

	_ = authorID
}

func TestGetUserTsumikis_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)

	h := handler.NewTsumikiHandler(&repository.Repositories{}, mediamock.NewMockMediaService(ctrl))

	req := withChiParam(httptest.NewRequest(http.MethodGet, "/users/abc/tsumikis", nil), "userID", "abc")
	w := httptest.NewRecorder()

	h.GetUserTsumikis(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: want %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetSpecifiedTsumiki(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleTsumiki()

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().GetTsumiki(gomock.Nil(), expected.ID).Return(expected, nil)

	mockMedia := mediamock.NewMockMediaService(ctrl)
	mockMedia.EXPECT().ResolveURL(gomock.Any()).Return("https://cdn.example.com/test").AnyTimes()

	h := handler.NewTsumikiHandler(&repository.Repositories{Tsumiki: mockTsumiki}, mockMedia)

	req := withChiParam(httptest.NewRequest(http.MethodGet, "/tsumikis/3", nil), "tsumikiID", "3")
	w := httptest.NewRecorder()

	h.GetSpecifiedTsumiki(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}

	var resp schema.Tsumiki
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, resp.ID)
	}
}

func TestGetSpecifiedTsumiki_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)

	h := handler.NewTsumikiHandler(&repository.Repositories{}, mediamock.NewMockMediaService(ctrl))

	req := withChiParam(httptest.NewRequest(http.MethodGet, "/tsumikis/abc", nil), "tsumikiID", "abc")
	w := httptest.NewRecorder()

	h.GetSpecifiedTsumiki(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: want %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetSpecifiedTsumiki_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().GetTsumiki(gomock.Nil(), 999).Return(nil, nil)

	h := handler.NewTsumikiHandler(&repository.Repositories{Tsumiki: mockTsumiki}, mediamock.NewMockMediaService(ctrl))

	req := withChiParam(httptest.NewRequest(http.MethodGet, "/tsumikis/999", nil), "tsumikiID", "999")
	w := httptest.NewRecorder()

	h.GetSpecifiedTsumiki(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: want %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetBlocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	tsumiki := sampleTsumiki()

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().GetTsumiki(gomock.Nil(), tsumiki.ID).Return(tsumiki, nil)
	mockTsumiki.EXPECT().GetTsumikiBlocks(tsumiki.ID, gomock.Any(), gomock.Any()).
		Return([]schema.TsumikiBlockView{}, nil)

	mockMedia := mediamock.NewMockMediaService(ctrl)

	h := handler.NewTsumikiHandler(&repository.Repositories{Tsumiki: mockTsumiki}, mockMedia)

	req := withChiParam(httptest.NewRequest(http.MethodGet, "/tsumikis/3/blocks", nil), "tsumikiID", "3")
	w := httptest.NewRecorder()

	h.GetBlocks(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetBlocks_TsumikiNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().GetTsumiki(gomock.Nil(), 999).Return(nil, nil)

	h := handler.NewTsumikiHandler(&repository.Repositories{Tsumiki: mockTsumiki}, mediamock.NewMockMediaService(ctrl))

	req := withChiParam(httptest.NewRequest(http.MethodGet, "/tsumikis/999/blocks", nil), "tsumikiID", "999")
	w := httptest.NewRecorder()

	h.GetBlocks(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: want %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestCreateTsumiki(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := sampleUser().ID
	expected := sampleTsumiki()
	thumbnail := sampleThumbnail()

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().
		CreateTsumiki(userID, expected.Title, expected.Visibility, gomock.Nil(), thumbnail.ID).
		Return(expected, nil)

	mockThumbnail := repomock.NewMockThumbnailRepository(ctrl)
	mockThumbnail.EXPECT().Get(thumbnail.ID).Return(thumbnail, nil)
	mockThumbnail.EXPECT().IsInUse(thumbnail.ID).Return(false, nil)

	mockMedia := mediamock.NewMockMediaService(ctrl)
	mockMedia.EXPECT().ResolveURL(gomock.Any()).Return("https://cdn.example.com/test").AnyTimes()

	repos := &repository.Repositories{
		Tsumiki:   mockTsumiki,
		Thumbnail: mockThumbnail,
	}
	repos.RunTxFn = func(fn repository.TxCommandFunc) error { return fn(repos) }

	h := handler.NewTsumikiHandler(repos, mockMedia)

	body := map[string]any{
		"title":        expected.Title,
		"visibility":   expected.Visibility,
		"thumbnail_id": thumbnail.ID,
	}
	req := jsonRequest(t, http.MethodPost, "/tsumikis", body)
	req = withUserID(req, userID)
	w := httptest.NewRecorder()

	h.CreateTsumiki(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}

	var resp schema.Tsumiki
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, resp.ID)
	}
}

func TestCreateTsumiki_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)

	h := handler.NewTsumikiHandler(&repository.Repositories{}, mediamock.NewMockMediaService(ctrl))

	body := map[string]any{"title": "t", "visibility": "public", "thumbnail_id": 10}
	req := jsonRequest(t, http.MethodPost, "/tsumikis", body)
	w := httptest.NewRecorder()

	h.CreateTsumiki(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestEditTsumiki(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := sampleUser().ID
	tsumiki := sampleTsumiki()
	tsumiki.User.ID = userID
	updated := sampleTsumiki()
	updated.Title = "Updated Tsumiki"

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	gomock.InOrder(
		mockTsumiki.EXPECT().GetTsumiki(gomock.Any(), tsumiki.ID).Return(tsumiki, nil),
		mockTsumiki.EXPECT().UpdateTsumiki(tsumiki.ID, updated.Title, tsumiki.Visibility, gomock.Nil()).Return(updated, nil),
	)

	mockMedia := mediamock.NewMockMediaService(ctrl)
	mockMedia.EXPECT().ResolveURL(gomock.Any()).Return("https://cdn.example.com/test").AnyTimes()

	h := handler.NewTsumikiHandler(&repository.Repositories{Tsumiki: mockTsumiki}, mockMedia)

	body := map[string]any{"title": updated.Title, "visibility": tsumiki.Visibility}
	req := jsonRequest(t, http.MethodPut, "/tsumikis/3", body)
	req = withUserID(req, userID)
	req = withChiParam(req, "tsumikiID", "3")
	w := httptest.NewRecorder()

	h.EditTsumiki(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}
}

func TestEditTsumiki_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().GetTsumiki(gomock.Any(), 999).Return(nil, nil)

	h := handler.NewTsumikiHandler(&repository.Repositories{Tsumiki: mockTsumiki}, mediamock.NewMockMediaService(ctrl))

	body := map[string]any{"title": "t", "visibility": "public"}
	req := jsonRequest(t, http.MethodPut, "/tsumikis/999", body)
	req = withUserID(req, sampleUser().ID)
	req = withChiParam(req, "tsumikiID", "999")
	w := httptest.NewRecorder()

	h.EditTsumiki(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: want %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestEditTsumiki_Forbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	tsumiki := sampleTsumiki()
	tsumiki.User.ID = 99 // different owner

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().GetTsumiki(gomock.Any(), tsumiki.ID).Return(tsumiki, nil)

	h := handler.NewTsumikiHandler(&repository.Repositories{Tsumiki: mockTsumiki}, mediamock.NewMockMediaService(ctrl))

	body := map[string]any{"title": "t", "visibility": "public"}
	req := jsonRequest(t, http.MethodPut, "/tsumikis/3", body)
	req = withUserID(req, 1) // user 1, owner is 99
	req = withChiParam(req, "tsumikiID", "3")
	w := httptest.NewRecorder()

	h.EditTsumiki(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status: want %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestUpdateTsumikiThumbnail(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := sampleUser().ID
	tsumiki := sampleTsumiki()
	tsumiki.User.ID = userID
	thumbnail := sampleThumbnail()

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().GetTsumiki(gomock.Any(), tsumiki.ID).Return(tsumiki, nil)
	mockTsumiki.EXPECT().UpdateTsumikiThumbnail(tsumiki.ID, thumbnail.ID).Return(nil)

	mockThumbnail := repomock.NewMockThumbnailRepository(ctrl)
	mockThumbnail.EXPECT().Get(thumbnail.ID).Return(thumbnail, nil)
	mockThumbnail.EXPECT().IsInUse(thumbnail.ID).Return(false, nil)

	repos := &repository.Repositories{Tsumiki: mockTsumiki, Thumbnail: mockThumbnail}
	h := handler.NewTsumikiHandler(repos, mediamock.NewMockMediaService(ctrl))

	body := map[string]any{"thumbnail_id": thumbnail.ID}
	req := jsonRequest(t, http.MethodPut, "/tsumikis/3/thumbnail", body)
	req = withUserID(req, userID)
	req = withChiParam(req, "tsumikiID", "3")
	w := httptest.NewRecorder()

	h.UpdateTsumikiThumbnail(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteTsumiki(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := sampleUser().ID
	tsumiki := sampleTsumiki()
	tsumiki.User.ID = userID

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	gomock.InOrder(
		mockTsumiki.EXPECT().GetTsumiki(gomock.Any(), tsumiki.ID).Return(tsumiki, nil),
		mockTsumiki.EXPECT().DeleteTsumiki(tsumiki.ID).Return(nil),
	)

	h := handler.NewTsumikiHandler(&repository.Repositories{Tsumiki: mockTsumiki}, mediamock.NewMockMediaService(ctrl))

	req := httptest.NewRequest(http.MethodDelete, "/tsumikis/3", nil)
	req = withUserID(req, userID)
	req = withChiParam(req, "tsumikiID", "3")
	w := httptest.NewRecorder()

	h.DeleteTsumiki(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteTsumiki_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().GetTsumiki(gomock.Any(), 999).Return(nil, nil)

	h := handler.NewTsumikiHandler(&repository.Repositories{Tsumiki: mockTsumiki}, mediamock.NewMockMediaService(ctrl))

	req := httptest.NewRequest(http.MethodDelete, "/tsumikis/999", nil)
	req = withUserID(req, sampleUser().ID)
	req = withChiParam(req, "tsumikiID", "999")
	w := httptest.NewRecorder()

	h.DeleteTsumiki(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: want %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteTsumiki_Forbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	tsumiki := sampleTsumiki()
	tsumiki.User.ID = 99

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().GetTsumiki(gomock.Any(), tsumiki.ID).Return(tsumiki, nil)

	h := handler.NewTsumikiHandler(&repository.Repositories{Tsumiki: mockTsumiki}, mediamock.NewMockMediaService(ctrl))

	req := httptest.NewRequest(http.MethodDelete, "/tsumikis/3", nil)
	req = withUserID(req, 1)
	req = withChiParam(req, "tsumikiID", "3")
	w := httptest.NewRecorder()

	h.DeleteTsumiki(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status: want %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestAddBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := sampleUser().ID
	tsumiki := sampleTsumiki()
	tsumiki.User.ID = userID
	block := sampleTsumikiBlock()

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().GetTsumiki(gomock.Any(), tsumiki.ID).Return(tsumiki, nil)

	mockBlock := repomock.NewMockTsumikiBlockRepository(ctrl)
	mockBlock.EXPECT().GetLatestBlockID(tsumiki.ID).Return(nil, nil)
	mockBlock.EXPECT().CreateBlock(tsumiki.ID, block.Message, block.Percentage, block.Condition).Return(block, nil)

	mockMedia := repomock.NewMockTsumikiBlockMediaRepository(ctrl)
	mockMedia.EXPECT().SetMediaRelation(block.ID, []int{}).Return([]schema.TsumikiBlockMedia{}, nil)

	repos := &repository.Repositories{
		Tsumiki:           mockTsumiki,
		TsumikiBlock:      mockBlock,
		TsumikiBlockMedia: mockMedia,
	}
	repos.RunTxFn = func(fn repository.TxCommandFunc) error { return fn(repos) }

	h := handler.NewTsumikiHandler(repos, mediamock.NewMockMediaService(ctrl))

	body := map[string]any{
		"message":    *block.Message,
		"percentage": block.Percentage,
		"condition":  block.Condition,
		"media_ids":  []int{},
	}
	req := jsonRequest(t, http.MethodPost, "/tsumikis/3/blocks", body)
	req = withUserID(req, userID)
	req = withChiParam(req, "tsumikiID", "3")
	w := httptest.NewRecorder()

	h.AddBlock(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAddBlock_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)

	h := handler.NewTsumikiHandler(&repository.Repositories{}, mediamock.NewMockMediaService(ctrl))

	body := map[string]any{"message": "m", "percentage": 50, "condition": 3, "media_ids": []int{}}
	req := jsonRequest(t, http.MethodPost, "/tsumikis/3/blocks", body)
	req = withChiParam(req, "tsumikiID", "3")
	w := httptest.NewRecorder()

	h.AddBlock(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAddBlock_InvalidCondition(t *testing.T) {
	ctrl := gomock.NewController(t)

	h := handler.NewTsumikiHandler(&repository.Repositories{}, mediamock.NewMockMediaService(ctrl))

	body := map[string]any{"message": "m", "percentage": 50, "condition": 6, "media_ids": []int{}}
	req := jsonRequest(t, http.MethodPost, "/tsumikis/3/blocks", body)
	req = withUserID(req, sampleUser().ID)
	req = withChiParam(req, "tsumikiID", "3")
	w := httptest.NewRecorder()

	h.AddBlock(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: want %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAddBlock_TsumikiNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().GetTsumiki(gomock.Any(), 999).Return(nil, nil)

	repos := &repository.Repositories{Tsumiki: mockTsumiki, TsumikiBlock: repomock.NewMockTsumikiBlockRepository(ctrl)}
	h := handler.NewTsumikiHandler(repos, mediamock.NewMockMediaService(ctrl))

	body := map[string]any{"message": "m", "percentage": 50, "condition": 3, "media_ids": []int{}}
	req := jsonRequest(t, http.MethodPost, "/tsumikis/999/blocks", body)
	req = withUserID(req, sampleUser().ID)
	req = withChiParam(req, "tsumikiID", "999")
	w := httptest.NewRecorder()

	h.AddBlock(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: want %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestEditBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := sampleUser().ID
	tsumiki := sampleTsumiki()
	tsumiki.User.ID = userID
	block := sampleTsumikiBlock()

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().GetTsumiki(gomock.Any(), tsumiki.ID).Return(tsumiki, nil)

	mockBlock := repomock.NewMockTsumikiBlockRepository(ctrl)
	mockBlock.EXPECT().IsBelongToTsumiki(tsumiki.ID, block.ID).Return(true, nil)
	mockBlock.EXPECT().UpdateBlock(block.ID, block.Message, block.Percentage, block.Condition).Return(block, nil)

	mockBlockMedia := repomock.NewMockTsumikiBlockMediaRepository(ctrl)
	mockBlockMedia.EXPECT().SetMediaRelation(block.ID, []int{}).Return([]schema.TsumikiBlockMedia{}, nil)

	repos := &repository.Repositories{
		Tsumiki:           mockTsumiki,
		TsumikiBlock:      mockBlock,
		TsumikiBlockMedia: mockBlockMedia,
	}
	repos.RunTxFn = func(fn repository.TxCommandFunc) error { return fn(repos) }

	h := handler.NewTsumikiHandler(repos, mediamock.NewMockMediaService(ctrl))

	body := map[string]any{
		"message":    *block.Message,
		"percentage": block.Percentage,
		"condition":  block.Condition,
		"media_ids":  []int{},
	}
	req := jsonRequest(t, http.MethodPut, "/tsumikis/3/blocks/5", body)
	req = withUserID(req, userID)
	req = withChiParams(req, "tsumikiID", "3", "blockID", "5")
	w := httptest.NewRecorder()

	h.EditBlock(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}
}

func TestEditBlock_NotBelong(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := sampleUser().ID
	tsumiki := sampleTsumiki()
	tsumiki.User.ID = userID

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().GetTsumiki(gomock.Any(), tsumiki.ID).Return(tsumiki, nil)

	mockBlock := repomock.NewMockTsumikiBlockRepository(ctrl)
	mockBlock.EXPECT().IsBelongToTsumiki(tsumiki.ID, 999).Return(false, nil)

	repos := &repository.Repositories{Tsumiki: mockTsumiki, TsumikiBlock: mockBlock}
	h := handler.NewTsumikiHandler(repos, mediamock.NewMockMediaService(ctrl))

	body := map[string]any{"message": "m", "percentage": 50, "condition": 3, "media_ids": []int{}}
	req := jsonRequest(t, http.MethodPut, "/tsumikis/3/blocks/999", body)
	req = withUserID(req, userID)
	req = withChiParams(req, "tsumikiID", "3", "blockID", "999")
	w := httptest.NewRecorder()

	h.EditBlock(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: want %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestOmitBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := sampleUser().ID
	tsumiki := sampleTsumiki()
	tsumiki.User.ID = userID
	block := sampleTsumikiBlock()

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().GetTsumiki(gomock.Any(), tsumiki.ID).Return(tsumiki, nil)

	mockBlock := repomock.NewMockTsumikiBlockRepository(ctrl)
	mockBlock.EXPECT().IsBelongToTsumiki(tsumiki.ID, block.ID).Return(true, nil)
	mockBlock.EXPECT().SoftDeleteBlock(block.ID).Return(nil)

	repos := &repository.Repositories{Tsumiki: mockTsumiki, TsumikiBlock: mockBlock}
	h := handler.NewTsumikiHandler(repos, mediamock.NewMockMediaService(ctrl))

	req := httptest.NewRequest(http.MethodDelete, "/tsumikis/3/blocks/5", nil)
	req = withUserID(req, userID)
	req = withChiParams(req, "tsumikiID", "3", "blockID", "5")
	w := httptest.NewRecorder()

	h.OmitBlock(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}
}

func TestOmitBlock_NotBelong(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := sampleUser().ID
	tsumiki := sampleTsumiki()
	tsumiki.User.ID = userID

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().GetTsumiki(gomock.Any(), tsumiki.ID).Return(tsumiki, nil)

	mockBlock := repomock.NewMockTsumikiBlockRepository(ctrl)
	mockBlock.EXPECT().IsBelongToTsumiki(tsumiki.ID, 999).Return(false, nil)

	repos := &repository.Repositories{Tsumiki: mockTsumiki, TsumikiBlock: mockBlock}
	h := handler.NewTsumikiHandler(repos, mediamock.NewMockMediaService(ctrl))

	req := httptest.NewRequest(http.MethodDelete, "/tsumikis/3/blocks/999", nil)
	req = withUserID(req, userID)
	req = withChiParams(req, "tsumikiID", "3", "blockID", "999")
	w := httptest.NewRecorder()

	h.OmitBlock(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: want %d, got %d", http.StatusBadRequest, w.Code)
	}
}
