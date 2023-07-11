package mysql

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"
)

func TestIdent(t *testing.T) {
	tests := []struct {
		src  string
		want string
		err  string
	}{
		// 正常ケース
		{"field", "`field`", ""},
		{"T.field", "`T`.`field`", ""},
		{" T . field ", "`T`.`field`", ""},

		// 構文エラー
		{"", "", "empty ident"},
		{".field", "", "invalid ident"},
		{"field.", "", "invalid ident"},
		{"T field", "", "invalid ident"},

		// 文字制約エラー
		{"*", "", "not allowed asterisk"},
		{"T.*", "", "not allowed asterisk"},
		{"`field`", "", "not allowed ident quote or brackets"},
		{"[field]", "", "not allowed ident quote or brackets"},
		{"foo\001bar", "", "invalid ident"},
	}

	for _, te := range tests {
		var (
			d      mysql
			w      strings.Builder
			goterr string
			got    string
		)

		name := fmt.Sprintf("WriteIdent(%q)", te.src)
		err := d.WriteIdent(&w, te.src)
		if err != nil {
			goterr = err.Error()
		} else {
			got = w.String()
		}
		if goterr != te.err {
			t.Errorf("%s errored %q, want %q", name, goterr, te.err)
		}
		if got != te.want {
			t.Errorf("%s returned %q, want %q", name, got, te.want)
		}
	}
}

func TestNull(t *testing.T) {
	var (
		d      mysql
		w      strings.Builder
		goterr string
		got    string
	)

	name, want, terr := "WriteNull()", "NULL", ""
	err := d.WriteNull(&w)
	if err != nil {
		goterr = err.Error()
	} else {
		got = w.String()
	}
	if goterr != terr {
		t.Errorf("%s errored %q, want %q", name, goterr, terr)
	}
	if got != want {
		t.Errorf("%s returned %q, want %q", name, got, want)
	}
}

func TestInt64(t *testing.T) {
	tests := []struct {
		src  int64
		want string
		err  string
	}{
		{0, "0", ""},
		{123, "123", ""},
		{math.MinInt64, "-9223372036854775808", ""},
		{math.MaxInt64, "9223372036854775807", ""},
	}

	for _, te := range tests {
		var (
			d      mysql
			w      strings.Builder
			goterr string
			got    string
		)

		name := fmt.Sprintf("WriteInt64(%d)", te.src)
		err := d.WriteInt64(&w, te.src)
		if err != nil {
			goterr = err.Error()
		} else {
			got = w.String()
		}
		if goterr != te.err {
			t.Errorf("%s errored %q, want %q", name, goterr, te.err)
		}
		if got != te.want {
			t.Errorf("%s returned %q, want %q", name, got, te.want)
		}
	}
}

func TestFloat64(t *testing.T) {
	tests := []struct {
		src  float64
		want string
		err  string
	}{
		{0, "0", ""},
		{1.5, "1.5", ""},
		{math.SmallestNonzeroFloat64, "5e-324", ""},
		{math.MaxFloat64, "1.7976931348623157e+308", ""},
		{-math.MaxFloat64, "-1.7976931348623157e+308", ""},
	}

	for _, te := range tests {
		var (
			d      mysql
			w      strings.Builder
			goterr string
			got    string
		)

		name := fmt.Sprintf("WriteFloat64(%g)", te.src)
		err := d.WriteFloat64(&w, te.src)
		if err != nil {
			goterr = err.Error()
		} else {
			got = w.String()
		}
		if goterr != te.err {
			t.Errorf("%s errored %q, want %q", name, goterr, te.err)
		}
		if got != te.want {
			t.Errorf("%s returned %q, want %q", name, got, te.want)
		}
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		src  bool
		want string
		err  string
	}{
		{false, "FALSE", ""},
		{true, "TRUE", ""},
	}

	for _, te := range tests {
		var (
			d      mysql
			w      strings.Builder
			goterr string
			got    string
		)

		name := fmt.Sprintf("WriteBool(%t)", te.src)
		err := d.WriteBool(&w, te.src)
		if err != nil {
			goterr = err.Error()
		} else {
			got = w.String()
		}
		if goterr != te.err {
			t.Errorf("%s errored %q, want %q", name, goterr, te.err)
		}
		if got != te.want {
			t.Errorf("%s returned %q, want %q", name, got, te.want)
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		src  string
		want string
		err  string
	}{
		{"", "''", ""},
		{"foo", "'foo'", ""},
		{"foo's", `'foo\'s'`, ""},
		{`foo\s`, `'foo\\s'`, ""},
		{"foo\000bar", "'foobar'", ""},
	}

	for _, te := range tests {
		var (
			d      mysql
			w      strings.Builder
			goterr string
			got    string
		)

		name := fmt.Sprintf("WriteString(%q)", te.src)
		err := d.WriteString(&w, te.src)
		if err != nil {
			goterr = err.Error()
		} else {
			got = w.String()
		}
		if goterr != te.err {
			t.Errorf("%s errored %q, want %q", name, goterr, te.err)
		}
		if got != te.want {
			t.Errorf("%s returned %q, want %q", name, got, te.want)
		}
	}
}

func TestBytes(t *testing.T) {
	tests := []struct {
		src  []byte
		want string
		err  string
	}{
		{nil, "NULL", ""},
		{[]byte{0x41, 0x42, 0x43}, "X'414243'", ""},
	}

	for _, te := range tests {
		var (
			d      mysql
			w      strings.Builder
			goterr string
			got    string
		)

		name := fmt.Sprintf("WriteBytes(%#v)", te.src)
		err := d.WriteBytes(&w, te.src)
		if err != nil {
			goterr = err.Error()
		} else {
			got = w.String()
		}
		if goterr != te.err {
			t.Errorf("%s errored %q, want %q", name, goterr, te.err)
		}
		if got != te.want {
			t.Errorf("%s returned %q, want %q", name, got, te.want)
		}
	}
}

func TestTime(t *testing.T) {
	tests := []struct {
		src  time.Time
		want string
		err  string
	}{
		{time.Time{}, "'0001-01-01 00:00:00'", ""},
		// 仕様: ロケール情報は失われます
		{time.Date(2001, time.January, 2, 13, 14, 15, 678901000, time.UTC), "'2001-01-02 13:14:15.678901'", ""},
		{time.Date(2001, time.January, 2, 13, 14, 15, 678901000, time.Local), "'2001-01-02 13:14:15.678901'", ""},
	}

	for _, te := range tests {
		var (
			d      mysql
			w      strings.Builder
			goterr string
			got    string
		)

		name := fmt.Sprintf("WriteTime(%#v)", te.src)
		err := d.WriteTime(&w, te.src)
		if err != nil {
			goterr = err.Error()
		} else {
			got = w.String()
		}
		if goterr != te.err {
			t.Errorf("%s errored %q, want %q", name, goterr, te.err)
		}
		if got != te.want {
			t.Errorf("%s returned %q, want %q", name, got, te.want)
		}
	}
}
