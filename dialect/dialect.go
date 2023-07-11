package dialect

import (
	"errors"
	"io"
	"time"
)

var (
	ErrEmptyIdent   = errors.New("empty ident")
	ErrNoIdentValue = errors.New("no ident type")
	ErrInvalidIdent = errors.New("invalid ident")
	ErrAsterisk     = errors.New("not allowed asterisk")
	ErrIdentQuote   = errors.New("not allowed ident quote or brackets")
)

type Writer interface {
	io.Writer
	io.ByteWriter
	io.StringWriter
}

type Dialect interface {
	WriteIdent(w Writer, v string) error
	WriteNull(w Writer) error
	WriteInt64(w Writer, v int64) error
	WriteFloat64(w Writer, v float64) error
	WriteBool(w Writer, v bool) error
	WriteString(w Writer, v string) error
	WriteBytes(w Writer, v []byte) error
	WriteTime(w Writer, v time.Time) error
}

var d Dialect

func SetDialect(di Dialect) {
	d = di
}

func GetDialect() Dialect {
	if d == nil {
		panic("dialect is not ready.")
	}
	return d
}
