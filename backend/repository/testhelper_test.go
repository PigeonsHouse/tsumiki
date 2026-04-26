package repository_test

import (
	"database/sql"
	"time"
	"tsumiki/repository/mock"
	"tsumiki/schema"

	"go.uber.org/mock/gomock"
)

// --- sql.Result stub ---

type stubResult struct{ lastInsertID int64 }

func (s *stubResult) LastInsertId() (int64, error) { return s.lastInsertID, nil }
func (s *stubResult) RowsAffected() (int64, error) { return 1, nil }

var _ sql.Result = (*stubResult)(nil)

// --- generic helpers ---

// makeAnyArgs は variadic な Scan の EXPECT に渡す n 個の gomock.Any() スライスを返す。
func makeAnyArgs(n int) []any {
	args := make([]any, n)
	for i := range args {
		args[i] = gomock.Any()
	}
	return args
}

// newErrNoRowsScanner は Scan(単一引数) が sql.ErrNoRows を返す RowScanner を返す。
// COUNT や EXISTS など 1 フィールドを返す QueryRow の「行なし」ケースに使う。
func newErrNoRowsScanner(ctrl *gomock.Controller) *mock.MockRowScanner {
	return newNotFoundRowScanner(ctrl, 1)
}

// newNotFoundRowScanner は numFields 個の引数を受け取る Scan が sql.ErrNoRows を返す RowScanner を返す。
// 複数フィールドを持つ QueryRow の「レコードなし」ケースに使う。
func newNotFoundRowScanner(ctrl *gomock.Controller, numFields int) *mock.MockRowScanner {
	row := mock.NewMockRowScanner(ctrl)
	row.EXPECT().Scan(makeAnyArgs(numFields)...).Return(sql.ErrNoRows)
	return row
}

// newIntRowScanner は Scan(単一 *int 引数) で value をセットする RowScanner を返す。
// COUNT や単一 ID を返す QueryRow に使う。
func newIntRowScanner(ctrl *gomock.Controller, value int) *mock.MockRowScanner {
	row := mock.NewMockRowScanner(ctrl)
	row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
		*dest[0].(*int) = value
		return nil
	})
	return row
}

// newSingleRowsScanner は 1 行だけを返す RowsScanner を生成する。
// Next→Scan→Next→Err→Close の順序を InOrder で保証する。
func newSingleRowsScanner(ctrl *gomock.Controller, numFields int, scanFn func(dest ...any) error) *mock.MockRowsScanner {
	rows := mock.NewMockRowsScanner(ctrl)
	gomock.InOrder(
		rows.EXPECT().Next().Return(true),
		rows.EXPECT().Scan(makeAnyArgs(numFields)...).DoAndReturn(scanFn),
		rows.EXPECT().Next().Return(false),
		rows.EXPECT().Err().Return(nil),
		rows.EXPECT().Close().Return(nil),
	)
	return rows
}

// --- sample data factories ---

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

func sampleThumbnail() *schema.ThumbnailUpload {
	return &schema.ThumbnailUpload{
		ID:        10,
		Url:       "thumbnails/1/abc.png",
		CreatedAt: time.Now().Truncate(time.Second),
		UpdatedAt: time.Now().Truncate(time.Second),
	}
}

func sampleWork() *schema.Work {
	owner := sampleUser()
	return &schema.Work{
		ID:          2,
		Title:       "Test Work",
		Description: "desc",
		Visibility:  "public",
		Owner:       *owner,
		CreatedAt:   time.Now().Truncate(time.Second),
		UpdatedAt:   time.Now().Truncate(time.Second),
	}
}

func sampleTsumikiBlock() *schema.TsumikiBlock {
	msg := "test message"
	return &schema.TsumikiBlock{
		ID:         5,
		Message:    &msg,
		Medias:     []schema.TsumikiBlockMedia{},
		Percentage: 50,
		Condition:  1,
		TsumikiId:  3,
		CreatedAt:  time.Now().Truncate(time.Second),
		UpdatedAt:  time.Now().Truncate(time.Second),
	}
}

func sampleTsumikiBlockMedia() *schema.TsumikiBlockMedia {
	return &schema.TsumikiBlockMedia{
		ID:        7,
		Type:      "image",
		Url:       "medias/abc.png",
		Order:     0,
		CreatedAt: time.Now().Truncate(time.Second),
		UpdatedAt: time.Now().Truncate(time.Second),
	}
}

