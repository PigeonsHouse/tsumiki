package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"tsumiki/handler"
	mediamock "tsumiki/media/mock"
	"tsumiki/repository"
	repomock "tsumiki/repository/mock"
	"tsumiki/schema"

	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"
)

func sampleUser() *schema.User {
	guildID := "guild123"
	return &schema.User{
		ID:            1,
		DiscordUserID: "discord123",
		Name:          "Test User",
		GuildID:       &guildID,
		AvatarUrl:     "avatars/1/abc.png",
		CreatedAt:     time.Now().Truncate(time.Second),
		UpdatedAt:     time.Now().Truncate(time.Second),
	}
}

func TestGetMyInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleUser()

	mockRepo := repomock.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().FindByID(expected.ID).Return(expected, nil)

	mockMedia := mediamock.NewMockMediaService(ctrl)
	mockMedia.EXPECT().ResolveURL(expected.AvatarUrl).Return("https://cdn.example.com/" + expected.AvatarUrl)

	h := handler.NewUserHandler(&repository.Repositories{User: mockRepo}, mockMedia)

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", expected.ID))
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

func TestGetUserInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleUser()

	mockRepo := repomock.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().FindByID(expected.ID).Return(expected, nil)

	mockMedia := mediamock.NewMockMediaService(ctrl)
	mockMedia.EXPECT().ResolveURL(expected.AvatarUrl).Return("https://cdn.example.com/" + expected.AvatarUrl)

	h := handler.NewUserHandler(&repository.Repositories{User: mockRepo}, mockMedia)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", "1")
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
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
