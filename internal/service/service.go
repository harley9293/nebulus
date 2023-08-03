package service

import (
	"context"
	"github.com/harley9293/nebulus/internal/exception"
	"reflect"
	"sync"
	"time"

	log "github.com/harley9293/blotlog"
	"github.com/harley9293/nebulus/pkg/def"
	"github.com/harley9293/nebulus/pkg/errors"
)

const (
	msgCap    = 100 // message queue capacity
	warnLimit = 0.8 // warning ratio
)

var ParamNumError = errors.New("parameter number error, need:%d, got:%d")
var ParamTypeMismatch = errors.New("parameter %d type mismatch, need:%v, got:%v")
var ParamNotPointerError = errors.New("parameter %d is not a pointer")

type service struct {
	name string   // service name
	args []any    // service startup parameters for automatic recovery
	ch   chan Msg // message queue

	wg     *sync.WaitGroup    // coroutine wait structure
	ctx    context.Context    // coroutine context
	cancel context.CancelFunc // coroutine cancel function

	def.Handler // service handle
}

func (c *service) run() {
	log.Info("%s service started successfully", c.name)
	defer func() {
		p := exception.TryE()
		if p != nil {
			log.Error("%s service exited abnormally", c.name)
			c.OnPanic(p)
			// TODO try to restart the service
		} else {
			c.wg.Done()
			close(c.ch)
		}
	}()
Loop:
	for {
		select {
		case msg, ok := <-c.ch:
			if ok {
				c.rev(msg)
			}
		case <-time.After(16 * time.Millisecond):
			c.OnTick()
		case <-c.ctx.Done():
			break Loop
		}
	}

	c.OnStop()
	log.Info("%s service exited normally", c.name)
}

func (c *service) rev(msg Msg) {
	log.Info("%s service received %s message request: %v", c.name, msg.Cmd, msg.InOut)
	f := reflect.ValueOf(c.Handler).MethodByName(msg.Cmd)
	in, err := c.parse(f, msg.InOut)
	if err != nil {
		log.Error(err.Error())
		c.rsp(err, []reflect.Value{}, msg)
		return
	}
	data := f.Call(in)
	c.rsp(nil, data, msg)
}

func (c *service) rsp(err error, data []reflect.Value, msg Msg) {
	if msg.Sync {
		log.Info("%s service returned %s message response, err:%s, data:%v", c.name, msg.Cmd, err, data)
		msg.Done <- Rsp{err: err, data: data}
	}
}

func (c *service) parse(f reflect.Value, params []any) (in []reflect.Value, err error) {
	if f.Type().NumIn()+f.Type().NumOut() != len(params) {
		err = ParamNumError.Fill(f.Type().NumIn()+f.Type().NumOut(), len(params))
		return
	}

	index := 0
	for i := 0; i < f.Type().NumIn(); i++ {
		t := reflect.TypeOf(params[index])
		if f.Type().In(i) != t {
			err = ParamTypeMismatch.Fill(index, f.Type().In(i), t)
			return
		}
		index++
		in = append(in, reflect.ValueOf(params[i]))
	}

	for i := 0; i < f.Type().NumOut(); i++ {
		t := reflect.TypeOf(params[index])
		if t.Kind() != reflect.Pointer {
			err = ParamNotPointerError.Fill(index)
			return
		}
		if f.Type().Out(i) != t.Elem() {
			err = ParamTypeMismatch.Fill(index, f.Type().Out(i), t.Elem())
			return
		}
		index++
	}
	return
}

func (c *service) send(msg Msg) {
	if len(c.ch) >= cap(c.ch) {
		log.Error("%s service message queue is full, message discarded: %s, %+v", c.name, msg.Cmd, msg.InOut)
		if msg.Sync {
			msg.Done <- Rsp{err: errors.New("message queue is full")}
		}
		return
	}

	c.ch <- msg
	if len(c.ch) > int(float64(cap(c.ch))*warnLimit) {
		log.Warn("service load is too high, name: %s, cur: %d, total: %d", c.name, len(c.ch), cap(c.ch))
	}
}

func (c *service) call(msg Msg) ([]reflect.Value, error) {
	c.send(msg)

	select {
	case rsp := <-msg.Done:
		return rsp.data, rsp.err
	case <-time.After(5 * time.Second):
		return nil, errors.New("call timeout")
	}
}
