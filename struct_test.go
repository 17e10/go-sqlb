package sqlb

import (
	"reflect"
	"testing"
)

func TestColumns(t *testing.T) {
	got := Columns((*person)(nil))
	want := []string{"id", "family_name", "given_name", "age"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s = %v, want %v", "test Columns #1", got, want)
	}

	got = Columns((*person)(nil), "id")
	want = []string{"family_name", "given_name", "age"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s = %v, want %v", "test Columns #2", got, want)
	}

	got = Columns((*person)(nil), "id", "age")
	want = []string{"family_name", "given_name"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s = %v, want %v", "test Columns #3", got, want)
	}
}

func TestValues(t *testing.T) {
	got := Values(&olivia)
	want := []any{uint64(1), "Williams", "Olivia", 54}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s = %v, want %v", "test Values #1", got, want)
	}

	got = Values(&olivia, "id")
	want = []any{"Williams", "Olivia", 54}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s = %v, want %v", "test Values #2", got, want)
	}

	got = Values(&olivia, "id", "age")
	want = []any{"Williams", "Olivia"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s = %v, want %v", "test Values #3", got, want)
	}
}

func TestGroupValues(t *testing.T) {
	got := GroupValues(persons)
	want := [][]any{
		{uint64(1), "Williams", "Olivia", 54},
		{uint64(2), "Loggins", "Kenny", 75},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s = %v, want %v", "test GroupValues #1", got, want)
	}

	got = GroupValues(persons, "id")
	want = [][]any{
		{"Williams", "Olivia", 54},
		{"Loggins", "Kenny", 75},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s = %v, want %v", "test GroupValues #2", got, want)
	}
}

func TestKeyValues(t *testing.T) {
	got := KeyValues(&olivia)
	want := []Kv{
		{"id", uint64(1)},
		{"family_name", "Williams"},
		{"given_name", "Olivia"},
		{"age", 54},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s = %v, want %v", "test KeyValues #1", got, want)
	}

	got = KeyValues(&olivia, "id")
	want = []Kv{
		{"family_name", "Williams"},
		{"given_name", "Olivia"},
		{"age", 54},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s = %v, want %v", "test KeyValues #2", got, want)
	}
}
