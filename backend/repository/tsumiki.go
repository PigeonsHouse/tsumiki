package repository

import (
	"database/sql"
	"tsumiki/schema"
)

type TsumikiRepository interface {
	GetTsumiki(watchUserID *int, tsumikiID int) (*schema.Tsumiki, error)
	GetTsumikiBlocks(tsumikiID int, pageSize, page int) ([]schema.TsumikiBlockView, error)
	GetTsumikis(watchUserID *int, pageSize, page int, authorID *int, workID *int, keyword string) ([]schema.Tsumiki, error)
	CreateTsumiki(userID int, title string, visibility string, workID *int, thumbnailID int) (*schema.Tsumiki, error)
	UpdateTsumikiThumbnail(tsumikiID int, thumbnailID int) error
	UpdateTsumiki(tsumikiID int, title string, visibility string, workID *int) (*schema.Tsumiki, error)
	DeleteTsumiki(tsumikiID int) error
	CreateMedia(tsumikiID int, mediaType string, url string) (*schema.TsumikiBlockMedia, error)
}

type tsumikiRepositoryImpl struct {
	db DBTX
}

func NewTsumikiRepository(db DBTX) TsumikiRepository {
	return &tsumikiRepositoryImpl{db: db}
}

const tsumikiSelectQuery = "SELECT t.id, t.title, t.visibility, t.created_at, t.updated_at, " +
	"u.id, u.discord_user_id, u.name, u.avatar_url, u.created_at, u.updated_at, " +
	"w.id, w.title, w.description, w.created_at, w.updated_at, " +
	"wu.id, wu.discord_user_id, wu.name, wu.avatar_url, wu.created_at, wu.updated_at, " +
	"wth.id, wth.path, wth.created_at, wth.updated_at, " +
	"tth.id, tth.path, tth.created_at, tth.updated_at " +
	"FROM tsumikis t " +
	"JOIN users u ON t.user_id = u.id " +
	"LEFT JOIN works w ON t.work_id = w.id " +
	"LEFT JOIN users wu ON w.owner_user_id = wu.id " +
	"LEFT JOIN thumbnails wth ON w.thumbnail_id = wth.id " +
	"LEFT JOIN thumbnails tth ON t.thumbnail_id = tth.id"

