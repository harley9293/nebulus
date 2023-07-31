package util

import (
	"testing"
	"time"
)

func TestSpinLock(t *testing.T) {
	lock := NewSpinLock()

	lock.Lock()
	go func() {
		lock.Lock()
		lock.Unlock()
	}()
	time.Sleep(1 * time.Second)
	lock.Unlock()
}
