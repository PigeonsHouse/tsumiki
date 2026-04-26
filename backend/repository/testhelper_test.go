package repository_test

import "database/sql"

type stubResult struct{ lastInsertID int64 }

func (s *stubResult) LastInsertId() (int64, error) { return s.lastInsertID, nil }
func (s *stubResult) RowsAffected() (int64, error) { return 1, nil }

var _ sql.Result = (*stubResult)(nil)
