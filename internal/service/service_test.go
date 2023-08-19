package service

import (
	"context"
	"errors"
	"github.com/harley9293/nebulus/pkg/def"
	"sync"
	"testing"
	"time"
)

type testHandler struct {
	def.DefaultHandler
}

func (h *testHandler) OnInit(in ...any) error {
	if len(in) != 0 {
		return errors.New("invalid input")
	}
	return nil
}

func (h *testHandler) TestFunc(a int, b float64) (float64, int) {
	return b, a
}

func (h *testHandler) TestLoad() {
}

func (h *testHandler) TestPanic() {
	panic("test panic")
}

func (h *testHandler) TestTimeout() {
	time.Sleep(10 * time.Hour)
}

func initContext(name string, args ...any) (*service, error) {
	h := &testHandler{}
	wg := sync.WaitGroup{}
	ctx, _ := context.WithCancel(context.Background())
	c := &service{name: name, wg: &wg, Handler: h, ch: make(chan Msg, msgCap)}
	c.ctx, c.cancel = context.WithCancel(ctx)
	err := c.OnInit(args...)
	if err != nil {
		return nil, err
	}
	wg.Add(1)
	go c.run()
	return c, nil
}

func TestContextStart(t *testing.T) {
	_, err := initContext("test")
	if err != nil {
		t.Fatal(err)
	}

	_, err = initContext("test", 1)
	if err == nil {
		t.Fatal("test2 should not start successfully")
	}
}

func TestContextOnTick(t *testing.T) {
	_, err := initContext("test")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)
}

func TestContextLoadWarn(t *testing.T) {
	c, err := initContext("test")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < msgCap*1.2; i++ {
		go func() {
			done := make(chan Rsp)
			_, _ = c.call(Msg{
				Cmd:   "TestTimeout",
				InOut: []any{},
				Sync:  true,
				Done:  done,
			})
		}()
		go func() {
			c.send(Msg{
				Cmd:   "TestLoad",
				InOut: []any{},
			})
		}()
	}

	time.Sleep(6 * time.Second)
}

func TestContextCall(t *testing.T) {
	c, err := initContext("test")
	if err != nil {
		t.Fatal(err)
	}

	in0 := 1
	in1 := 2.0
	out0 := 0.0
	out1 := 0
	out2 := ""

	// test func success
	done := make(chan Rsp)
	data, err := c.call(Msg{
		Cmd:   "TestFunc",
		Sync:  true,
		InOut: []any{in0, in1, &out0, &out1},
		Done:  done,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 2 {
		t.Fatal("TestFunc should return 2 values")
	}
	if data[0].Float() != in1 || data[1].Int() != int64(in0) {
		t.Fatal("TestFunc return value error")
	}

	// test func param num error
	_, err = c.call(Msg{
		Cmd:   "TestFunc",
		InOut: []any{in0, in1, &out0, &out1, &out0},
		Sync:  true,
		Done:  done,
	})
	if err == nil {
		t.Fatal("TestFunc should not be called successfully")
	}

	// test func in param type mismatch
	_, err = c.call(Msg{
		Cmd:   "TestFunc",
		InOut: []any{2.0, in1, &out0, &out0},
		Sync:  true,
		Done:  done,
	})
	if err == nil {
		t.Fatal("TestFunc should not be called successfully")
	}

	// test func out param not pointer
	_, err = c.call(Msg{
		Cmd:   "TestFunc",
		InOut: []any{in0, in1, out0, &out0},
		Sync:  true,
		Done:  done,
	})
	if err == nil {
		t.Fatal("TestFunc should not be called successfully")
	}

	// test func out param type mismatch
	_, err = c.call(Msg{
		Cmd:   "TestFunc",
		InOut: []any{in0, in1, &out0, &out2},
		Sync:  true,
		Done:  done,
	})
	if err == nil {
		t.Fatal("TestFunc should not be called successfully")
	}
}

func TestCall_Timeout(t *testing.T) {
	c, err := initContext("test")
	if err != nil {
		t.Fatal(err)
	}

	// test func timeout
	done := make(chan Rsp)
	_, err = c.call(Msg{
		Cmd:   "TestTimeout",
		InOut: []any{},
		Sync:  true,
		Done:  done,
	})
	if err == nil {
		t.Fatal("TestTimeout should not be called successfully")
	}
}

func TestService_Panic(t *testing.T) {
	c, err := initContext("test")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 6; i++ {
		// test func panic
		done := make(chan Rsp)
		go func() {
			_, _ = c.call(Msg{
				Cmd:   "TestPanic",
				InOut: []any{},
				Sync:  true,
				Done:  done,
			})
		}()
	}
}
