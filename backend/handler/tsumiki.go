package handler

import "tsumiki/repository"

type TsumikiHandler interface {
}

type tsumikiHandlerImpl struct {
	repository repository.TsumikiRepository
}

func NewTsumikiHandler(tsumikiRepo repository.TsumikiRepository) TsumikiHandler {
	return &tsumikiHandlerImpl{
		repository: tsumikiRepo,
	}
}