func sampleTsumiki() *schema.Tsumiki {
	return &schema.Tsumiki{
		ID:         3,
		Title:      "Test Tsumiki",
		Visibility: "public",
		User:       *sampleUser(),
		CreatedAt:  time.Now().Truncate(time.Second),
		UpdatedAt:  time.Now().Truncate(time.Second),
	}
}

// --- scan functions (QueryRow / RowsScanner 両方から再利用) ---

// scanUser: id, discord_user_id, name, guild_id, avatar_url, created_at, updated_at (7 fields)
func makeUserScanFn(u *schema.User) func(dest ...any) error {
	return func(dest ...any) error {
		*dest[0].(*int) = u.ID
		*dest[1].(*string) = u.DiscordUserID
		*dest[2].(*string) = u.Name
		*dest[3].(**string) = u.GuildID
		*dest[4].(*string) = u.AvatarUrl
		*dest[5].(*time.Time) = u.CreatedAt
		*dest[6].(*time.Time) = u.UpdatedAt
		return nil
	}
}

// scanThumbnail: id, path, created_at, updated_at (4 fields)
func makeThumbnailScanFn(th *schema.ThumbnailUpload) func(dest ...any) error {
	return func(dest ...any) error {
		*dest[0].(*int) = th.ID
		*dest[1].(*string) = th.Url
		*dest[2].(*time.Time) = th.CreatedAt
		*dest[3].(*time.Time) = th.UpdatedAt
		return nil
	}
}

// scanWork: 17 fields
// w.ID, Title, Description, Visibility, CreatedAt, UpdatedAt,
// Owner.ID, DiscordUserID, Name, GuildID, AvatarUrl, CreatedAt, UpdatedAt,
// thID(NullInt64), thPath(NullString), thCreatedAt(NullTime), thUpdatedAt(NullTime)
func makeWorkScanFn(w *schema.Work) func(dest ...any) error {
	return func(dest ...any) error {
		*dest[0].(*int) = w.ID
		*dest[1].(*string) = w.Title
		*dest[2].(*string) = w.Description
		*dest[3].(*string) = w.Visibility
		*dest[4].(*time.Time) = w.CreatedAt
		*dest[5].(*time.Time) = w.UpdatedAt
		*dest[6].(*int) = w.Owner.ID
		*dest[7].(*string) = w.Owner.DiscordUserID
		*dest[8].(*string) = w.Owner.Name
		*dest[9].(**string) = w.Owner.GuildID
		*dest[10].(*string) = w.Owner.AvatarUrl
		*dest[11].(*time.Time) = w.Owner.CreatedAt
		*dest[12].(*time.Time) = w.Owner.UpdatedAt
		*dest[13].(*sql.NullInt64) = sql.NullInt64{}
		*dest[14].(*sql.NullString) = sql.NullString{}
		*dest[15].(*sql.NullTime) = sql.NullTime{}
		*dest[16].(*sql.NullTime) = sql.NullTime{}
		return nil
	}
}

// scanBlock: id, message(*string), percentage, condition, next_block_id(*int), tsumiki_id, created_at, updated_at (8 fields)
func makeBlockScanFn(b *schema.TsumikiBlock) func(dest ...any) error {
	return func(dest ...any) error {
		*dest[0].(*int) = b.ID
		*dest[1].(**string) = b.Message
		*dest[2].(*int) = b.Percentage
		*dest[3].(*int) = b.Condition
		*dest[4].(**int) = b.NextBlockId
		*dest[5].(*int) = b.TsumikiId
		*dest[6].(*time.Time) = b.CreatedAt
		*dest[7].(*time.Time) = b.UpdatedAt
		return nil
	}
}

// scanMedia: id, type, url, created_at, updated_at (5 fields)
func makeMediaScanFn(m *schema.TsumikiBlockMedia) func(dest ...any) error {
	return func(dest ...any) error {
		*dest[0].(*int) = m.ID
		*dest[1].(*string) = m.Type
		*dest[2].(*string) = m.Url
		*dest[3].(*time.Time) = m.CreatedAt
		*dest[4].(*time.Time) = m.UpdatedAt
		return nil
	}
}

