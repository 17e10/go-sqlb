package sqlb

import (
	"reflect"
	"sync"
)

type columnInfo struct {
	name  string
	index []int
}

// getColumnInfos の結果キャッシュ
var cicache = make(map[string][]columnInfo)

// getColumnInfos を排他制御する mutex
var cimu = sync.Mutex{}

// getColumnInfos は構造体の exported なカラム情報を返します.
func getColumnInfos[V any](v *V) []columnInfo {
	var cols []columnInfo

	cimu.Lock()
	defer cimu.Unlock()

	rt := reflect.TypeOf(v).Elem()
	if rt.Kind() != reflect.Struct {
		panic(errNoStruct)
	}

	key := rt.String()
	if cols = cicache[key]; cols != nil {
		return cols
	}

	fields := reflect.VisibleFields(rt)
	cols = make([]columnInfo, 0, len(fields))
	for _, f := range fields {
		if f.Anonymous || !f.IsExported() {
			continue
		}
		name := f.Tag.Get("sqlb")
		if name == "" {
			name = columnName(f.Name)
		}
		cols = append(cols, columnInfo{name, f.Index})
	}
	cicache[key] = cols
	return cols
}

// exportedColumns は構造体の exported なカラムから excludes を除外したカラム情報を返します.
func exportedColumns[V any](v *V, excludes []string) []columnInfo {
	cols := getColumnInfos(v)
	if len(excludes) == 0 {
		return cols
	}

	exclude := make(map[string]bool)
	for _, s := range excludes {
		exclude[s] = true
	}

	exps := make([]columnInfo, 0, len(cols))
	for _, col := range cols {
		if exclude[col.name] {
			continue
		}
		exps = append(exps, col)
	}
	return exps
}

// Columns は構造体からカラムリストを作成します.
//
// excludes でリストから除外するカラムを指定できます.
func Columns[V any](v *V, excludes ...string) []string {
	cols := exportedColumns(v, excludes)
	r := make([]string, len(cols))
	for i, c := range cols {
		r[i] = c.name
	}
	return r
}

// Values は構造体から値リストを作成します.
//
// excludes でリストから除外するカラムを指定できます.
func Values[V any](v *V, excludes ...string) []any {
	cols := exportedColumns(v, excludes)
	rv := reflect.ValueOf(v).Elem()
	r := make([]any, len(cols))
	for i, c := range cols {
		r[i] = rv.FieldByIndex(c.index).Interface()
	}
	return r
}

// GroupValues は構造体の配列からグループリストを作成します.
//
// excludes でリストから除外するカラムを指定できます.
func GroupValues[V any](v []V, excludes ...string) [][]any {
	cols := exportedColumns((*V)(nil), excludes)
	group := make([][]any, len(v))
	for i, val := range v {
		rv := reflect.ValueOf(&val).Elem()
		r := make([]any, len(cols))
		for i, c := range cols {
			r[i] = rv.FieldByIndex(c.index).Interface()
		}
		group[i] = r
	}
	return group
}

// KeyValues は構造体から Key-Value リストを作成します.
//
// excludes でリストから除外するカラムを指定できます.
func KeyValues[V any](v *V, excludes ...string) []Kv {
	cols := exportedColumns(v, excludes)
	rv := reflect.ValueOf(v).Elem()
	r := make([]Kv, len(cols))
	for i, c := range cols {
		r[i].K = c.name
		r[i].V = rv.FieldByIndex(c.index).Interface()
	}
	return r
}