func scanTsumikiRow(scan func(...any) error) (*schema.Tsumiki, error) {
	var t schema.Tsumiki
	var workID sql.NullInt64
	var workTitle, workDesc sql.NullString
	var workCreatedAt, workUpdatedAt sql.NullTime
	var ownerID sql.NullInt64
	var ownerDiscordID, ownerName, ownerAvatarUrl sql.NullString
	var ownerCreatedAt, ownerUpdatedAt sql.NullTime
	var wthID sql.NullInt64
	var wthPath sql.NullString
	var wthCreatedAt, wthUpdatedAt sql.NullTime
	var tthID sql.NullInt64
	var tthPath sql.NullString
	var tthCreatedAt, tthUpdatedAt sql.NullTime

	err := scan(
		&t.ID, &t.Title, &t.Visibility, &t.CreatedAt, &t.UpdatedAt,
		&t.User.ID, &t.User.DiscordUserID, &t.User.Name, &t.User.AvatarUrl, &t.User.CreatedAt, &t.User.UpdatedAt,
		&workID, &workTitle, &workDesc, &workCreatedAt, &workUpdatedAt,
		&ownerID, &ownerDiscordID, &ownerName, &ownerAvatarUrl, &ownerCreatedAt, &ownerUpdatedAt,
		&wthID, &wthPath, &wthCreatedAt, &wthUpdatedAt,
		&tthID, &tthPath, &tthCreatedAt, &tthUpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if tthID.Valid {
		t.Thumbnail = &schema.ThumbnailUpload{
			ID:        int(tthID.Int64),
			Url:       tthPath.String,
			CreatedAt: tthCreatedAt.Time,
			UpdatedAt: tthUpdatedAt.Time,
		}
	}

	if workID.Valid {
		w := &schema.Work{
			ID:          int(workID.Int64),
			Title:       workTitle.String,
			Description: workDesc.String,
			CreatedAt:   workCreatedAt.Time,
			UpdatedAt:   workUpdatedAt.Time,
		}
		if wthID.Valid {
			w.Thumbnail = &schema.ThumbnailUpload{
				ID:        int(wthID.Int64),
				Url:       wthPath.String,
				CreatedAt: wthCreatedAt.Time,
				UpdatedAt: wthUpdatedAt.Time,
			}
		}
		if ownerID.Valid {
			w.Owner = schema.User{
				ID:            int(ownerID.Int64),
				DiscordUserID: ownerDiscordID.String,
				Name:          ownerName.String,
				AvatarUrl:     ownerAvatarUrl.String,
				CreatedAt:     ownerCreatedAt.Time,
				UpdatedAt:     ownerUpdatedAt.Time,
			}
		}
		t.Work = w
	}

	return &t, nil
}

func (tr *tsumikiRepositoryImpl) fetchTsumikiByID(tsumikiID int) (*schema.Tsumiki, error) {
	row := tr.db.QueryRow(tsumikiSelectQuery+" WHERE t.id = ?", tsumikiID)
	t, err := scanTsumikiRow(row.Scan)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return t, err
}

func (tr *tsumikiRepositoryImpl) GetTsumiki(watchUserID *int, tsumikiID int) (*schema.Tsumiki, error) {
	query := tsumikiSelectQuery + " WHERE t.id = ?"
	args := []any{tsumikiID}

	if watchUserID != nil {
		query += " AND (t.visibility = 'public' OR t.user_id = ?)"
		args = append(args, *watchUserID)
	} else {
		query += " AND t.visibility = 'public'"
	}

	row := tr.db.QueryRow(query, args...)
	t, err := scanTsumikiRow(row.Scan)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return t, err
}

func (tr *tsumikiRepositoryImpl) GetTsumikiBlocks(tsumikiID int, pageSize, page int) ([]schema.TsumikiBlockView, error) {
	// deleted_at IS NULL の条件をJOIN側に置くことで、削除済みブロックのメディアは取得しない
	rows, err := tr.db.Query(
		"SELECT b.id, b.deleted_at, b.message, b.percentage, b.condition, b.created_at, b.updated_at, "+
			"m.id, m.type, m.url, m.`order`, m.created_at, m.updated_at "+
			"FROM tsumiki_blocks b "+
			"LEFT JOIN tsumiki_block_medias m ON b.id = m.tsumiki_block_id AND b.deleted_at IS NULL "+
			"WHERE b.tsumiki_id = ? "+
			"ORDER BY b.id ASC, m.`order` "+
			"LIMIT ? OFFSET ?",
		tsumikiID, pageSize, (page-1)*pageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blockMap := make(map[int]*schema.TsumikiBlockView)
	var blockOrder []int

	for rows.Next() {
		var id int
		var deletedAt sql.NullTime
		var message sql.NullString
		var percentage, condition sql.NullInt64
		var createdAt, updatedAt sql.NullTime
		var mediaID sql.NullInt64
		var mediaType, mediaUrl sql.NullString
		var mediaOrder sql.NullInt64
		var mediaCreatedAt, mediaUpdatedAt sql.NullTime

		if err := rows.Scan(
			&id, &deletedAt, &message, &percentage, &condition, &createdAt, &updatedAt,
			&mediaID, &mediaType, &mediaUrl, &mediaOrder, &mediaCreatedAt, &mediaUpdatedAt,
		); err != nil {
			return nil, err
		}

		if _, exists := blockMap[id]; !exists {
			b := &schema.TsumikiBlockView{ID: id, IsDeleted: deletedAt.Valid}
			if !deletedAt.Valid {
				if message.Valid {
					b.Message = &message.String
				}
				b.Medias = []schema.TsumikiBlockMedia{}
				p := int(percentage.Int64)
				c := int(condition.Int64)
				b.Percentage = &p
				b.Condition = &c
				b.CreatedAt = &createdAt.Time
				b.UpdatedAt = &updatedAt.Time
			}
			blockMap[id] = b
			blockOrder = append(blockOrder, id)
		}

		if mediaID.Valid {
			blockMap[id].Medias = append(blockMap[id].Medias, schema.TsumikiBlockMedia{
				ID:        int(mediaID.Int64),
				Type:      mediaType.String,
				Url:       mediaUrl.String,
				Order:     int(mediaOrder.Int64),
				CreatedAt: mediaCreatedAt.Time,
				UpdatedAt: mediaUpdatedAt.Time,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	blocks := make([]schema.TsumikiBlockView, 0, len(blockOrder))
	for _, id := range blockOrder {
		blocks = append(blocks, *blockMap[id])
	}
	return blocks, nil
}

func (tr *tsumikiRepositoryImpl) GetTsumikis(watchUserID *int, pageSize, page int, authorID *int, workID *int, keyword string) ([]schema.Tsumiki, error) {
	query := tsumikiSelectQuery + " WHERE 1=1"
	args := []any{}

	if watchUserID != nil {
		query += " AND (t.visibility = 'public' OR t.user_id = ?)"
		args = append(args, *watchUserID)
	} else {
		query += " AND t.visibility = 'public'"
	}

	if authorID != nil {
		query += " AND t.user_id = ?"
		args = append(args, *authorID)
	}

	if workID != nil {
		query += " AND t.work_id = ?"
		args = append(args, *workID)
	}

	if keyword != "" {
		query += " AND t.title LIKE ?"
		args = append(args, "%"+keyword+"%")
	}

	query += " ORDER BY t.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := tr.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tsumikis := make([]schema.Tsumiki, 0)
	for rows.Next() {
		t, err := scanTsumikiRow(rows.Scan)
		if err != nil {
			return nil, err
		}
		tsumikis = append(tsumikis, *t)
	}
	return tsumikis, rows.Err()
}

func (tr *tsumikiRepositoryImpl) CreateTsumiki(userID int, title string, visibility string, workID *int, thumbnailID int) (*schema.Tsumiki, error) {
	result, err := tr.db.Exec(
		"INSERT INTO tsumikis (user_id, title, visibility, work_id, thumbnail_id) VALUES (?, ?, ?, ?, ?)",
		userID, title, visibility, workID, thumbnailID,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return tr.fetchTsumikiByID(int(id))
}

func (tr *tsumikiRepositoryImpl) UpdateTsumikiThumbnail(tsumikiID int, thumbnailID int) error {
	_, err := tr.db.Exec(
		"UPDATE tsumikis SET thumbnail_id = ? WHERE id = ?",
		thumbnailID, tsumikiID,
	)
	return err
}

func (tr *tsumikiRepositoryImpl) UpdateTsumiki(tsumikiID int, title string, visibility string, workID *int) (*schema.Tsumiki, error) {
	_, err := tr.db.Exec(
		"UPDATE tsumikis SET title = ?, visibility = ?, work_id = ? WHERE id = ?",
		title, visibility, workID, tsumikiID,
	)
	if err != nil {
		return nil, err
	}
	return tr.fetchTsumikiByID(tsumikiID)
}

func (tr *tsumikiRepositoryImpl) DeleteTsumiki(tsumikiID int) error {
	_, err := tr.db.Exec("DELETE FROM tsumikis WHERE id = ?", tsumikiID)
	return err
}

func (tr *tsumikiRepositoryImpl) CreateMedia(tsumikiID int, mediaType string, url string) (*schema.TsumikiBlockMedia, error) {
	result, err := tr.db.Exec(
		"INSERT INTO tsumiki_block_medias (type, url) VALUES (?, ?)",
		mediaType, url,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	var m schema.TsumikiBlockMedia
	err = tr.db.QueryRow(
		"SELECT id, type, url, created_at, updated_at FROM tsumiki_block_medias WHERE id = ?",
		id,
	).Scan(&m.ID, &m.Type, &m.Url, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
