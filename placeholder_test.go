package sqlb

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/17e10/go-sqlb/sqlt"
)

func TestPutIdent(t *testing.T) {
	tests := []struct {
		src  any
		want string
		err  string
	}{
		{"a", "`a`", ""},
		{[]string{"a", "b"}, "`a`, `b`", ""},
		{"*", "", `ident: *: not allowed asterisk`},
		{123, "", `ident: 123: got type int: no ident type`},
		{[]string{}, "", `ident: []: empty array or slice`},
		{[]string{"a", ""}, "", `ident: [a ]: empty ident`},
	}

	w := &bytes.Buffer{}
	for _, te := range tests {
		var (
			goterr string
			got    string
		)
		w.Reset()
		name := fmt.Sprintf("putIdent(%#v)", te.src)
		err := putIdent(w, te.src)
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

func TestPutValue(t *testing.T) {
	tests := []struct {
		src  any
		want string
		err  string
	}{
		// primitive
		{"a", "'a'", ""},

		// value list
		{[]any{"a", "b"}, "'a', 'b'", ""},
		{[]any{}, "", `value: []: empty array or slice`},
		{[]any{noint(0)}, "", `value: [0]: got type sqlb.noint: no value type`},

		// group list
		{[][]any{{"a", "b"}, {"c", "d"}}, "('a', 'b'), ('c', 'd')", ""},
		{[][]any{}, "", `value: []: empty array or slice`},
		{[][]any{{noint(0)}}, "", `value: [[0]]: got type sqlb.noint: no value type`},

		// key-value list
		{[]Kv{{"k1", "v1"}, {"k2", "v2"}}, "`k1` = 'v1', `k2` = 'v2'", ""},
		{[]Kv{}, "", `value: []: empty array or slice`},
		{[]Kv{{"k", noint(0)}}, "", `value: [{k 0}]: got type sqlb.noint: no value type`},
	}

	w := &bytes.Buffer{}
	for _, te := range tests {
		var (
			goterr string
			got    string
		)

		w.Reset()
		name := fmt.Sprintf("putValue(%#v)", te.src)
		err := putValue(w, te.src)
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

func TestPutSqler(t *testing.T) {
	tests := []struct {
		src  any
		want string
		err  string
	}{
		{StringSqler("abc"), "abc", ""},
		{nil, "", "got type <nil>, want Sqler: no value type"},
	}

	w := &bytes.Buffer{}
	for _, te := range tests {
		var (
			goterr string
			got    string
		)

		w.Reset()
		name := fmt.Sprintf("putSqler(%#v)", te.src)
		err := putSqler(w, te.src)
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

func TestEqValue(t *testing.T) {
	tests := []struct {
		eq   string
		v    any
		want string
		err  string
	}{
		// primitive
		{"==", 0, "= 0", ""},
		{"!==", 0, "!= 0", ""},

		// nil
		{"==", nil, "IS NULL", ""},
		{"!==", nil, "IS NOT NULL", ""},

		// value list
		{"==", []any{}, "", "empty array or slice"},
		{"==", []any{"a"}, "= 'a'", ""},
		{"==", []any{"a", "b"}, "IN ('a', 'b')", ""},
		{"!==", []any{"a", "b"}, "NOT IN ('a', 'b')", ""},
		{"==", []string{"a", "b"}, "", `got type []string: no value type`},
	}

	w := &bytes.Buffer{}
	for _, te := range tests {
		var (
			goterr string
			got    string
		)

		w.Reset()
		name := fmt.Sprintf("putEqValue(%q, %#v)", te.eq, te.v)
		err := putEqValue(w, te.eq, te.v)
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

func TestT(t *testing.T) {
	tests := []struct {
		s    string
		a    []any
		want string
		err  string
	}{
		// basis test
		{"a $ b", []any{StringSqler("and")}, "a and b", ""},
		{"`age` = @", []any{20}, "`age` = 20", ""},
		{"# = @", []any{"age", 20}, "`age` = 20", ""},
		{"# == @", []any{"age", nil}, "`age` IS NULL", ""},
		{"# !== @", []any{"age", []any{123, 234}}, "`age` NOT IN (123, 234)", ""},
		{"#0 <= @1 AND #0 < @2", []any{"age", 20, 30}, "`age` <= 20 AND `age` < 30", ""},
		{"#1 = @", []any{"nop", "age", 20}, "`age` = 20", ""},
		{"#2 = @", []any{"age", 20}, "", "T extract #2, index 2: out of range"},
	}

	w := &bytes.Buffer{}
	for i, te := range tests {
		var (
			goterr string
			got    string
		)

		w.Reset()
		name := fmt.Sprintf("test T #%d", i)
		err := T(te.s, te.a...).Sql(w)
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

func TestM(t *testing.T) {
	tests := []struct {
		s      string
		params map[string]any
		want   string
		err    string
	}{
		{"a $op b", map[string]any{"$op": StringSqler("and")}, "a and b", ""},
		{"#field = @value", map[string]any{"#field": "age", "@value": 20}, "`age` = 20", ""},
		{"#field == @value", map[string]any{"#field": "age", "@value": nil}, "`age` IS NULL", ""},
		{"#noname = @value", map[string]any{"#field": "age", "@value": 20}, "", `M extract #noname: no such key`},
	}

	w := &bytes.Buffer{}
	for i, te := range tests {
		var (
			goterr string
			got    string
		)

		w.Reset()
		name := fmt.Sprintf("test M #%d", i)
		err := M(te.s, te.params).Sql(w)
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

func BenchmarkInsert(b *testing.B) {
	ctx := context.TODO()
	execer := sqlt.NullExecer{}
	// cols := []string{"family_name", "given_name", "age"}
	// values := []any{"Williams", "Olivia", 54}

	for i := 0; i < b.N; i++ {
		sqler := T(
			"INSERT INTO person (#) VALUES (@)",
			Columns(&olivia, "id"),
			Values(&olivia, "id"),
		)
		Exec(execer, ctx, sqler)
	}
}

// hard coding: BenchmarkInsert-4   	 1455375	       808.0 ns/op	     904 B/op	      15 allocs/op
// use pat:     BenchmarkInsert-4   	 1382932	       860.6 ns/op	     872 B/op	      13 allocs/op
//
// hard coding vs use pat ... 処理速度 1:1.07, メモリ 1:0.96, allocs 1:0.87
