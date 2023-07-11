package sqlb

import (
	"reflect"

	"github.com/17e10/go-sqlb/sqlt"
)

// Scan は database/sql の Row(s).Scan メソッドの結果を構造体に入れます.
//
// Scan で受け取るフィールド順序は Columns で生成されるカラム列の順序と一致します.
func Scan[V any](row sqlt.RowsScanner, dest *V) error {
	d := makeFieldsAddr(dest)
	return row.Scan(d...)
}

// makeFieldsAddr は構造体から各要素のポインタ配列を作成します.
func makeFieldsAddr[V any](v *V) []any {
	cols := exportedColumns(v, nil)
	rv := reflect.ValueOf(v).Elem()
	r := make([]any, len(cols))
	for i, c := range cols {
		r[i] = rv.FieldByIndex(c.index).Addr().Interface()
	}
	return r
}
