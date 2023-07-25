package exception

import "testing"

func TestTryE_Panic(t *testing.T) {
	defer TryE()

	panic("test panic")
}

func TestTryE_Pass(t *testing.T) {
	defer TryE()
}
