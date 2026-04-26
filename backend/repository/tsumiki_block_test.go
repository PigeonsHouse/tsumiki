package repository_test

import (
	"testing"
	"tsumiki/repository"
	"tsumiki/repository/mock"

	"go.uber.org/mock/gomock"
)

func TestTsumikiBlockRepository_IsBelongToTsumiki(t *testing.T) {
	ctrl := gomock.NewController(t)
	block := sampleTsumikiBlock()

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		QueryRow(gomock.Any(), block.ID, block.TsumikiId).
		Return(newIntRowScanner(ctrl, 1))

	belongs, err := repository.NewTsumikiBlockRepository(db).IsBelongToTsumiki(block.TsumikiId, block.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !belongs {
		t.Error("want belongs=true, got false")
	}
}

func TestTsumikiBlockRepository_GetLatestBlockID(t *testing.T) {
	ctrl := gomock.NewController(t)
	block := sampleTsumikiBlock()

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		QueryRow(gomock.Any(), block.TsumikiId).
		Return(newIntRowScanner(ctrl, block.ID))

	id, err := repository.NewTsumikiBlockRepository(db).GetLatestBlockID(block.TsumikiId)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == nil {
		t.Fatal("expected id, got nil")
	}
	if *id != block.ID {
		t.Errorf("ID: want %d, got %d", block.ID, *id)
	}
}

func TestTsumikiBlockRepository_IsBelongToTsumiki_NotBelong(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		QueryRow(gomock.Any(), 999, 3).
		Return(newIntRowScanner(ctrl, 0))

	belongs, err := repository.NewTsumikiBlockRepository(db).IsBelongToTsumiki(3, 999)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if belongs {
		t.Error("want belongs=false, got true")
	}
}

func TestTsumikiBlockRepository_GetLatestBlockID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		QueryRow(gomock.Any(), 3).
		Return(newErrNoRowsScanner(ctrl))

	id, err := repository.NewTsumikiBlockRepository(db).GetLatestBlockID(3)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != nil {
		t.Errorf("want nil, got %d", *id)
	}
}

func TestTsumikiBlockRepository_CreateBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	block := sampleTsumikiBlock()
	const insertedID = int64(5)

	db := mock.NewMockDBTX(ctrl)
	gomock.InOrder(
		// fetchTailBlockID → 末尾ブロックなし
		db.EXPECT().
			QueryRow(gomock.Any(), block.TsumikiId).
			Return(newErrNoRowsScanner(ctrl)),
		// INSERT
		db.EXPECT().
			Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&stubResult{lastInsertID: insertedID}, nil),
		// fetchBlock
		db.EXPECT().
			QueryRow(gomock.Any(), int(insertedID)).
			Return(setupBlockRow(ctrl, block)),
	)

	result, err := repository.NewTsumikiBlockRepository(db).CreateBlock(
		block.TsumikiId, block.Message, block.Percentage, block.Condition,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected block, got nil")
	}
	if result.ID != block.ID {
		t.Errorf("ID: want %d, got %d", block.ID, result.ID)
	}
}

func TestTsumikiBlockRepository_UpdateBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	block := sampleTsumikiBlock()

	db := mock.NewMockDBTX(ctrl)
	gomock.InOrder(
		db.EXPECT().
			Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&stubResult{}, nil),
		db.EXPECT().
			QueryRow(gomock.Any(), block.ID).
			Return(setupBlockRow(ctrl, block)),
	)

	result, err := repository.NewTsumikiBlockRepository(db).UpdateBlock(
		block.ID, block.Message, block.Percentage, block.Condition,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected block, got nil")
	}
	if result.ID != block.ID {
		t.Errorf("ID: want %d, got %d", block.ID, result.ID)
	}
}

func TestTsumikiBlockRepository_SoftDeleteBlock(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		Exec(gomock.Any(), 5).
		Return(&stubResult{}, nil)

	err := repository.NewTsumikiBlockRepository(db).SoftDeleteBlock(5)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
