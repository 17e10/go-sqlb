package sqlb

import (
	"testing"
)

func equalPtrs(t *testing.T, name string, x, y []any) {
	if len(x) != len(y) {
		t.Errorf("%s len %d, want %d", name, len(x), len(y))
		return
	}
	for i, l := 0, len(x); i < l; i++ {
		if x[i] != y[i] {
			t.Errorf("%s index %d %p, want %p", name, i, x[i], y[i])
		}
	}
}

func TestMakeFieldsAddr(t *testing.T) {
	type users struct {
		FirstName string
		LastName  string
		fullname  string
	}

	type data struct {
		users
		Age int
		gen int
	}

	// note: use fullname and gen to avoid warnings
	record := data{
		users: users{
			FirstName: "",
			LastName:  "",
			fullname:  "",
		},
		Age: 0,
		gen: 0,
	}
	result := makeFieldsAddr(&record)
	want := []any{&record.FirstName, &record.LastName, &record.Age}
	equalPtrs(t, "test makeFieldsAddr", result, want)
}

func BenchmarkMakeFieldsAddr(b *testing.B) {
	dest := struct {
		FirstName string
		LastName  string
		Age       int
	}{}

	for i := 0; i < b.N; i++ {
		makeFieldsAddr(&dest)
	}
}

// BenchmarkMakeFieldsAddr-4   	16529665	        69.59 ns/op	      48 B/op	       1 allocs/op
