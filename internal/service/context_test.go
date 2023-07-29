package service

import (
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
		return ParamNumError.Fill(0, len(in))
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

func initContext(name string, args ...any) (*context, error) {
	h := &testHandler{}
	wg := sync.WaitGroup{}
	c := &context{name: name, args: args, wg: &wg, Handler: h}
	err := c.start()
	if err == nil {
		wg.Add(1)
	}
	return c, err
}

func TestContextStart(t *testing.T) {
	c, err := initContext("test")
	if err != nil {
		t.Fatal(err)
	}

	c, err = initContext("test", 1)
	err = c.start()
	if err == nil {
		t.Fatal("test2 should not start successfully")
	}
}

func TestContextStop(t *testing.T) {
	c, err := initContext("test")
	if err != nil {
		t.Fatal(err)
	}

	if !c.status() {
		t.Fatal("test should be running")
	}

	c.stop()
	time.Sleep(1 * time.Second)

	if c.status() {
		t.Fatal("test should not be running")
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

	for i := 0; i < msgCap; i++ {
		c.send(Msg{
			Cmd: "TestLoad",
		})
	}

	time.Sleep(1 * time.Second)
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
