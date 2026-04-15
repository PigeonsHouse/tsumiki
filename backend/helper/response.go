package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ResponseJSON(w http.ResponseWriter, body any, status int) error {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	resp, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	w.Write(resp)

	return nil
}

func ResponseOk(w http.ResponseWriter, body any) error {
	return ResponseJSON(w, body, http.StatusOK)
}

type failureBody struct {
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func ResponseBadRequest(w http.ResponseWriter, detail string) error {
	return ResponseJSON(w, failureBody{Status: http.StatusBadRequest, Detail: detail}, http.StatusBadRequest)
}

func ResponseForbidden(w http.ResponseWriter, detail string) error {
	return ResponseJSON(w, failureBody{Status: http.StatusForbidden, Detail: detail}, http.StatusForbidden)
}

func ResponseInternalServerError(w http.ResponseWriter, detail string) error {
	return ResponseJSON(w, failureBody{Status: http.StatusBadRequest, Detail: detail}, http.StatusInternalServerError)
}
