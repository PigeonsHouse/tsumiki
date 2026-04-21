package repository

import (
	"database/sql"
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

const workSelectQuery = "SELECT w.id, w.title, w.description, w.thumbnail_url, w.created_at, w.updated_at, " +
	"u.id, u.discord_user_id, u.name, u.avatar_url, u.created_at, u.updated_at " +
	"FROM works w " +
	"JOIN users u ON w.owner_user_id = u.id"

func scanWork(scan func(...any) error) (*schema.Work, error) {
	var w schema.Work
	var thumbUrl sql.NullString
	err := scan(
		&w.ID, &w.Title, &w.Description, &thumbUrl, &w.CreatedAt, &w.UpdatedAt,
		&w.Owner.ID, &w.Owner.DiscordUserID, &w.Owner.Name, &w.Owner.AvatarUrl, &w.Owner.CreatedAt, &w.Owner.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if thumbUrl.Valid {
		w.ThumbnailUrl = &thumbUrl.String
	}
	return &w, nil
}

func (wr *workRepositoryImpl) GetWorks(pageSize, page int) ([]schema.Work, error) {
	rows, err := wr.db.Query(
		workSelectQuery+" ORDER BY w.created_at DESC LIMIT ? OFFSET ?",
		pageSize, (page-1)*pageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	works := make([]schema.Work, 0)
	for rows.Next() {
		w, err := scanWork(rows.Scan)
		if err != nil {
			return nil, err
		}
		works = append(works, *w)
	}
	return works, rows.Err()
}

func (wr *workRepositoryImpl) GetWork(workID int) (*schema.Work, error) {
	row := wr.db.QueryRow(workSelectQuery+" WHERE w.id = ?", workID)
	w, err := scanWork(row.Scan)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return w, err
}

func (wr *workRepositoryImpl) GetWorkTsumikis(workID int, pageSize, page int) ([]schema.Tsumiki, error) {
	return nil, nil
}

func (wr *workRepositoryImpl) CreateWork(userID int, title string, description string) (*schema.Work, error) {
	result, err := wr.db.Exec(
		"INSERT INTO works (owner_user_id, title, description) VALUES (?, ?, ?)",
		userID, title, description,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return wr.GetWork(int(id))
}

func (wr *workRepositoryImpl) UpdateWork(workID int, title string, description string) (*schema.Work, error) {
	_, err := wr.db.Exec(
		"UPDATE works SET title = ?, description = ? WHERE id = ?",
		title, description, workID,
	)
	if err != nil {
		return nil, err
	}
	return wr.GetWork(workID)
}

func (wr *workRepositoryImpl) DeleteWork(workID int) error {
	_, err := wr.db.Exec("DELETE FROM works WHERE id = ?", workID)
	return err
}
