package errors

import "testing"

func TestNew(t *testing.T) {
	err := New("test")
	if err == nil {
		t.Fatal(err)
	}
}

func TestFill(t *testing.T) {
	err := New("test %d %d %d").Fill(1, 2, 3)
	if err == nil {
		t.Fatal(err)
	}

	if err.Error() != "test 1 2 3" {
		t.Fatal(err)
	}
}
