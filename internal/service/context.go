package service

import (
	"reflect"
	"sync"
	"time"

	log "github.com/harley9293/blotlog"
	"github.com/harley9293/nebulus/internal/exception"
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

type context struct {
	name    string          // service name
	args    []any           // service startup parameters for automatic recovery
	ch      chan Msg        // message queue
	wg      *sync.WaitGroup // coroutine wait structure
	running bool            // running status

	def.Handler // service handle
}

func (c *context) status() bool {
	return c.running
}

func (c *context) start() error {
	c.ch = make(chan Msg, msgCap)
	err := c.OnInit(c.args...)
	if err != nil {
		return err
	}
	c.running = true
	log.Info("%s service started successfully", c.name)
	go c.run()
	return nil
}

func (c *context) stop() {
	if c.running == true {
		close(c.ch)
	}
	c.running = false
}

func (c *context) run() {
	defer c.stop()
	defer c.wg.Done()
	defer exception.TryE()
	for c.running {
		select {
		case msg, ok := <-c.ch:
			if ok {
				c.rev(msg)
			}
		case <-time.After(16 * time.Millisecond):
			c.OnTick()
		}
	}

	c.OnStop()
	log.Info("%s service exited normally", c.name)
}

func (c *context) rev(msg Msg) {
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

func (c *context) rsp(err error, data []reflect.Value, msg Msg) {
	if msg.Sync {
		log.Info("%s service returned %s message response, err:%s, data:%v", c.name, msg.Cmd, err, data)
		msg.Done <- Rsp{err: err, data: data}
	}
}

func (c *context) parse(f reflect.Value, params []any) (in []reflect.Value, err error) {
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

func (c *context) send(msg Msg) {
	c.ch <- msg
	if len(c.ch) > int(float64(cap(c.ch))*warnLimit) {
		log.Warn("service load is too high, name: %s, cur: %d, total: %d", c.name, len(c.ch), cap(c.ch))
	}
}

func (c *context) call(msg Msg) ([]reflect.Value, error) {
	c.send(msg)

	rsp := <-msg.Done
	return rsp.data, rsp.err
}
