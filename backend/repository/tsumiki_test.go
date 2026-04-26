package repository_test

import (
	"database/sql"
	"testing"
	"time"
	"tsumiki/repository"
	"tsumiki/repository/mock"

	"go.uber.org/mock/gomock"
)

func TestTsumikiRepository_GetTsumiki(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleTsumiki()

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		QueryRow(gomock.Any(), expected.ID).
		Return(setupTsumikiRow(ctrl, expected))

	tsumiki, err := repository.NewTsumikiRepository(db).GetTsumiki(nil, expected.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tsumiki == nil {
		t.Fatal("expected tsumiki, got nil")
	}
	if tsumiki.ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, tsumiki.ID)
	}
	if tsumiki.Title != expected.Title {
		t.Errorf("Title: want %s, got %s", expected.Title, tsumiki.Title)
	}
}

func TestTsumikiRepository_GetTsumiki_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		QueryRow(gomock.Any(), 999).
		Return(newNotFoundRowScanner(ctrl, 31))

	tsumiki, err := repository.NewTsumikiRepository(db).GetTsumiki(nil, 999)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tsumiki != nil {
		t.Errorf("want nil, got %+v", tsumiki)
	}
}

func TestTsumikiRepository_GetTsumikis(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleTsumiki()

	rows := newSingleRowsScanner(ctrl, 31, makeTsumikiScanFn(expected))
	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		Query(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(rows, nil)

	tsumikis, err := repository.NewTsumikiRepository(db).GetTsumikis(nil, 10, 1, nil, nil, "")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tsumikis) != 1 {
		t.Fatalf("want 1 tsumiki, got %d", len(tsumikis))
	}
	if tsumikis[0].ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, tsumikis[0].ID)
	}
}

func TestTsumikiRepository_GetTsumikiBlocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	block := sampleTsumikiBlock()

	// GetTsumikiBlocks の Scan: 13 フィールド (nullable 多数)
	// id, deletedAt, message, percentage, condition, createdAt, updatedAt,
	// mediaID, mediaType, mediaUrl, mediaOrder, mediaCreatedAt, mediaUpdatedAt
	blockScanFn := func(dest ...any) error {
		*dest[0].(*int) = block.ID
		*dest[1].(*sql.NullTime) = sql.NullTime{}
		*dest[2].(*sql.NullString) = sql.NullString{String: *block.Message, Valid: true}
		*dest[3].(*sql.NullInt64) = sql.NullInt64{Int64: int64(block.Percentage), Valid: true}
		*dest[4].(*sql.NullInt64) = sql.NullInt64{Int64: int64(block.Condition), Valid: true}
		*dest[5].(*sql.NullTime) = sql.NullTime{Time: block.CreatedAt, Valid: true}
		*dest[6].(*sql.NullTime) = sql.NullTime{Time: block.UpdatedAt, Valid: true}
		*dest[7].(*sql.NullInt64) = sql.NullInt64{} // media なし
		*dest[8].(*sql.NullString) = sql.NullString{}
		*dest[9].(*sql.NullString) = sql.NullString{}
		*dest[10].(*sql.NullInt64) = sql.NullInt64{}
		*dest[11].(*sql.NullTime) = sql.NullTime{}
		*dest[12].(*sql.NullTime) = sql.NullTime{}
		return nil
	}
	rows := newSingleRowsScanner(ctrl, 13, blockScanFn)

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		Query(gomock.Any(), block.TsumikiId, gomock.Any(), gomock.Any()).
		Return(rows, nil)

	blocks, err := repository.NewTsumikiRepository(db).GetTsumikiBlocks(block.TsumikiId, 10, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blocks) != 1 {
		t.Fatalf("want 1 block, got %d", len(blocks))
	}
	if blocks[0].ID != block.ID {
		t.Errorf("ID: want %d, got %d", block.ID, blocks[0].ID)
	}
	if blocks[0].IsDeleted {
		t.Error("want IsDeleted=false")
	}
}

func TestTsumikiRepository_CreateTsumiki(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleTsumiki()
	const insertedID = int64(3)

	db := mock.NewMockDBTX(ctrl)
	gomock.InOrder(
		db.EXPECT().
			Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&stubResult{lastInsertID: insertedID}, nil),
		db.EXPECT().
			QueryRow(gomock.Any(), int(insertedID)).
			Return(setupTsumikiRow(ctrl, expected)),
	)

	tsumiki, err := repository.NewTsumikiRepository(db).CreateTsumiki(
		expected.User.ID, expected.Title, expected.Visibility, nil, 0,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tsumiki == nil {
		t.Fatal("expected tsumiki, got nil")
	}
	if tsumiki.Title != expected.Title {
		t.Errorf("Title: want %s, got %s", expected.Title, tsumiki.Title)
	}
}

func TestTsumikiRepository_UpdateTsumikiThumbnail(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		Exec(gomock.Any(), 10, 3).
		Return(&stubResult{}, nil)

	err := repository.NewTsumikiRepository(db).UpdateTsumikiThumbnail(3, 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTsumikiRepository_UpdateTsumiki(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleTsumiki()

	db := mock.NewMockDBTX(ctrl)
	gomock.InOrder(
		db.EXPECT().
			Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&stubResult{}, nil),
		db.EXPECT().
			QueryRow(gomock.Any(), expected.ID).
			Return(setupTsumikiRow(ctrl, expected)),
	)

	tsumiki, err := repository.NewTsumikiRepository(db).UpdateTsumiki(
		expected.ID, expected.Title, expected.Visibility, nil,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tsumiki == nil {
		t.Fatal("expected tsumiki, got nil")
	}
	if tsumiki.ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, tsumiki.ID)
	}
}

func TestTsumikiRepository_DeleteTsumiki(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		Exec(gomock.Any(), 3).
		Return(&stubResult{}, nil)

	err := repository.NewTsumikiRepository(db).DeleteTsumiki(3)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTsumikiRepository_CreateMedia(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleTsumikiBlockMedia()
	const insertedID = int64(7)

	// CreateMedia in tsumiki.go: Scan 5 fields (id, type, url, created_at, updated_at)
	mediaRow := func() *mock.MockRowScanner {
		row := mock.NewMockRowScanner(ctrl)
		row.EXPECT().Scan(makeAnyArgs(5)...).DoAndReturn(func(dest ...any) error {
			*dest[0].(*int) = expected.ID
			*dest[1].(*string) = expected.Type
			*dest[2].(*string) = expected.Url
			*dest[3].(*time.Time) = expected.CreatedAt
			*dest[4].(*time.Time) = expected.UpdatedAt
			return nil
		})
		return row
	}()

	db := mock.NewMockDBTX(ctrl)
	gomock.InOrder(
		db.EXPECT().
			Exec(gomock.Any(), expected.Type, expected.Url).
			Return(&stubResult{lastInsertID: insertedID}, nil),
		db.EXPECT().
			QueryRow(gomock.Any(), insertedID).
			Return(mediaRow),
	)

	media, err := repository.NewTsumikiRepository(db).CreateMedia(3, expected.Type, expected.Url)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if media == nil {
		t.Fatal("expected media, got nil")
	}
	if media.ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, media.ID)
	}
}
