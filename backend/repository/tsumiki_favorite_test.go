package repository_test

import (
	"errors"
	"testing"
	"tsumiki/repository"
	"tsumiki/repository/mock"
	"tsumiki/schema"

	"go.uber.org/mock/gomock"
)

func TestTsumikiFavoriteRepository_SetFavoriteCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	tsumiki := sampleTsumiki()
	userID := 12
	expected := &schema.Favorite{
		TotalFavoriteCount: 5,
		MyFavoriteCount:    5,
	}

	db := mock.NewMockDBTX(ctrl)
	gomock.InOrder(
		db.EXPECT().
			Exec(gomock.Any(), tsumiki.ID, userID, expected.MyFavoriteCount).
			Return(&stubResult{}, nil),
		db.EXPECT().
			QueryRow(gomock.Any(), tsumiki.ID, userID).
			Return(setupFavoriteRow(ctrl, expected)),
	)

	countData, err := repository.NewTsumikiFavoriteRepository(db).SetFavoriteCount(tsumiki.ID, userID, expected.MyFavoriteCount)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if countData == nil {
		t.Fatal("expected countData, got nil")
	}
	if countData.MyFavoriteCount != expected.MyFavoriteCount {
		t.Errorf("MyFavoriteCount: want %d, got %d", expected.MyFavoriteCount, countData.MyFavoriteCount)
	}
	if countData.TotalFavoriteCount != expected.TotalFavoriteCount {
		t.Errorf("TotalFavoriteCount: want %d, got %d", expected.TotalFavoriteCount, countData.TotalFavoriteCount)
	}
}

func TestTsumikiFavoriteRepository_SetFavoriteCount_ExecError(t *testing.T) {
	ctrl := gomock.NewController(t)
	tsumiki := sampleTsumiki()
	userID := 12
	count := 5

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		Exec(gomock.Any(), tsumiki.ID, userID, count).
		Return(nil, errors.New("exec error"))

	countData, err := repository.NewTsumikiFavoriteRepository(db).SetFavoriteCount(tsumiki.ID, userID, count)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if countData != nil {
		t.Errorf("expected nil countData, got %v", countData)
	}
}

func TestTsumikiFavoriteRepository_SetFavoriteCount_ScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	tsumiki := sampleTsumiki()
	userID := 12
	count := 5

	db := mock.NewMockDBTX(ctrl)
	gomock.InOrder(
		db.EXPECT().
			Exec(gomock.Any(), tsumiki.ID, userID, count).
			Return(&stubResult{}, nil),
		db.EXPECT().
			QueryRow(gomock.Any(), tsumiki.ID, userID).
			Return(newErrNoRowsScanner(ctrl)),
	)

	countData, err := repository.NewTsumikiFavoriteRepository(db).SetFavoriteCount(tsumiki.ID, userID, count)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if countData != nil {
		t.Errorf("expected nil countData, got %v", countData)
	}
}
