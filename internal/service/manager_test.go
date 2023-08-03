package service

import (
	"context"
	"testing"
)

func reInit() {
	m = &Mgr{serviceByName: make(map[string]*service)}
	m.ctx, m.cancelAll = context.WithCancel(context.Background())
}

func TestServiceRegister(t *testing.T) {
	reInit()
	// test register success
	err := Register("Test", &testHandler{})
	if err != nil {
		t.Fatal(err)
	}

	// test repeated registration
	err = Register("Test", &testHandler{})
	if err == nil {
		t.Fatal("Test service should not be registered successfully")
	}

	// test register with err args
	err = Register("Test2", &testHandler{}, 1)
	if err == nil {
		t.Fatal("Test2 service should not be registered successfully")
	}
}

func TestServiceDestroy(t *testing.T) {
	reInit()
	err := Register("Test", &testHandler{})
	if err != nil {
		t.Fatal(err)
	}

	Destroy("Test")

	err = Register("Test", &testHandler{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestServiceStop(t *testing.T) {
	reInit()
	err := Register("Test", &testHandler{})
	if err != nil {
		t.Fatal(err)
	}

	Stop()

	err = Register("Test", &testHandler{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestServiceSend(t *testing.T) {
	reInit()
	err := Register("Test", &testHandler{})
	if err != nil {
		t.Fatal(err)
	}

	Send("TestLoad")
	Send("Test.TestLoad")
	Send("Test2.TestLoad")
}

func TestServiceCall(t *testing.T) {
	reInit()
	err := Register("Test", &testHandler{})
	if err != nil {
		t.Fatal(err)
	}

	err = Call("TestFunc")
	if err == nil {
		t.Fatal(err)
	}

	err = Call("Test2.TestFunc")
	if err == nil {
		t.Fatal(err)
	}

	out1 := 0.0
	out2 := 0
	err = Call("Test.TestFunc", 1, 2.0, &out1, &out2)
	if err != nil {
		t.Fatal(err)
	}
}
