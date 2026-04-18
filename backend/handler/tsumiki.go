package handler

import (
	"net/http"
	"tsumiki/repository"
)

type TsumikiHandler interface {
	GetMyTsumikis(w http.ResponseWriter, r *http.Request)
	GetUserTsumikis(w http.ResponseWriter, r *http.Request)
	GetTsumikis(w http.ResponseWriter, r *http.Request)
	GetSpecifiedTsumiki(w http.ResponseWriter, r *http.Request)
	CreateTsumiki(w http.ResponseWriter, r *http.Request)
	EditTsumiki(w http.ResponseWriter, r *http.Request)
	DeleteTsumiki(w http.ResponseWriter, r *http.Request)
	PostMedia(w http.ResponseWriter, r *http.Request)
	AddBlock(w http.ResponseWriter, r *http.Request)
	EditBlock(w http.ResponseWriter, r *http.Request)
	OmitBlock(w http.ResponseWriter, r *http.Request)
}

type tsumikiHandlerImpl struct {
	tsumikiRepo repository.TsumikiRepository
	blockRepo   repository.TsumikiRepository
}

func NewTsumikiHandler(tsumikiRepo repository.TsumikiRepository, blockRepo repository.TsumikiBlockRepository) TsumikiHandler {
	return &tsumikiHandlerImpl{
		tsumikiRepo: tsumikiRepo,
	}
}

func (th *tsumikiHandlerImpl) GetMyTsumikis(w http.ResponseWriter, r *http.Request)       {}
func (th *tsumikiHandlerImpl) GetUserTsumikis(w http.ResponseWriter, r *http.Request)     {}
func (th *tsumikiHandlerImpl) GetTsumikis(w http.ResponseWriter, r *http.Request)         {}
func (th *tsumikiHandlerImpl) GetSpecifiedTsumiki(w http.ResponseWriter, r *http.Request) {}
func (th *tsumikiHandlerImpl) CreateTsumiki(w http.ResponseWriter, r *http.Request)       {}
func (th *tsumikiHandlerImpl) EditTsumiki(w http.ResponseWriter, r *http.Request)         {}
func (th *tsumikiHandlerImpl) DeleteTsumiki(w http.ResponseWriter, r *http.Request)       {}
func (th *tsumikiHandlerImpl) PostMedia(w http.ResponseWriter, r *http.Request)           {}
func (th *tsumikiHandlerImpl) AddBlock(w http.ResponseWriter, r *http.Request)            {}
func (th *tsumikiHandlerImpl) EditBlock(w http.ResponseWriter, r *http.Request)           {}
func (th *tsumikiHandlerImpl) OmitBlock(w http.ResponseWriter, r *http.Request)           {}
