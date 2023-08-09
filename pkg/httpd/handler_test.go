package httpd

import "testing"

func TestHandler_Success(t *testing.T) {
	mng := newHandlerMng()
	err := mng.add("GET", "/hello", func(req *EmptyTestStruct, ctx *Context) string {
		return ""
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandler_HandlerTypeError(t *testing.T) {
	mng := newHandlerMng()
	err := mng.add("GET", "/hello", 1)
	if err == nil {
		t.Fatal("add handler error")
	}
}

func TestHandler_HandlerParamPointerError(t *testing.T) {
	mng := newHandlerMng()
	err := mng.add("GET", "/hello", func(req EmptyTestStruct, ctx Context) string {
		return ""
	})
	if err == nil {
		t.Fatal("add handler error")
	}
}

func TestHandler_HandlerSecondParamTypeError(t *testing.T) {
	mng := newHandlerMng()
	err := mng.add("GET", "/hello", func(req *EmptyTestStruct, ctx *EmptyTestStruct) string {
		return ""
	})
	if err == nil {
		t.Fatal("add handler error")
	}
}

func TestHandler_HandlerRepeatedError(t *testing.T) {
	mng := newHandlerMng()
	err := mng.add("GET", "/hello", func(req *EmptyTestStruct, ctx *Context) string {
		return ""
	})
	if err != nil {
		t.Fatal(err)
	}

	err = mng.add("GET", "/hello", func(req *EmptyTestStruct, ctx *Context) string {
		return ""
	})
	if err == nil {
		t.Fatal("add handler error")
	}
}
