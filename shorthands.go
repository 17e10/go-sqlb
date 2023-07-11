package sqlb

import (
	"context"
	"database/sql"

	"github.com/17e10/go-sqlb/sqlt"
)

// Query は database/sql の QueryContext メソッドを Sqler で呼び出すショートハンドです.
func Query(conn sqlt.Queryer, ctx context.Context, sqler Sqler) (*sql.Rows, error) {
	query, err := Stringify(sqler)
	if err != nil {
		return nil, err
	}
	return conn.QueryContext(ctx, query)
}

// QueryRow は database/sql の QueryRowContext メソッドを Sqler で呼び出すショートハンドです.
func QueryRow(conn sqlt.QueryRower, ctx context.Context, sqler Sqler) *sql.Row {
	query, _ := Stringify(sqler)
	return conn.QueryRowContext(ctx, query)
}

// Exec は database/sql の ExecContext メソッドを Sqler で呼び出すショートハンドです.
func Exec(conn sqlt.Execer, ctx context.Context, sqler Sqler) (sql.Result, error) {
	query, err := Stringify(sqler)
	if err != nil {
		return nil, err
	}
	return conn.ExecContext(ctx, query)
}
