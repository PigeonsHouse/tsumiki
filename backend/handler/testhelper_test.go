package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
	"time"
	"tsumiki/schema"

	"github.com/go-chi/chi/v5"
)

// --- sample data ---

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

func sampleWork() *schema.Work {
	return &schema.Work{
		ID:          2,
		Title:       "Test Work",
		Description: "Test Description",
		Visibility:  "public",
		Owner:       *sampleUser(),
		CreatedAt:   time.Now().Truncate(time.Second),
		UpdatedAt:   time.Now().Truncate(time.Second),
	}
}

func sampleTsumiki() *schema.Tsumiki {
	return &schema.Tsumiki{
		ID:         3,
		Title:      "Test Tsumiki",
		Visibility: "public",
		User:       *sampleUser(),
		CreatedAt:  time.Now().Truncate(time.Second),
		UpdatedAt:  time.Now().Truncate(time.Second),
	}
}

func sampleTsumikiBlock() *schema.TsumikiBlock {
	msg := "test message"
	return &schema.TsumikiBlock{
		ID:         5,
		Message:    &msg,
		Medias:     []schema.TsumikiBlockMedia{},
		Percentage: 50,
		Condition:  3,
		TsumikiId:  3,
		CreatedAt:  time.Now().Truncate(time.Second),
		UpdatedAt:  time.Now().Truncate(time.Second),
	}
}

func sampleThumbnail() *schema.ThumbnailUpload {
	return &schema.ThumbnailUpload{
		ID:        10,
		Url:       "thumbnails/10.png",
		CreatedAt: time.Now().Truncate(time.Second),
		UpdatedAt: time.Now().Truncate(time.Second),
	}
}

func sampleBlockMedia() *schema.TsumikiBlockMedia {
	return &schema.TsumikiBlockMedia{
		ID:        7,
		Type:      "image",
		Url:       "media/7.jpg",
		CreatedAt: time.Now().Truncate(time.Second),
		UpdatedAt: time.Now().Truncate(time.Second),
	}
}

// --- request helpers ---

func withUserID(r *http.Request, id int) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "user_id", id))
}

// withChiParam sets a chi URL parameter on the request, preserving existing context values.
func withChiParam(r *http.Request, key, val string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, val)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

// withChiParams sets multiple chi URL parameters (key, value pairs).
func withChiParams(r *http.Request, pairs ...string) *http.Request {
	rctx := chi.NewRouteContext()
	for i := 0; i+1 < len(pairs); i += 2 {
		rctx.URLParams.Add(pairs[i], pairs[i+1])
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func jsonRequest(t *testing.T, method, path string, v any) *http.Request {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// --- multipart helpers ---

func createMinimalPNG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func newMultipartRequest(t *testing.T, fieldName, filename, contentType string, data []byte) *http.Request {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, filename))
	h.Set("Content-Type", contentType)

	fw, err := mw.CreatePart(h)
	if err != nil {
		t.Fatalf("CreatePart: %v", err)
	}
	if _, err := fw.Write(data); err != nil {
		t.Fatalf("write multipart data: %v", err)
	}
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/thumbnails", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return req
}
