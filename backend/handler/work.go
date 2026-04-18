package handler

import (
	"net/http"
	"tsumiki/repository"
)

type WorkHandler interface {
	GetWorks(w http.ResponseWriter, r *http.Request)
	GetSpecifiedWork(w http.ResponseWriter, r *http.Request)
	GetWorkTsumiki(w http.ResponseWriter, r *http.Request)
	CreateWork(w http.ResponseWriter, r *http.Request)
	EditWork(w http.ResponseWriter, r *http.Request)
	DeleteWork(w http.ResponseWriter, r *http.Request)
}

type workHandlerImpl struct {
	repository repository.WorkRepository
}

func NewWorkHandler(WorkRepo repository.WorkRepository) WorkHandler {
	return &workHandlerImpl{
		repository: WorkRepo,
	}
}

func (wh *workHandlerImpl) GetWorks(w http.ResponseWriter, r *http.Request)         {}
func (wh *workHandlerImpl) GetSpecifiedWork(w http.ResponseWriter, r *http.Request) {}
func (wh *workHandlerImpl) GetWorkTsumiki(w http.ResponseWriter, r *http.Request)   {}
func (wh *workHandlerImpl) CreateWork(w http.ResponseWriter, r *http.Request)       {}
func (wh *workHandlerImpl) EditWork(w http.ResponseWriter, r *http.Request)         {}
func (wh *workHandlerImpl) DeleteWork(w http.ResponseWriter, r *http.Request)       {}
