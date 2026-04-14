package handler

import (
	"net/http"
	"tsumiki/helper"
)

type pingResponse struct {
	Status string `json:"status"`
}

func Ping(w http.ResponseWriter, r *http.Request) {
	response := pingResponse{Status: "ok"}
	helper.ResponseOk(w, response)
}
