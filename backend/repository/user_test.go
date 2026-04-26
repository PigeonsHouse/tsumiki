package repository_test

import (
	"testing"
	"tsumiki/repository"
	"tsumiki/repository/mock"

	"go.uber.org/mock/gomock"
)

func TestUserRepository_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleUser()

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		QueryRow(gomock.Any(), expected.ID).
		Return(setupUserRow(ctrl, expected))

	user, err := repository.NewUserRepository(db).FindByID(expected.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.ID != expected.ID {
		t.Errorf("ID: want %d, got %d", expected.ID, user.ID)
	}
	if user.Name != expected.Name {
		t.Errorf("Name: want %s, got %s", expected.Name, user.Name)
	}
	if user.AvatarUrl != expected.AvatarUrl {
		t.Errorf("AvatarUrl: want %s, got %s", expected.AvatarUrl, user.AvatarUrl)
	}
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		QueryRow(gomock.Any(), 999).
		Return(newNotFoundRowScanner(ctrl, 7))

	user, err := repository.NewUserRepository(db).FindByID(999)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user != nil {
		t.Errorf("want nil, got %+v", user)
	}
}

func TestUserRepository_FindByDiscordUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := sampleUser()

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		QueryRow(gomock.Any(), expected.DiscordUserID).
		Return(setupUserRow(ctrl, expected))

	user, err := repository.NewUserRepository(db).FindByDiscordUserId(expected.DiscordUserID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.DiscordUserID != expected.DiscordUserID {
		t.Errorf("DiscordUserID: want %s, got %s", expected.DiscordUserID, user.DiscordUserID)
	}
}

func TestUserRepository_FindByDiscordUserId_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		QueryRow(gomock.Any(), "unknown").
		Return(newNotFoundRowScanner(ctrl, 7))

	user, err := repository.NewUserRepository(db).FindByDiscordUserId("unknown")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user != nil {
		t.Errorf("want nil, got %+v", user)
	}
}

func TestUserRepository_CreateUserByDiscord(t *testing.T) {
	ctrl := gomock.NewController(t)
	const insertedID = int64(42)

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&stubResult{lastInsertID: insertedID}, nil)

	user, err := repository.NewUserRepository(db).CreateUserByDiscord(
		"New User", "avatars/42/img.png", "discord456", "guild456",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.ID != int(insertedID) {
		t.Errorf("ID: want %d, got %d", insertedID, user.ID)
	}
	if user.Name != "New User" {
		t.Errorf("Name: want %s, got %s", "New User", user.Name)
	}
}

func TestUserRepository_UpdateAvatarUrl(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := mock.NewMockDBTX(ctrl)
	db.EXPECT().
		Exec(gomock.Any(), "avatars/1/new.png", 1).
		Return(&stubResult{}, nil)

	err := repository.NewUserRepository(db).UpdateAvatarUrl(1, "avatars/1/new.png")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
