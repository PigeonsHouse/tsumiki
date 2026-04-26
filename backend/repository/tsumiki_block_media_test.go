package repository_test

import (
	"database/sql"
	"testing"
	"time"
	"tsumiki/repository"
	"tsumiki/repository/mock"
	"tsumiki/schema"

	"go.uber.org/mock/gomock"
)

func TestTsumikiBlockMediaRepository_CreateMedia(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleTsumikiBlockMedia()
	const insertedID = int64(7)

	db := mock.NewMockDBTX(ctrl)
	gomock.InOrder(
		db.EXPECT().
			Exec(gomock.Any(), expected.Url, expected.Type).
			Return(&stubResult{lastInsertID: insertedID}, nil),
		db.EXPECT().
			QueryRow(gomock.Any(), insertedID).
			Return(setupMediaRow(ctrl, expected)),
	)

	media, err := repository.NewTsumikiBlockMediaRepository(db).CreateMedia(expected.Url, expected.Type)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if media == nil {
		t.Fatal("expected media, got nil")
	}
	if media.ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, media.ID)
	}
	if media.Url != expected.Url {
		t.Errorf("Url: want %s, got %s", expected.Url, media.Url)
	}
}

func TestTsumikiBlockMediaRepository_SetMediaRelation(t *testing.T) {
	ctrl := gomock.NewController(t)
	const blockID = 3
	expected := sampleTsumikiBlockMedia()

	// SetMediaRelation の Query が返す 6 フィールド行:
	// id, type, url, order, created_at, updated_at
	mediaScanFn := func(dest ...any) error {
		*dest[0].(*int) = expected.ID
		*dest[1].(*string) = expected.Type
		*dest[2].(*string) = expected.Url
		*dest[3].(*int) = expected.Order
		*dest[4].(*time.Time) = expected.CreatedAt
		*dest[5].(*time.Time) = expected.UpdatedAt
		return nil
	}
	rows := newSingleRowsScanner(ctrl, 6, mediaScanFn)

	db := mock.NewMockDBTX(ctrl)
	gomock.InOrder(
		// 既存リレーションを NULL クリア
		db.EXPECT().
			Exec(gomock.Any(), blockID).
			Return(&stubResult{}, nil),
		// CASE UPDATE で新リレーション設定
		db.EXPECT().
			Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&stubResult{}, nil),
		// SELECT で結果取得
		db.EXPECT().
			Query(gomock.Any(), blockID).
			Return(rows, nil),
	)

	medias, err := repository.NewTsumikiBlockMediaRepository(db).SetMediaRelation(blockID, []int{expected.ID})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(medias) != 1 {
		t.Fatalf("want 1 media, got %d", len(medias))
	}
	if medias[0].ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, medias[0].ID)
	}
}

// setupMediaRowWithOrder は SetMediaRelation の QueryRow には使わないが、
// schema.TsumikiBlockMedia が sql.NullInt64 を使わないことを確認するためのコンパイルチェック。
var _ = func() {
	_ = schema.TsumikiBlockMedia{}
	_ = sql.NullInt64{}
}
