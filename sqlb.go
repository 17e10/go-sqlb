package sqlb

import (
	"errors"
	"strings"

	"github.com/17e10/go-nameb"
	d "github.com/17e10/go-sqlb/dialect"
)

var (
	errNoIdentType = errors.New("no ident type")
	errNoValueType = errors.New("no value type")
	errNoSuchKey   = errors.New("no such key")
	errOutRange    = errors.New("out of range")
	errEmptySlice  = errors.New("empty array or slice")
	errNoStruct    = errors.New("no struct")
)

var dialect = d.GetDialect

// Writer は SQL を書き込むインターフェイスです.
//
// 標準的な Writer は bytes.Buffer や strings.Builder です.
type Writer = d.Writer

// Sqler は SQL を構築するインターフェイスを表します.
type Sqler interface {
	Sql(w Writer) error
}

// SqlerFunc 型は関数を Sqler として使用できるようにするアダプタです.
type SqlerFunc func(w Writer) error

// Sql は fn(w) を呼び出します.
func (fn SqlerFunc) Sql(w Writer) error {
	return fn(w)
}

// Sqler から文字列を得ます.
func Stringify(sqler Sqler) (string, error) {
	const cap = 256

	w := &strings.Builder{}
	w.Grow(cap)
	if err := sqler.Sql(w); err != nil {
		return "", err
	}
	return w.String(), nil
}

// StringSqler は Sqler インターフェイスを持つ文字列型です.
type StringSqler string

func (s StringSqler) Sql(w Writer) error {
	w.WriteString(string(s))
	return nil
}

// ColumnCase は sqlb が生成するカラム名の形式を保持します.
var ColumnCase = nameb.Snake

// 構造体のフィールド名をデータベースのカラム名に変換します.
func columnName(s string) string {
	return ColumnCase(s)
}
