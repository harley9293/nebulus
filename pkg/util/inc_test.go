package util

import "testing"

func TestAutoInc(t *testing.T) {
	ai := NewAutoInc(1, 5)
	defer ai.Close()

	if ai.Id() != 1 {
		t.Error("ai.Id() != 1")
	}

	if ai.Id() != 6 {
		t.Error("ai.Id() != 6")
	}
}
