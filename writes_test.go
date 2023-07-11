package sqlb

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"
)

func TestWriteIdent(t *testing.T) {
	tests := []struct {
		src  string
		want string
		err  string
	}{
		{"field", "`field`", ""},
		{"T.field", "`T`.`field`", ""},
		{"*", "", `not allowed asterisk`},
	}

	for _, te := range tests {
		var (
			goterr string
			got    string
		)

		name := fmt.Sprintf("writeIdent(%v)", te.src)

		w := strings.Builder{}
		err := writeIdent(&w, te.src)
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

func TestWriteValue(t *testing.T) {
	tests := []struct {
		src  any
		want string
		err  string
	}{
		{nil, "NULL", ""},
		{123, "123", ""},
		{uint64(123), "123", ""},
		{uint64(math.MaxUint64), "'18446744073709551615'", ""},
		{1.5, "1.5", ""},
		{false, "FALSE", ""},
		{true, "TRUE", ""},
		{"foo's", `'foo\'s'`, ""},
		{[]byte{0x41, 0x42, 0x43}, "X'414243'", ""},
		{time.Time{}, "'0001-01-01 00:00:00'", ""},
		{testValuer("valuer"), "'valuer'", ""},
		{&kenny, "", `got type *sqlb.person: no value type`},
		{errValuer("error"), "", `string.Value() failed: error`},
	}

	for _, te := range tests {
		var (
			goterr string
			got    string
		)

		name := fmt.Sprintf("writeValue(%v)", te.src)
		w := strings.Builder{}
		err := writeValue(&w, te.src)
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
