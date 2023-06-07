# nebulus

![](https://github.com/harley9293/nebulus/workflows/Go/badge.svg)

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