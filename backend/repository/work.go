package repository

import (
	"database/sql"
	"tsumiki/schema"
)

type WorkRepository interface {
	GetWorks(pageSize, page int) ([]schema.Work, error)
	GetWork(workID int) (*schema.Work, error)
	CreateWork(userID int, title string, description string, thumbnailID *int) (*schema.Work, error)
	UpdateWorkThumbnail(workID int, thumbnailID int) error
	UpdateWork(workID int, title string, description string) (*schema.Work, error)
	DeleteWork(workID int) error
}

type workRepositoryImpl struct {
	db DBTX
}

func NewWorkRepository(db DBTX) WorkRepository {
	return &workRepositoryImpl{db: db}
}

const workSelectQuery = "SELECT w.id, w.title, w.description, w.created_at, w.updated_at, " +
	"u.id, u.discord_user_id, u.name, u.avatar_url, u.created_at, u.updated_at, " +
	"th.id, th.path, th.created_at, th.updated_at " +
	"FROM works w " +
	"JOIN users u ON w.owner_user_id = u.id " +
	"LEFT JOIN thumbnails th ON w.thumbnail_id = th.id"

func scanWork(scan func(...any) error) (*schema.Work, error) {
	var w schema.Work
	var thID sql.NullInt64
	var thPath sql.NullString
	var thCreatedAt, thUpdatedAt sql.NullTime
	err := scan(
		&w.ID, &w.Title, &w.Description, &w.CreatedAt, &w.UpdatedAt,
		&w.Owner.ID, &w.Owner.DiscordUserID, &w.Owner.Name, &w.Owner.AvatarUrl, &w.Owner.CreatedAt, &w.Owner.UpdatedAt,
		&thID, &thPath, &thCreatedAt, &thUpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if thID.Valid {
		w.Thumbnail = &schema.ThumbnailUpload{
			ID:        int(thID.Int64),
			Url:       thPath.String,
			CreatedAt: thCreatedAt.Time,
			UpdatedAt: thUpdatedAt.Time,
		}
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

func (wr *workRepositoryImpl) CreateWork(userID int, title string, description string, thumbnailID *int) (*schema.Work, error) {
	result, err := wr.db.Exec(
		"INSERT INTO works (owner_user_id, title, description, thumbnail_id) VALUES (?, ?, ?, ?)",
		userID, title, description, thumbnailID,
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

func (wr *workRepositoryImpl) UpdateWorkThumbnail(workID int, thumbnailID int) error {
	_, err := wr.db.Exec(
		"UPDATE works SET thumbnail_id = ? WHERE id = ?",
		thumbnailID, workID,
	)
	return err
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
