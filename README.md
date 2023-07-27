# nebulus

![](https://github.com/harley9293/nebulus/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/harley9293/nebulus/branch/master/graph/badge.svg?token=UB9yvfrUP9)](https://codecov.io/gh/harley9293/nebulus)

A simple and easy-to-use distributed universal server framework.

## Installation

Install nebulus using the go get command:

```shell
go get -u github.com/harley9293/nebulus
```

## Usage

First, import the blotlog library:

```go
import "github.com/harley9293/nebulus"
```

You need to create a service that implements the Handler interface and some business interfaces.

```go
type Handler interface {
	OnInit(args ...any) error
	OnStop()
	OnTick()
}

type EchoService struct {
}

func (m *EchoService) OnInit(args ...any) error {
	return nil
}

func (m *EchoService) OnStop() {
}

func (m *EchoService) OnTick() {
}

func (m *EchoService) Business(req string) string {
	return req
}
```

If your service is very simple, you can also inherit the predefined DefaultHandler and just implement the business interface.

```go
import "github.com/harley9293/nebulus/service"

type EchoService struct {
    service.DefaultHandler
}

func (m *EchoService) Business(req string) string {
	return req
}
```

Finally, register the service and start nebulus

```go
nebulus.Register("Echo", &EchoService{})
nebulus.Run()
```

nebulus supports synchronous and asynchronous method invocation for service interfaces

```go
Send("ServiceName.FuncName", 1, 2, "3")
Call("ServiceName.FuncName", 1, 2, "3", &rsp1, &rsp2)
```

## Example

You can easily create a built-in HTTP service

```go
package http

import (
	"github.com/harley9293/nebulus"
	"github.com/harley9293/nebulus/pkg/httpd"
)

type LoginReq struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type LoginRsp struct {
	Result int    `json:"result"`
	Token  string `json:"token"`
}

func HandleLoginReq(req *LoginReq, ctx *httpd.Context) LoginRsp {
	// do something
	return LoginRsp{}
}

type GetInfoReq struct {
}

type GetInfoRsp struct {
	Info string `json:"info"`
}

func HandleGetInfoReq(req *GetInfoReq, ctx *httpd.Context) GetInfoRsp {
	// do something
	return GetInfoRsp{}
}

func Example() {
	s := httpd.NewHttpService()
	if s == nil {
		panic("NewHttpService() failed")
	}

	s.AddGlobalMiddleWare(httpd.LogMW, httpd.CookieMW, httpd.CorsMW)
	s.AddHandler("POST", "/login", HandleLoginReq)
	s.AddHandler("GET", "/info", HandleGetInfoReq, httpd.AuthMW)

	err := nebulus.Register("http", s, "127.0.0.1:36000")
	if err != nil {
		panic(err)
	}

	nebulus.Run()
}

The framework automatically parses request parameters into predefined structures, eliminating the need for further parameter parsing and packaging in the business handler, allowing you to focus on business logic.
```