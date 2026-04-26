package repository_test

import (
	"testing"
	"tsumiki/repository"
	"tsumiki/repository/mock"

	"go.uber.org/mock/gomock"
)

func TestThumbnailRepository_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleThumbnail()

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		QueryRow(gomock.Any(), expected.ID).
		Return(setupThumbnailRow(ctrl, expected))

	th, err := repository.NewThumbnailRepository(db).Get(expected.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected thumbnail, got nil")
	}
	if th.ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, th.ID)
	}
	if th.Url != expected.Url {
		t.Errorf("Url: want %s, got %s", expected.Url, th.Url)
	}
}

func TestThumbnailRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleThumbnail()
	const insertedID = int64(10)

	db := mock.NewMockDBTX(ctrl)
	gomock.InOrder(
		db.EXPECT().
			Exec(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&stubResult{lastInsertID: insertedID}, nil),
		db.EXPECT().
			QueryRow(gomock.Any(), int(insertedID)).
			Return(setupThumbnailRow(ctrl, expected)),
	)

	th, err := repository.NewThumbnailRepository(db).Create(1, expected.Url)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected thumbnail, got nil")
	}
	if th.ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, th.ID)
	}
}

func TestThumbnailRepository_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		QueryRow(gomock.Any(), 999).
		Return(newNotFoundRowScanner(ctrl, 4))

	th, err := repository.NewThumbnailRepository(db).Get(999)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th != nil {
		t.Errorf("want nil, got %+v", th)
	}
}

func TestThumbnailRepository_IsInUse(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleThumbnail()

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		QueryRow(gomock.Any(), expected.ID, expected.ID).
		Return(newIntRowScanner(ctrl, 1))

	inUse, err := repository.NewThumbnailRepository(db).IsInUse(expected.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !inUse {
		t.Error("want inUse=true, got false")
	}
}

func TestThumbnailRepository_IsInUse_NotInUse(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		QueryRow(gomock.Any(), 999, 999).
		Return(newErrNoRowsScanner(ctrl))

	inUse, err := repository.NewThumbnailRepository(db).IsInUse(999)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inUse {
		t.Error("want inUse=false, got true")
	}
}
