package repository_test

import (
	"testing"
	"tsumiki/repository"
	"tsumiki/repository/mock"

	"go.uber.org/mock/gomock"
)

func TestWorkRepository_GetWork(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleWork()

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		QueryRow(gomock.Any(), expected.ID).
		Return(setupWorkRow(ctrl, expected))

	work, err := repository.NewWorkRepository(db).GetWork(expected.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if work == nil {
		t.Fatal("expected work, got nil")
	}
	if work.ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, work.ID)
	}
	if work.Title != expected.Title {
		t.Errorf("Title: want %s, got %s", expected.Title, work.Title)
	}
	if work.Owner.ID != expected.Owner.ID {
		t.Errorf("Owner.ID: want %d, got %d", expected.Owner.ID, work.Owner.ID)
	}
}

func TestWorkRepository_GetWork_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		QueryRow(gomock.Any(), 999).
		Return(newNotFoundRowScanner(ctrl, 17))

	work, err := repository.NewWorkRepository(db).GetWork(999)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if work != nil {
		t.Errorf("want nil, got %+v", work)
	}
}

func TestWorkRepository_GetWorks(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleWork()

	rows := newSingleRowsScanner(ctrl, 17, makeWorkScanFn(expected))
	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		Query(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(rows, nil)

	works, err := repository.NewWorkRepository(db).GetWorks(nil, 10, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(works) != 1 {
		t.Fatalf("want 1 work, got %d", len(works))
	}
	if works[0].ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, works[0].ID)
	}
}

func TestWorkRepository_CreateWork(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleWork()
	const insertedID = int64(2)

	db := mock.NewMockDBTX(ctrl)
	gomock.InOrder(
		db.EXPECT().
			Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&stubResult{lastInsertID: insertedID}, nil),
		db.EXPECT().
			QueryRow(gomock.Any(), int(insertedID)).
			Return(setupWorkRow(ctrl, expected)),
	)

	work, err := repository.NewWorkRepository(db).CreateWork(
		expected.Owner.ID, expected.Title, expected.Visibility, expected.Description, nil,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if work == nil {
		t.Fatal("expected work, got nil")
	}
	if work.Title != expected.Title {
		t.Errorf("Title: want %s, got %s", expected.Title, work.Title)
	}
}

func TestWorkRepository_UpdateWorkThumbnail(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		Exec(gomock.Any(), 10, 2).
		Return(&stubResult{}, nil)

	err := repository.NewWorkRepository(db).UpdateWorkThumbnail(2, 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkRepository_UpdateWork(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleWork()

	db := mock.NewMockDBTX(ctrl)
	gomock.InOrder(
		db.EXPECT().
			Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&stubResult{}, nil),
		db.EXPECT().
			QueryRow(gomock.Any(), expected.ID).
			Return(setupWorkRow(ctrl, expected)),
	)

	work, err := repository.NewWorkRepository(db).UpdateWork(
		expected.ID, expected.Title, expected.Visibility, expected.Description,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if work == nil {
		t.Fatal("expected work, got nil")
	}
	if work.ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, work.ID)
	}
}

func TestWorkRepository_DeleteWork(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		Exec(gomock.Any(), 2).
		Return(&stubResult{}, nil)

	err := repository.NewWorkRepository(db).DeleteWork(2)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
