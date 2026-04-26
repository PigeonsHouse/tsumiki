package repository_test

import (
	"testing"
	"time"
	"tsumiki/repository"
	"tsumiki/repository/mock"
	"tsumiki/schema"

	"go.uber.org/mock/gomock"
)

func sampleUser() *schema.User {
	guildID := "guild123"
	return &schema.User{
		ID:            1,
		DiscordUserID: "discord123",
		Name:          "Test User",
		GuildID:       &guildID,
		AvatarUrl:     "avatars/1/abc.png",
		CreatedAt:     time.Now().Truncate(time.Second),
		UpdatedAt:     time.Now().Truncate(time.Second),
	}
}

// setupUserRow は MockRowScanner に schema.User の値を Scan させるセットアップをまとめたヘルパー。
func setupUserRow(ctrl *gomock.Controller, u *schema.User) *mock.MockRowScanner {
	row := mock.NewMockRowScanner(ctrl)
	row.EXPECT().Scan(
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(),
	).DoAndReturn(func(dest ...any) error {
		*dest[0].(*int) = u.ID
		*dest[1].(*string) = u.DiscordUserID
		*dest[2].(*string) = u.Name
		*dest[3].(**string) = u.GuildID
		*dest[4].(*string) = u.AvatarUrl
		*dest[5].(*time.Time) = u.CreatedAt
		*dest[6].(*time.Time) = u.UpdatedAt
		return nil
	})
	return row
}

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
