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

func TestGetMyInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleUser()

	mockRepo := repomock.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().FindByID(expected.ID).Return(expected, nil)

	mockMedia := mediamock.NewMockMediaService(ctrl)
	mockMedia.EXPECT().ResolveURL(expected.AvatarUrl).Return("https://cdn.example.com/" + expected.AvatarUrl)

	h := handler.NewUserHandler(&repository.Repositories{User: mockRepo}, mockMedia)

	req := withUserID(httptest.NewRequest(http.MethodGet, "/users/me", nil), expected.ID)
	w := httptest.NewRecorder()

	h.GetMyInfo(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}

	var resp schema.User
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, resp.ID)
	}
	if resp.Name != expected.Name {
		t.Errorf("Name: want %s, got %s", expected.Name, resp.Name)
	}
	if resp.AvatarUrl != "https://cdn.example.com/"+expected.AvatarUrl {
		t.Errorf("AvatarUrl: want resolved URL, got %s", resp.AvatarUrl)
	}
}

func TestGetMyInfo_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)

	h := handler.NewUserHandler(&repository.Repositories{}, mediamock.NewMockMediaService(ctrl))

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	w := httptest.NewRecorder()

	h.GetMyInfo(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestGetMyInfo_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := repomock.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().FindByID(999).Return(nil, nil)

	h := handler.NewUserHandler(&repository.Repositories{User: mockRepo}, mediamock.NewMockMediaService(ctrl))

	req := withUserID(httptest.NewRequest(http.MethodGet, "/users/me", nil), 999)
	w := httptest.NewRecorder()

	h.GetMyInfo(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: want %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetUserInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleUser()

	mockRepo := repomock.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().FindByID(expected.ID).Return(expected, nil)

	mockMedia := mediamock.NewMockMediaService(ctrl)
	mockMedia.EXPECT().ResolveURL(expected.AvatarUrl).Return("https://cdn.example.com/" + expected.AvatarUrl)

	h := handler.NewUserHandler(&repository.Repositories{User: mockRepo}, mockMedia)

	req := withChiParam(httptest.NewRequest(http.MethodGet, "/users/1", nil), "userID", "1")
	w := httptest.NewRecorder()

	h.GetUserInfo(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}

	var resp schema.User
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, resp.ID)
	}
	if resp.Name != expected.Name {
		t.Errorf("Name: want %s, got %s", expected.Name, resp.Name)
	}
	if resp.AvatarUrl != "https://cdn.example.com/"+expected.AvatarUrl {
		t.Errorf("AvatarUrl: want resolved URL, got %s", resp.AvatarUrl)
	}
}

func TestGetUserInfo_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)

	h := handler.NewUserHandler(&repository.Repositories{}, mediamock.NewMockMediaService(ctrl))

	req := withChiParam(httptest.NewRequest(http.MethodGet, "/users/abc", nil), "userID", "abc")
	w := httptest.NewRecorder()

	h.GetUserInfo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: want %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetUserInfo_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := repomock.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().FindByID(999).Return(nil, nil)

	h := handler.NewUserHandler(&repository.Repositories{User: mockRepo}, mediamock.NewMockMediaService(ctrl))

	req := withChiParam(httptest.NewRequest(http.MethodGet, "/users/999", nil), "userID", "999")
	w := httptest.NewRecorder()

	h.GetUserInfo(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: want %d, got %d", http.StatusNotFound, w.Code)
	}
}
