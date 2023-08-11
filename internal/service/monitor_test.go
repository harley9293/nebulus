package service

import (
	"testing"
	"time"
)

func TestMonitorLoop(t *testing.T) {
	reInit()
	go monitorLoop(m.ctx, &m.wg)

	_ = Register("Test", &testHandler{})

	time.Sleep(100 * time.Millisecond)

	for i := 0; i < 6; i++ {
		go func() {
			_ = Call("Test.TestPanic")
		}()
	}

	time.Sleep(100 * time.Millisecond)
}
