package httpd

import "testing"

func SuccessFuncTest(req *LoginReq, ctx *Context) LoginRsp {
	return LoginRsp{}
}

func FailFuncTest1(req *LoginReq) {
}

func FailFuncTest2(req LoginReq, ctx *Context) LoginRsp {
	return LoginRsp{}
}

func FailFuncTest3(req *LoginReq, ctx *LoginReq) LoginRsp {
	return LoginRsp{}
}

func FailFuncTest4(req *LoginReq, ctx *Context) func() {
	return func() {}
}

func TestHandler(t *testing.T) {
	mng := newHandlerMng()

	err := mng.add("GET", "/hello", SuccessFuncTest)
	if err != nil {
		t.Fatal(err)
	}

	err = mng.add("GET", "/hello", SuccessFuncTest)
	if err == nil {
		t.Fatal(err)
	}

	err = mng.add("GET", "/test", FailFuncTest1)
	if err == nil {
		t.Fatal("add handler error")
	}

	err = mng.add("GET", "/test", FailFuncTest2)
	if err == nil {
		t.Fatal("add handler error")
	}

	err = mng.add("GET", "/test", FailFuncTest3)
	if err == nil {
		t.Fatal("add handler error")
	}

	err = mng.add("GET", "/test", 1)
	if err == nil {
		t.Fatal("add handler error")
	}

	err = mng.add("GET", "/test", FailFuncTest4)
	if err == nil {
		t.Fatal("add handler error")
	}
}
