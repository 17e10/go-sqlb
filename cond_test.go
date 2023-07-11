package sqlb

import (
	"bytes"
	"testing"
)

func TestCond(t *testing.T) {
	w := &bytes.Buffer{}

	joined := And(StringSqler("cond1"), StringSqler("cond2"))
	joined.Sql(w)
	got := w.String()
	want := "cond1 AND cond2"
	if got != want {
		t.Errorf("%s = %q, want %q", "test Cond #1", got, want)
	}

	w.Reset()

	joined = Bracket(
		Or(StringSqler("cond1"), StringSqler("cond2"), StringSqler("cond3")),
	)
	joined.Sql(w)
	got = w.String()
	want = "(cond1 OR cond2 OR cond3)"
	if got != want {
		t.Errorf("%s = %q, want %q", "test Cond #2", got, want)
	}
}
