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

func TestGetWorks(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleWork()

	mockWork := repomock.NewMockWorkRepository(ctrl)
	mockWork.EXPECT().GetWorks(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]schema.Work{*expected}, nil)

	mockMedia := mediamock.NewMockMediaService(ctrl)
	mockMedia.EXPECT().ResolveURL(gomock.Any()).Return("https://cdn.example.com/test").AnyTimes()

	repos := &repository.Repositories{Work: mockWork}
	h := handler.NewWorkHandler(repos, mockMedia)

	req := httptest.NewRequest(http.MethodGet, "/works", nil)
	w := httptest.NewRecorder()

	h.GetWorks(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}

	var resp []schema.Work
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("want 1 work, got %d", len(resp))
	}
	if resp[0].ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, resp[0].ID)
	}
}

func TestGetSpecifiedWork(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleWork()

	mockWork := repomock.NewMockWorkRepository(ctrl)
	mockWork.EXPECT().GetWork(expected.ID).Return(expected, nil)

	mockMedia := mediamock.NewMockMediaService(ctrl)
	mockMedia.EXPECT().ResolveURL(gomock.Any()).Return("https://cdn.example.com/test").AnyTimes()

	repos := &repository.Repositories{Work: mockWork}
	h := handler.NewWorkHandler(repos, mockMedia)

	req := withChiParam(httptest.NewRequest(http.MethodGet, "/works/2", nil), "workID", "2")
	w := httptest.NewRecorder()

	h.GetSpecifiedWork(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}

	var resp schema.Work
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, resp.ID)
	}
}

