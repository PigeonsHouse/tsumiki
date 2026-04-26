package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"tsumiki/handler"
	"tsumiki/schema"
)

func TestPing(t *testing.T) {
	h := handler.NewPingHandler()

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	h.Ping(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}

	var resp schema.PingResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("Status: want ok, got %s", resp.Status)
	}
}
