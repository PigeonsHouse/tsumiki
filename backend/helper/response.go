package helper

import (
	"encoding/json"
	"net/http"
)

func ResponseJSON(w http.ResponseWriter, body any, status int) {
	resp, err := json.Marshal(body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"status":500,"detail":"application error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(resp)
}

func ResponseOk(w http.ResponseWriter, body any) {
	ResponseJSON(w, body, http.StatusOK)
}

type failureBody struct {
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func ResponseBadRequest(w http.ResponseWriter, detail string) {
	ResponseJSON(w, failureBody{Status: http.StatusBadRequest, Detail: detail}, http.StatusBadRequest)
}

func ResponseForbidden(w http.ResponseWriter, detail string) {
	ResponseJSON(w, failureBody{Status: http.StatusForbidden, Detail: detail}, http.StatusForbidden)
}

func ResponseUnauthorized(w http.ResponseWriter, detail string) {
	ResponseJSON(w, failureBody{Status: http.StatusUnauthorized, Detail: detail}, http.StatusUnauthorized)
}

func ResponseNotFound(w http.ResponseWriter, detail string) {
	ResponseJSON(w, failureBody{Status: http.StatusNotFound, Detail: detail}, http.StatusNotFound)
}

func ResponseConflict(w http.ResponseWriter, detail string) {
	ResponseJSON(w, failureBody{Status: http.StatusConflict, Detail: detail}, http.StatusConflict)
}

func ResponseInternalServerError(w http.ResponseWriter, detail string) {
	ResponseJSON(w, failureBody{Status: http.StatusInternalServerError, Detail: detail}, http.StatusInternalServerError)
}