// scanTsumiki: 31 fields (work=nil のケース)
// t.ID, Title, Visibility, CreatedAt, UpdatedAt,
// User.ID, DiscordUserID, Name, AvatarUrl, CreatedAt, UpdatedAt,
// workID...(NullInt64/NullString/NullTime ×6),
// ownerID...(NullInt64/NullString/NullTime ×6),
// wthID...(NullInt64/NullString/NullTime ×4),
// tthID...(NullInt64/NullString/NullTime ×4)
func makeTsumikiScanFn(t *schema.Tsumiki) func(dest ...any) error {
	return func(dest ...any) error {
		*dest[0].(*int) = t.ID
		*dest[1].(*string) = t.Title
		*dest[2].(*string) = t.Visibility
		*dest[3].(*time.Time) = t.CreatedAt
		*dest[4].(*time.Time) = t.UpdatedAt
		*dest[5].(*int) = t.User.ID
		*dest[6].(*string) = t.User.DiscordUserID
		*dest[7].(*string) = t.User.Name
		*dest[8].(*string) = t.User.AvatarUrl
		*dest[9].(*time.Time) = t.User.CreatedAt
		*dest[10].(*time.Time) = t.User.UpdatedAt
		// work = nil
		*dest[11].(*sql.NullInt64) = sql.NullInt64{}
		*dest[12].(*sql.NullString) = sql.NullString{}
		*dest[13].(*sql.NullString) = sql.NullString{}
		*dest[14].(*sql.NullString) = sql.NullString{}
		*dest[15].(*sql.NullTime) = sql.NullTime{}
		*dest[16].(*sql.NullTime) = sql.NullTime{}
		// work owner = nil
		*dest[17].(*sql.NullInt64) = sql.NullInt64{}
		*dest[18].(*sql.NullString) = sql.NullString{}
		*dest[19].(*sql.NullString) = sql.NullString{}
		*dest[20].(*sql.NullString) = sql.NullString{}
		*dest[21].(*sql.NullTime) = sql.NullTime{}
		*dest[22].(*sql.NullTime) = sql.NullTime{}
		// work thumbnail = nil
		*dest[23].(*sql.NullInt64) = sql.NullInt64{}
		*dest[24].(*sql.NullString) = sql.NullString{}
		*dest[25].(*sql.NullTime) = sql.NullTime{}
		*dest[26].(*sql.NullTime) = sql.NullTime{}
		// tsumiki thumbnail = nil
		*dest[27].(*sql.NullInt64) = sql.NullInt64{}
		*dest[28].(*sql.NullString) = sql.NullString{}
		*dest[29].(*sql.NullTime) = sql.NullTime{}
		*dest[30].(*sql.NullTime) = sql.NullTime{}
		return nil
	}
}

// --- RowScanner setup helpers (QueryRow 用) ---

func setupUserRow(ctrl *gomock.Controller, u *schema.User) *mock.MockRowScanner {
	row := mock.NewMockRowScanner(ctrl)
	row.EXPECT().Scan(makeAnyArgs(7)...).DoAndReturn(makeUserScanFn(u))
	return row
}

func setupThumbnailRow(ctrl *gomock.Controller, th *schema.ThumbnailUpload) *mock.MockRowScanner {
	row := mock.NewMockRowScanner(ctrl)
	row.EXPECT().Scan(makeAnyArgs(4)...).DoAndReturn(makeThumbnailScanFn(th))
	return row
}

func setupWorkRow(ctrl *gomock.Controller, w *schema.Work) *mock.MockRowScanner {
	row := mock.NewMockRowScanner(ctrl)
	row.EXPECT().Scan(makeAnyArgs(17)...).DoAndReturn(makeWorkScanFn(w))
	return row
}

func setupBlockRow(ctrl *gomock.Controller, b *schema.TsumikiBlock) *mock.MockRowScanner {
	row := mock.NewMockRowScanner(ctrl)
	row.EXPECT().Scan(makeAnyArgs(8)...).DoAndReturn(makeBlockScanFn(b))
	return row
}

func setupMediaRow(ctrl *gomock.Controller, m *schema.TsumikiBlockMedia) *mock.MockRowScanner {
	row := mock.NewMockRowScanner(ctrl)
	row.EXPECT().Scan(makeAnyArgs(5)...).DoAndReturn(makeMediaScanFn(m))
	return row
}

func setupTsumikiRow(ctrl *gomock.Controller, t *schema.Tsumiki) *mock.MockRowScanner {
	row := mock.NewMockRowScanner(ctrl)
	row.EXPECT().Scan(makeAnyArgs(31)...).DoAndReturn(makeTsumikiScanFn(t))
	return row
}