func TestGetSpecifiedWork_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)

	h := handler.NewWorkHandler(&repository.Repositories{}, mediamock.NewMockMediaService(ctrl))

	req := withChiParam(httptest.NewRequest(http.MethodGet, "/works/abc", nil), "workID", "abc")
	w := httptest.NewRecorder()

	h.GetSpecifiedWork(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: want %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetSpecifiedWork_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockWork := repomock.NewMockWorkRepository(ctrl)
	mockWork.EXPECT().GetWork(999).Return(nil, nil)

	h := handler.NewWorkHandler(&repository.Repositories{Work: mockWork}, mediamock.NewMockMediaService(ctrl))

	req := withChiParam(httptest.NewRequest(http.MethodGet, "/works/999", nil), "workID", "999")
	w := httptest.NewRecorder()

	h.GetSpecifiedWork(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: want %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetSpecifiedWork_Forbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	work := sampleWork()
	work.Visibility = "limited"

	mockWork := repomock.NewMockWorkRepository(ctrl)
	mockWork.EXPECT().GetWork(work.ID).Return(work, nil)

	// user_id なし → canAccessWork が 403 を返す
	h := handler.NewWorkHandler(&repository.Repositories{Work: mockWork}, mediamock.NewMockMediaService(ctrl))

	req := withChiParam(httptest.NewRequest(http.MethodGet, "/works/2", nil), "workID", "2")
	w := httptest.NewRecorder()

	h.GetSpecifiedWork(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status: want %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestGetWorkTsumiki(t *testing.T) {
	ctrl := gomock.NewController(t)
	work := sampleWork()
	tsumiki := sampleTsumiki()

	mockWork := repomock.NewMockWorkRepository(ctrl)
	mockWork.EXPECT().GetWork(work.ID).Return(work, nil)

	mockTsumiki := repomock.NewMockTsumikiRepository(ctrl)
	mockTsumiki.EXPECT().GetTsumikis(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]schema.Tsumiki{*tsumiki}, nil)

	mockMedia := mediamock.NewMockMediaService(ctrl)
	mockMedia.EXPECT().ResolveURL(gomock.Any()).Return("https://cdn.example.com/test").AnyTimes()

	repos := &repository.Repositories{Work: mockWork, Tsumiki: mockTsumiki}
	h := handler.NewWorkHandler(repos, mockMedia)

	req := withChiParam(httptest.NewRequest(http.MethodGet, "/works/2/tsumikis", nil), "workID", "2")
	w := httptest.NewRecorder()

	h.GetWorkTsumiki(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetWorkTsumiki_WorkNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockWork := repomock.NewMockWorkRepository(ctrl)
	mockWork.EXPECT().GetWork(999).Return(nil, nil)

	h := handler.NewWorkHandler(&repository.Repositories{Work: mockWork}, mediamock.NewMockMediaService(ctrl))

	req := withChiParam(httptest.NewRequest(http.MethodGet, "/works/999/tsumikis", nil), "workID", "999")
	w := httptest.NewRecorder()

	h.GetWorkTsumiki(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: want %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestCreateWork(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleWork()

	mockWork := repomock.NewMockWorkRepository(ctrl)
	mockWork.EXPECT().
		CreateWork(sampleUser().ID, expected.Title, expected.Visibility, expected.Description, nil).
		Return(expected, nil)

	mockMedia := mediamock.NewMockMediaService(ctrl)
	mockMedia.EXPECT().ResolveURL(gomock.Any()).Return("https://cdn.example.com/test").AnyTimes()

	repos := &repository.Repositories{Work: mockWork}
	h := handler.NewWorkHandler(repos, mockMedia)

	body := map[string]any{
		"title":       expected.Title,
		"visibility":  expected.Visibility,
		"description": expected.Description,
	}
	req := jsonRequest(t, http.MethodPost, "/works", body)
	req = withUserID(req, sampleUser().ID)
	w := httptest.NewRecorder()

	h.CreateWork(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}

	var resp schema.Work
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Title != expected.Title {
		t.Errorf("Title: want %s, got %s", expected.Title, resp.Title)
	}
}

func TestCreateWork_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)

	h := handler.NewWorkHandler(&repository.Repositories{}, mediamock.NewMockMediaService(ctrl))

	body := map[string]any{"title": "Test", "visibility": "public", "description": "desc"}
	req := jsonRequest(t, http.MethodPost, "/works", body)
	w := httptest.NewRecorder()

	h.CreateWork(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestEditWork(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := sampleUser().ID
	work := sampleWork()
	work.Owner.ID = userID
	updated := sampleWork()
	updated.Title = "Updated Title"

	mockWork := repomock.NewMockWorkRepository(ctrl)
	gomock.InOrder(
		mockWork.EXPECT().GetWork(work.ID).Return(work, nil),
		mockWork.EXPECT().UpdateWork(work.ID, updated.Title, work.Visibility, work.Description).Return(updated, nil),
	)

	mockMedia := mediamock.NewMockMediaService(ctrl)
	mockMedia.EXPECT().ResolveURL(gomock.Any()).Return("https://cdn.example.com/test").AnyTimes()

	repos := &repository.Repositories{Work: mockWork}
	h := handler.NewWorkHandler(repos, mockMedia)

	body := map[string]any{"title": updated.Title, "visibility": work.Visibility, "description": work.Description}
	req := jsonRequest(t, http.MethodPut, "/works/2", body)
	req = withUserID(req, userID)
	req = withChiParam(req, "workID", "2")
	w := httptest.NewRecorder()

	h.EditWork(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}
}

func TestEditWork_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)

	h := handler.NewWorkHandler(&repository.Repositories{}, mediamock.NewMockMediaService(ctrl))

	body := map[string]any{"title": "t", "visibility": "public", "description": "d"}
	req := jsonRequest(t, http.MethodPut, "/works/2", body)
	req = withChiParam(req, "workID", "2")
	w := httptest.NewRecorder()

	h.EditWork(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestEditWork_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockWork := repomock.NewMockWorkRepository(ctrl)
	mockWork.EXPECT().GetWork(999).Return(nil, nil)

	h := handler.NewWorkHandler(&repository.Repositories{Work: mockWork}, mediamock.NewMockMediaService(ctrl))

	body := map[string]any{"title": "t", "visibility": "public", "description": "d"}
	req := jsonRequest(t, http.MethodPut, "/works/999", body)
	req = withUserID(req, sampleUser().ID)
	req = withChiParam(req, "workID", "999")
	w := httptest.NewRecorder()

	h.EditWork(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: want %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestEditWork_Forbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	work := sampleWork()
	work.Owner.ID = 99 // different from request user

	mockWork := repomock.NewMockWorkRepository(ctrl)
	mockWork.EXPECT().GetWork(work.ID).Return(work, nil)

	h := handler.NewWorkHandler(&repository.Repositories{Work: mockWork}, mediamock.NewMockMediaService(ctrl))

	body := map[string]any{"title": "t", "visibility": "public", "description": "d"}
	req := jsonRequest(t, http.MethodPut, "/works/2", body)
	req = withUserID(req, 1) // user 1, owner is 99
	req = withChiParam(req, "workID", "2")
	w := httptest.NewRecorder()

	h.EditWork(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status: want %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestUpdateWorkThumbnail(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := sampleUser().ID
	work := sampleWork()
	work.Owner.ID = userID
	thumbnail := sampleThumbnail()

	mockWork := repomock.NewMockWorkRepository(ctrl)
	mockWork.EXPECT().GetWork(work.ID).Return(work, nil)
	mockWork.EXPECT().UpdateWorkThumbnail(work.ID, thumbnail.ID).Return(nil)

	mockThumbnail := repomock.NewMockThumbnailRepository(ctrl)
	mockThumbnail.EXPECT().Get(thumbnail.ID).Return(thumbnail, nil)
	mockThumbnail.EXPECT().IsInUse(thumbnail.ID).Return(false, nil)

	repos := &repository.Repositories{Work: mockWork, Thumbnail: mockThumbnail}
	h := handler.NewWorkHandler(repos, mediamock.NewMockMediaService(ctrl))

	body := map[string]any{"thumbnail_id": thumbnail.ID}
	req := jsonRequest(t, http.MethodPut, "/works/2/thumbnail", body)
	req = withUserID(req, userID)
	req = withChiParam(req, "workID", "2")
	w := httptest.NewRecorder()

	h.UpdateWorkThumbnail(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteWork(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := sampleUser().ID
	work := sampleWork()
	work.Owner.ID = userID

	mockWork := repomock.NewMockWorkRepository(ctrl)
	gomock.InOrder(
		mockWork.EXPECT().GetWork(work.ID).Return(work, nil),
		mockWork.EXPECT().DeleteWork(work.ID).Return(nil),
	)

	h := handler.NewWorkHandler(&repository.Repositories{Work: mockWork}, mediamock.NewMockMediaService(ctrl))

	req := httptest.NewRequest(http.MethodDelete, "/works/2", nil)
	req = withUserID(req, userID)
	req = withChiParam(req, "workID", "2")
	w := httptest.NewRecorder()

	h.DeleteWork(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteWork_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockWork := repomock.NewMockWorkRepository(ctrl)
	mockWork.EXPECT().GetWork(999).Return(nil, nil)

	h := handler.NewWorkHandler(&repository.Repositories{Work: mockWork}, mediamock.NewMockMediaService(ctrl))

	req := httptest.NewRequest(http.MethodDelete, "/works/999", nil)
	req = withUserID(req, sampleUser().ID)
	req = withChiParam(req, "workID", "999")
	w := httptest.NewRecorder()

	h.DeleteWork(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: want %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteWork_Forbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	work := sampleWork()
	work.Owner.ID = 99

	mockWork := repomock.NewMockWorkRepository(ctrl)
	mockWork.EXPECT().GetWork(work.ID).Return(work, nil)

	h := handler.NewWorkHandler(&repository.Repositories{Work: mockWork}, mediamock.NewMockMediaService(ctrl))

	req := httptest.NewRequest(http.MethodDelete, "/works/2", nil)
	req = withUserID(req, 1)
	req = withChiParam(req, "workID", "2")
	w := httptest.NewRecorder()

	h.DeleteWork(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status: want %d, got %d", http.StatusForbidden, w.Code)
	}
}
