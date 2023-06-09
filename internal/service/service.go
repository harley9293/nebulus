package service

import (
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
	Req   any
	InOut []any
	Sync  bool
	Done  chan Rsp
}

var RegisterExistError = errors.New("%s service has already been registered")
var InvalidCallFuncError = errors.New("invalid call function: %s")
var NotRegisterError = errors.New("%s service is not registered")

var m *mgr

type mgr struct {
	serviceByName map[string]*context
	wg            sync.WaitGroup
}

func init() {
	m = new(mgr)
	m.serviceByName = make(map[string]*context)
}

func Stop() {
	for _, v := range m.serviceByName {
		v.stop()
	}

	m.wg.Wait()
	m.serviceByName = make(map[string]*context)
}

func Tick() {
	for _, v := range m.serviceByName {
		if !v.status() {
			log.Warn("%s Service is attempting to restart", v.name)
			err := v.start()
			if err != nil {
				log.Error("%s Service restart failed: %s", v.name, err.Error())
				continue
			}
			m.wg.Add(1)
		}
	}
}

func Register(name string, h def.Handler, args ...any) error {
	_, ok := m.serviceByName[name]
	if ok {
		return RegisterExistError.Fill("name")
	}

	c := context{name, args, nil, &m.wg, true, h}
	err := c.start()
	if err != nil {
		return err
	}
	m.serviceByName[name] = &c
	m.wg.Add(1)
	return nil
}

func Destroy(name string) {
	s, ok := m.serviceByName[name]
	if ok {
		s.stop()
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
