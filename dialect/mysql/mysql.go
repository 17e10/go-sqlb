package mysql

import (
	"encoding/hex"
	"fmt"
	"time"

	d "github.com/17e10/go-sqlb/dialect"
)

type Writer = d.Writer

func init() {
	d.SetDialect(mysql{})
}

type mysql struct{}

// WriteIdent は識別子の SQL 文字列を w に書き込みます.
func (mysql) WriteIdent(w Writer, s string) error {
	const (
		begin = iota
		ident
		end
	)

	if s == "" {
		return d.ErrEmptyIdent
	}

	state := begin
	w.WriteByte('`')
	for i, l := 0, len(s); i < l; i++ {
		c := s[i]
		switch c {
		case '`', '[', ']':
			return d.ErrIdentQuote
		case '*':
			return d.ErrAsterisk
		case ' ', '\t', '\n', '\r':
			if state == ident {
				state = end
			}
		case '.':
			if state == begin {
				return d.ErrInvalidIdent
			}
			w.WriteString("`.`")
			state = begin
		default:
			if c < ' ' || state == end {
				return d.ErrInvalidIdent
			}
			w.WriteByte(c)
			state = ident
		}
	}
	if state == begin {
		return d.ErrInvalidIdent
	}
	w.WriteByte('`')
	return nil
}

// WriteNull は NULL の SQL 文字列を w に書き込みます.
func (mysql) WriteNull(w Writer) error {
	w.WriteString("NULL")
	return nil
}

// WriteInt64 は整数の SQL 文字列を w に書き込みます.
func (mysql) WriteInt64(w Writer, v int64) error {
	fmt.Fprintf(w, "%d", v)
	return nil
}

// WriteFloat64 は浮動小数点の SQL 文字列を w に書き込みます.
func (mysql) WriteFloat64(w Writer, v float64) error {
	fmt.Fprintf(w, "%g", v)
	return nil
}

// WriteBool は bool の SQL 文字列を w に書き込みます.
func (mysql) WriteBool(w Writer, v bool) error {
	var s string
	if v {
		s = "TRUE"
	} else {
		s = "FALSE"
	}
	w.WriteString(s)
	return nil
}

// WriteBool は文字列の SQL 文字列を w に書き込みます.
// 書き込むとき適切にエスケープ処理をします.
func (mysql) WriteString(w Writer, s string) error {
	w.WriteByte('\'')
	for i, l := 0, len(s); i < l; i++ {
		c := s[i]
		switch c {
		case '\000':
			// nothing to do
		case '\'':
			w.WriteString(`\'`)
		case '\\':
			w.WriteString(`\\`)
		default:
			w.WriteByte(c)
		}
	}
	w.WriteByte('\'')
	return nil
}

func (d mysql) WriteBytes(w Writer, v []byte) error {
	if v == nil {
		return d.WriteNull(w)
	}

	w.WriteString("X'")
	hex.NewEncoder(w).Write(v)
	w.WriteByte('\'')
	return nil
}

func (mysql) WriteTime(w Writer, tm time.Time) error {
	w.WriteString(tm.Format("'2006-01-02 15:04:05.999999'"))
	return nil
}
