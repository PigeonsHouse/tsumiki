package repository

import (
	"tsumiki/schema"
)

type WorkRepository interface {
	GetWorks(pageSize, page int) ([]schema.Work, error)
	GetWork(workID int) (*schema.Work, error)
	CreateWork(userID int, title string, description string) (*schema.Work, error)
	UpdateWork(workID int, title string, description string) (*schema.Work, error)
	DeleteWork(workID int) error
}

type workRepositoryImpl struct {
	db DBTX
}

func NewWorkRepository(db DBTX) WorkRepository {
	return &workRepositoryImpl{db: db}
}

func (wr *workRepositoryImpl) GetWorks(pageSize, page int) ([]schema.Work, error) {
	return nil, nil
}
func (wr *workRepositoryImpl) GetWork(workID int) (*schema.Work, error) {
	return nil, nil
}
func (wr *workRepositoryImpl) GetWorkTsumikis(workID int, pageSize, page int) ([]schema.Tsumiki, error) {
	return nil, nil
}
func (wr *workRepositoryImpl) CreateWork(userID int, title string, description string) (*schema.Work, error) {
	return nil, nil
}
func (wr *workRepositoryImpl) UpdateWork(workID int, title string, description string) (*schema.Work, error) {
	return nil, nil
}
func (wr *workRepositoryImpl) DeleteWork(workID int) error {
	return nil
}
