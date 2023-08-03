package service

import (
	"context"
	log "github.com/harley9293/blotlog"
	"github.com/harley9293/nebulus/pkg/def"
	"github.com/harley9293/nebulus/pkg/errors"
	"reflect"
	"strings"
	"sync"
)

type Rsp struct {
	err  error
	data []reflect.Value
}

type Msg struct {
	Cmd   string
	InOut []any
	Sync  bool
	Done  chan Rsp
}

var RegisterExistError = errors.New("%s service has already been registered")
var InvalidCallFuncError = errors.New("invalid call function: %s")
var NotRegisterError = errors.New("%s service is not registered")

var m *Mgr

type Mgr struct {
	serviceByName map[string]*service
	wg            sync.WaitGroup

	ctx       context.Context
	cancelAll context.CancelFunc
}

func init() {
	m = &Mgr{serviceByName: make(map[string]*service)}
	m.ctx, m.cancelAll = context.WithCancel(context.Background())
}

func Stop() {
	m.cancelAll()
	m.wg.Wait()
	m.serviceByName = make(map[string]*service)
}

func Register(name string, h def.Handler, args ...any) error {
	check(name)
	_, ok := m.serviceByName[name]
	if ok {
		return RegisterExistError.Fill(name)
	}

	c := &service{name: name, wg: &m.wg, Handler: h, ch: make(chan Msg, msgCap), exit: make(chan bool, 1)}
	c.ctx, c.cancel = context.WithCancel(m.ctx)
	err := c.OnInit(args...)
	if err != nil {
		return err
	}
	m.serviceByName[name] = c
	m.wg.Add(1)
	go c.run()
	return nil
}

func Destroy(name string) {
	c, ok := m.serviceByName[name]
	if ok {
		c.cancel()
		log.Info("%s service is deleted by destroy", name)
		delete(m.serviceByName, name)
	}
}

func Send(f string, in ...any) {
	l := strings.Split(f, ".")
	if len(l) != 2 {
		log.Warn("Invalid call object: %s, should be service.func", f)
		return
	}
	name := l[0]
	cmd := l[1]

	check(name)
	c, ok := m.serviceByName[name]
	if !ok {
		log.Warn("%s service is not registered, send failed", name)
		return
	}

	var msg Msg
	msg.Cmd = cmd
	msg.InOut = in
	c.send(msg)
}

func Call(f string, inout ...any) error {
	l := strings.Split(f, ".")
	if len(l) != 2 {
		err := InvalidCallFuncError.Fill(f)
		log.Warn(err.Error())
		return err
	}
	name := l[0]
	cmd := l[1]

	check(name)
	c, ok := m.serviceByName[name]
	if !ok {
		err := NotRegisterError.Fill(name)
		log.Warn(err.Error())
		return err
	}

	var msg Msg
	msg.Cmd = cmd
	msg.InOut = inout
	done := make(chan Rsp)
	defer close(done)
	msg.Sync = true
	msg.Done = done
	data, err := c.call(msg)
	if err == nil {
		for i := 0; i < len(data); i++ {
			v := reflect.ValueOf(inout[len(inout)-len(data)+i])
			v.Elem().Set(data[i])
		}
	}
	return err
}

func check(name string) {
	c, ok := m.serviceByName[name]
	if ok {
		select {
		case <-c.exit:
			log.Info("%s service is deleted by check", name)
			delete(m.serviceByName, name)
		default:
			break
		}
	}
}
