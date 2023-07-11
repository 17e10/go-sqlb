package sqlb

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/17e10/go-patb"
)

// Kv は Key-Value ペア `key` = val を表す構造体です.
type Kv struct {
	K string
	V any
}

// putIdent はプレースホルダの識別子を展開します.
func putIdent(w Writer, v any) (err error) {
	switch val := v.(type) {
	case string:
		err = writeIdent(w, val)
	case []string:
		err = putIdentList(w, val)
	default:
		err = fmt.Errorf("got type %T: %w", v, errNoIdentType)
	}
	if err != nil {
		return fmt.Errorf("ident: %v: %w", v, err)
	}
	return nil
}

// putIdent はプレースホルダの識別子リスト `field1`, `field2`, ... を展開します.
func putIdentList(w Writer, v []string) error {
	if len(v) == 0 {
		return errEmptySlice
	}
	for i, s := range v {
		if i > 0 {
			w.WriteString(", ")
		}
		if err := writeIdent(w, s); err != nil {
			return err
		}
	}
	return nil
}

// putValue はプレースホルダの値を展開します.
func putValue(w Writer, v any) (err error) {
	switch val := v.(type) {
	case []any:
		err = putValueList(w, val)
	case [][]any:
		err = putGroupList(w, val)
	case []Kv:
		err = putKvList(w, val)
	default:
		err = writeValue(w, v)
	}
	if err != nil {
		return fmt.Errorf("value: %v: %w", v, err)
	}
	return nil
}

// putValueList はプレースホルダの値リスト value1, value2, ... を展開します.
func putValueList(w Writer, v []any) error {
	if len(v) == 0 {
		return errEmptySlice
	}
	for i, val := range v {
		if i > 0 {
			w.WriteString(", ")
		}
		if err := writeValue(w, val); err != nil {
			return err
		}
	}
	return nil
}

// putGroupList はプレースホルダのグループリスト
// (value1, value2, ...), (value3, value4, ...) を展開します.
func putGroupList(w Writer, v [][]any) error {
	if len(v) == 0 {
		return errEmptySlice
	}
	for i, list := range v {
		if i == 0 {
			w.WriteByte('(')
		} else {
			w.WriteString(", (")
		}
		if err := putValueList(w, list); err != nil {
			return err
		}
		w.WriteByte(')')
	}
	return nil
}

// putKvList はプレースホルダの Key-Value ペアリスト
// `key1` = value1, `key2` = value2, ...を展開します.
func putKvList(w Writer, v []Kv) error {
	if len(v) == 0 {
		return errEmptySlice
	}
	for i, kv := range v {
		if i > 0 {
			w.WriteString(", ")
		}
		if err := writeIdent(w, kv.K); err != nil {
			return err
		}
		w.WriteString(" = ")
		if err := writeValue(w, kv.V); err != nil {
			return err
		}
	}
	return nil
}

// putSqler は Sqler を展開します.
func putSqler(w Writer, v any) error {
	sqler, ok := v.(Sqler)
	if !ok {
		return fmt.Errorf("got type %T, want Sqler: %w", v, errNoValueType)
	}
	return sqler.Sql(w)
}

// putEqValue は擬似イコール構文を含めた値を展開します.
func putEqValue(w Writer, eq string, v any) error {
	switch val := v.(type) {
	case nil:
		return putIsNull(w, eq)
	case []any:
		switch len(val) {
		case 0:
			return errEmptySlice
		case 1:
			return putEqual(w, eq, val[0])
		}
		return putIn(w, eq, val)
	}
	return putEqual(w, eq, v)
}

// putIsNull は擬似イコール構文を IS (NOT) NULL に展開します.
func putIsNull(w Writer, eq string) error {
	if eq == "==" {
		w.WriteString("IS NULL")
	} else {
		w.WriteString("IS NOT NULL")
	}
	return nil
}

// putIn は擬似イコール構文を (NOT) IN (...) に展開します.
func putIn(w Writer, eq string, v []any) error {
	if eq == "==" {
		w.WriteString("IN (")
	} else {
		w.WriteString("NOT IN (")
	}
	if err := putValueList(w, v); err != nil {
		return err
	}
	w.WriteByte(')')
	return nil
}

// putEqual は擬似イコール構文を = 値, != 値 に展開します.
func putEqual(w Writer, eq string, v any) error {
	if eq == "==" {
		w.WriteString("= ")
	} else {
		w.WriteString("!= ")
	}
	return writeValue(w, v)
}

// splitEqValue は擬似イコール構文をイコールとプレースホルダに分離します.
func splitEqValue(s string) (eq string, v string) {
	const space = " \t\n\r"
	if len(s) > 3 {
		if s[:3] == "!==" {
			return s[:3], strings.TrimLeft(s[3:], space)
		}
		if s[:2] == "==" {
			return s[:2], strings.TrimLeft(s[2:], space)
		}
	}
	return "", s
}

