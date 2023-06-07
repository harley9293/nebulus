package service

import "reflect"

type Rsp struct {
	err  error
	data []reflect.Value
}

type Msg struct {
	Cmd   string
	Req   any
	InOut []any
	Sync  bool
	Done  chan Rsp
}

type Handler interface {
	OnInit(args ...any) error
	OnStop()
	OnTick()
}

type DefaultHandler struct {
}

func (h *DefaultHandler) OnInit(args ...any) error {
	return nil
}

func (h *DefaultHandler) OnStop() {
}

func (h *DefaultHandler) OnTick() {
}
