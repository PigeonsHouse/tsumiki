package handler

import (
	"net/http"
	"tsumiki/helper"
	"tsumiki/schema"
)

type PingHandler interface {
	Ping(w http.ResponseWriter, r *http.Request)
}

type pingHandlerImpl struct {
}

func NewPingHandler() PingHandler {
	return &pingHandlerImpl{}
}

func (ph *pingHandlerImpl) Ping(w http.ResponseWriter, r *http.Request) {
	response := schema.PingResponse{Status: "ok"}
	helper.ResponseOk(w, response)
}