// T はプレースホルダを展開します.
//
// T は比較的短い SQL を簡潔に記述するための関数です.
// tmpl には次の構文が利用できます.
//
//	$(index) // Sqler
//	#(index) // 識別子
//	@(index) // 値
//
// index は省略可能で 省略した場合は引数の先頭から順に展開します.
// index を指定した場合は指定した位置の引数から順に展開します.
func T(tmpl string, a ...any) Sqler {
	return &texec{tmpl, a, 0}
}

// texec は T の Sqler を実装します.
type texec struct {
	tmpl string
	a    []any
	i    int
}

var tc = patb.C("@#$=!")
var tpat = patb.Any(
	patb.Block(patb.S("@"), patb.Ch(0, 2, patb.Digit())),
	patb.Block(patb.S("#"), patb.Ch(0, 2, patb.Digit())),
	patb.Block(patb.S("$"), patb.Ch(0, 2, patb.Digit())),
	patb.Block(patb.Ch(0, 1, patb.C("!")), patb.S("=="), patb.Ch(1, 16, patb.Space()), patb.S("@"), patb.Ch(0, 2, patb.Digit())),
)

// Sql はプレースホルダを展開します.
func (te *texec) Sql(w Writer) error {
	te.i = 0
	return patb.ReplaceWrite(w, tc, tpat, te.tmpl, func(w patb.Writer, m string) error {
		return te.extract(w, m)
	})
}

// getParam は s に対応するパラメータを返します.
func (te *texec) getParam(s string) (int, any, error) {
	if len(s) > 1 {
		te.i, _ = strconv.Atoi(s[1:])
	}
	if te.i < 0 || te.i >= len(te.a) {
		return te.i, nil, errOutRange
	}

	v := te.a[te.i]
	te.i++
	return te.i - 1, v, nil
}

// extract は T がサポートする構文に応じてパラメータを展開します.
func (te *texec) extract(w Writer, m string) error {
	eq, m := splitEqValue(m)
	i, v, err := te.getParam(m)
	if err == nil {
		switch {
		case m[0] == '$':
			err = putSqler(w, v)
		case m[0] == '#':
			err = putIdent(w, v)
		case m[0] == '@' && eq == "":
			err = putValue(w, v)
		default:
			err = putEqValue(w, eq, v)
		}
	}
	if err != nil {
		return fmt.Errorf("T extract %s, index %d: %w", m, i, err)
	}
	return nil
}

// M は名前付きプレースホルダを展開します.
//
// M は比較的長い SQL を記述するための関数です.
// tmpl には次の構文が利用できます.
//
//	$sqler_name // Sqler
//	#ident_name // 識別子
//	@value_name // 値
func M(tmpl string, params map[string]any) Sqler {
	return &mexec{tmpl, params}
}

// mexec は M の Sqler を実装します.
type mexec struct {
	tmpl   string
	params map[string]any
}

var mc = patb.C("@#$=!")
var mpat = patb.Any(
	patb.Block(patb.S("@"), patb.Ch(1, 64, patb.Word())),
	patb.Block(patb.S("#"), patb.Ch(1, 64, patb.Word())),
	patb.Block(patb.S("$"), patb.Ch(1, 64, patb.Word())),
	patb.Block(patb.Ch(0, 1, patb.C("!")), patb.S("=="), patb.Ch(1, 16, patb.Space()), patb.S("@"), patb.Ch(1, 64, patb.Word())),
)

// Sql は名前付きプレースホルダを展開します.
func (me *mexec) Sql(w Writer) error {
	return patb.ReplaceWrite(w, mc, mpat, me.tmpl, func(w patb.Writer, m string) error {
		return me.extract(w, m)
	})
}

// getParam は k に対応するパラメータを返します.
func (me *mexec) getParam(k string) (string, any, error) {
	v, has := me.params[k]
	if !has {
		return k, nil, errNoSuchKey
	}
	return k, v, nil
}

// extract は M がサポートする構文に応じてパラメータを展開します.
func (me *mexec) extract(w Writer, m string) error {
	eq, m := splitEqValue(m)
	_, v, err := me.getParam(m)
	if err == nil {
		switch {
		case m[0] == '$':
			err = putSqler(w, v)
		case m[0] == '#':
			err = putIdent(w, v)
		case m[0] == '@' && eq == "":
			err = putValue(w, v)
		default:
			err = putEqValue(w, eq, v)
		}
	}
	if err != nil {
		return fmt.Errorf("M extract %s: %w", m, err)
	}
	return nil
}
