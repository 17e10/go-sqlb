package sqlt

import (
	"context"
	"database/sql"
)

type testResult struct{}

func (testResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (testResult) RowsAffected() (int64, error) {
	return 0, nil
}

type TestExeced struct {
	Query string
	Args  []any
}

type TestExecer struct {
	Execed []TestExeced
}

func (e *TestExecer) ExecContext(_ context.Context, query string, args ...any) (sql.Result, error) {
	e.Execed = append(e.Execed, TestExeced{query, args})
	return testResult{}, nil
}

type NullExecer struct{}

func (e NullExecer) ExecContext(_ context.Context, query string, args ...any) (sql.Result, error) {
	return testResult{}, nil
}
