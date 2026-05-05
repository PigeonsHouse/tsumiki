package repository

import (
	"tsumiki/schema"
)

//go:generate mockgen -source=tsumiki_favorite.go -destination=mock/mock_tsumiki_favorite.go -package=mock

type TsumikiFavoriteRepository interface {
	SetFavoriteCount(tsumikiID int, userID int, count int) (*schema.Favorite, error)
}

type tsumikiFavoriteRepositoryImpl struct {
	db DBTX
}

func NewTsumikiFavoriteRepository(db DBTX) TsumikiFavoriteRepository {
	return &tsumikiFavoriteRepositoryImpl{db: db}
}

func (tfr *tsumikiFavoriteRepositoryImpl) SetFavoriteCount(tsumikiID int, userID int, count int) (*schema.Favorite, error) {
	var f schema.Favorite
	_, err := tfr.db.Exec(
		"INSERT INTO tsumiki_favorites (tsumiki_id, user_id, counts) VALUES (?, ?, ?) AS new_data"+
			"ON DUPLICATE KEY UPDATE counts = new_data.counts",
		tsumikiID, userID, count,
	)
	if err != nil {
		return nil, err
	}
	f.MyFavoriteCount = count

	err = tfr.db.QueryRow("SELECT counts FROM tsumiki_favorites WHERE tsumiki_id = ? AND user_id = ? LIMIT 1", tsumikiID, userID).Scan(&f.TotalFavoriteCount)
	if err != nil {
		return nil, err
	}

	return &f, nil
}
