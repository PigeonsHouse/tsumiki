package repository

import (
	"database/sql"
	"tsumiki/schema"
)

type TsumikiBlockRepository interface {
	IsBelongToTsumiki(tsumikiID int, blockID int) (bool, error)
	GetLatestBlockID(tsumikiID int) (*int, error)
	CreateBlock(tsumikiID int, message *string, percentage int, condition int) (*schema.TsumikiBlock, error)
	UpdateBlock(blockID int, message *string, percentage int, condition int) (*schema.TsumikiBlock, error)
	SoftDeleteBlock(blockID int) error
}

type tsumikiBlockRepositoryImpl struct {
	db DBTX
}

func NewTsumikiBlockRepository(db DBTX) TsumikiBlockRepository {
	return &tsumikiBlockRepositoryImpl{db: db}
}

func (tbr *tsumikiBlockRepositoryImpl) fetchBlock(blockID int) (*schema.TsumikiBlock, error) {
	var b schema.TsumikiBlock
	err := tbr.db.QueryRow(
		"SELECT id, message, percentage, `condition`, next_block_id, tsumiki_id, created_at, updated_at "+
			"FROM tsumiki_blocks WHERE id = ? AND deleted_at IS NULL",
		blockID,
	).Scan(&b.ID, &b.Message, &b.Percentage, &b.Condition, &b.NextBlockId, &b.TsumikiId, &b.CreatedAt, &b.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	b.Medias = []schema.TsumikiBlockMedia{}
	return &b, nil
}

func (tbr *tsumikiBlockRepositoryImpl) GetLatestBlockID(tsumikiID int) (*int, error) {
	return tbr.fetchTailBlockID(tsumikiID)
}

func (tbr *tsumikiBlockRepositoryImpl) IsBelongToTsumiki(tsumikiID int, blockID int) (bool, error) {
	var count int
	err := tbr.db.QueryRow(
		"SELECT COUNT(*) FROM tsumiki_blocks WHERE id = ? AND tsumiki_id = ? AND deleted_at IS NULL",
		blockID, tsumikiID,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (tbr *tsumikiBlockRepositoryImpl) fetchTailBlockID(tsumikiID int) (*int, error) {
	var id int
	err := tbr.db.QueryRow(
		"SELECT id FROM tsumiki_blocks WHERE tsumiki_id = ? AND next_block_id IS NULL",
		tsumikiID,
	).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (tbr *tsumikiBlockRepositoryImpl) CreateBlock(tsumikiID int, message *string, percentage int, condition int) (*schema.TsumikiBlock, error) {
	prevTailID, err := tbr.fetchTailBlockID(tsumikiID)
	if err != nil {
		return nil, err
	}

	result, err := tbr.db.Exec(
		"INSERT INTO tsumiki_blocks (tsumiki_id, message, percentage, `condition`) VALUES (?, ?, ?, ?)",
		tsumikiID, message, percentage, condition,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if prevTailID != nil {
		_, err = tbr.db.Exec(
			"UPDATE tsumiki_blocks SET next_block_id = ? WHERE id = ?",
			id, *prevTailID,
		)
		if err != nil {
			return nil, err
		}
	}

	return tbr.fetchBlock(int(id))
}

func (tbr *tsumikiBlockRepositoryImpl) UpdateBlock(blockID int, message *string, percentage int, condition int) (*schema.TsumikiBlock, error) {
	_, err := tbr.db.Exec(
		"UPDATE tsumiki_blocks SET message = ?, percentage = ?, `condition` = ? WHERE id = ? AND deleted_at IS NULL",
		message, percentage, condition, blockID,
	)
	if err != nil {
		return nil, err
	}
	return tbr.fetchBlock(blockID)
}

func (tbr *tsumikiBlockRepositoryImpl) SoftDeleteBlock(blockID int) error {
	_, err := tbr.db.Exec(
		"UPDATE tsumiki_blocks SET deleted_at = NOW() WHERE id = ?",
		blockID,
	)
	return err
}
