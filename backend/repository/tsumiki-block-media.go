package repository

//go:generate mockgen -source=tsumiki-block-media.go -destination=mock/mock_tsumiki_block_media.go -package=mock

import (
	"strings"
	"tsumiki/schema"
)

type TsumikiBlockMediaRepository interface {
	CreateMedia(url string, mediaType string) (*schema.TsumikiBlockMedia, error)
	SetMediaRelation(blockID int, updatedMediaIDs []int) ([]schema.TsumikiBlockMedia, error)
}

type tsumikiBlockMediaRepositoryImpl struct {
	db DBTX
}

func NewTsumikiBlockMediaRepository(db DBTX) TsumikiBlockMediaRepository {
	return &tsumikiBlockMediaRepositoryImpl{db: db}
}

func (tbmr *tsumikiBlockMediaRepositoryImpl) CreateMedia(url string, mediaType string) (*schema.TsumikiBlockMedia, error) {
	result, err := tbmr.db.Exec(
		"INSERT INTO tsumiki_block_medias (url, type) VALUES (?, ?)",
		url, mediaType,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	var m schema.TsumikiBlockMedia
	err = tbmr.db.QueryRow(
		"SELECT id, type, url, created_at, updated_at FROM tsumiki_block_medias WHERE id = ?",
		id,
	).Scan(&m.ID, &m.Type, &m.Url, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (tbmr *tsumikiBlockMediaRepositoryImpl) SetMediaRelation(blockID int, updatedMediaIDs []int) ([]schema.TsumikiBlockMedia, error) {
	// tsumiki_block_medias は (tsumiki_block_id, order) にユニーク制約を持つ。
	// 新しい紐付けを設定する前に既存の紐付けを解除しないと制約違反になるため、先に NULL にリセットする。
	_, err := tbmr.db.Exec(
		"UPDATE tsumiki_block_medias SET tsumiki_block_id = NULL, `order` = NULL WHERE tsumiki_block_id = ?",
		blockID,
	)
	if err != nil {
		return nil, err
	}

	if len(updatedMediaIDs) > 0 {
		caseParts := strings.Repeat(" WHEN ? THEN ?", len(updatedMediaIDs))
		inPlaceholders := strings.Repeat(",?", len(updatedMediaIDs))[1:]

		args := make([]any, 0, 1+len(updatedMediaIDs)*3)
		args = append(args, blockID)
		for i, mediaID := range updatedMediaIDs {
			args = append(args, mediaID, i)
		}
		for _, mediaID := range updatedMediaIDs {
			args = append(args, mediaID)
		}

		_, err = tbmr.db.Exec(
			"UPDATE tsumiki_block_medias SET tsumiki_block_id = ?, `order` = CASE id"+caseParts+" END WHERE id IN ("+inPlaceholders+")",
			args...,
		)
		if err != nil {
			return nil, err
		}
	}

	rows, err := tbmr.db.Query(
		"SELECT id, type, url, `order`, created_at, updated_at FROM tsumiki_block_medias "+
			"WHERE tsumiki_block_id = ? ORDER BY `order`",
		blockID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	medias := make([]schema.TsumikiBlockMedia, 0)
	for rows.Next() {
		var m schema.TsumikiBlockMedia
		if err := rows.Scan(&m.ID, &m.Type, &m.Url, &m.Order, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		medias = append(medias, m)
	}
	return medias, rows.Err()
}
