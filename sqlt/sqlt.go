package sqlt

import (
	"context"
	"database/sql"
)

// Preparer は database/sql の PrepareContext メソッドをラップするインターフェイスです.
type Preparer interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

// Queryer は database/sql の QueryContext メソッドをラップするインターフェイスです.
type Queryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// QueryRower は database/sql の QueryRowContext メソッドをラップするインターフェイスです.
type QueryRower interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// RowsScanner は database/sql の Scan メソッドをラップするインターフェイスです.
type RowsScanner interface {
	Scan(dest ...any) error
}

// Execer は database/sql の ExecContext メソッドをラップするインターフェイスです.
type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
