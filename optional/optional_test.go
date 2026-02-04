package optional

import "testing"

func TestOptional(t *testing.T) {
	o := Some(42)
	if !o.Set || o.Value != 42 {
		t.Fatalf("expected set value")
	}

	n := None[int]()
	if n.Set {
		t.Fatalf("expected none")
	}
}
