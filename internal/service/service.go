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
	serviceByName sync.Map
	wg            sync.WaitGroup
}

func init() {
	m = new(mgr)
}

//=============== tick goroutine =================

func Stop() {
	m.serviceByName.Range(func(key, value interface{}) bool {
		v := value.(*context)
		v.stop()
		return true
	})
	m.wg.Wait()
	m.serviceByName = sync.Map{}
}

func Tick() {
	m.serviceByName.Range(func(key, value interface{}) bool {
		v := value.(*context)
		if !v.status() {
			log.Warn("%s Service is attempting to restart", v.name)
			err := v.start()
			if err != nil {
				log.Error("%s Service restart failed: %s", v.name, err.Error())
			} else {
				m.wg.Add(1)
			}
		}
		return true
	})
}

//=============== other goroutine =================

func Register(name string, h def.Handler, args ...any) error {
	c := context{name: name, args: args, wg: &m.wg, Handler: h}
	err := c.start()
	if err != nil {
		return err
	}

	_, loaded := m.serviceByName.LoadOrStore(name, &c)
	if loaded {
		c.stop() // Stop the context if a service by the same name already exists
		return RegisterExistError.Fill(name)
	}

	m.wg.Add(1)
	return nil
}

func Destroy(name string) {
	value, ok := m.serviceByName.Load(name) // Load returns the value stored in the map for a key.
	if ok {
		s := value.(*context)
		s.stop()
		m.serviceByName.Delete(name) // Delete removes the value for a key from the map.
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

	value, ok := m.serviceByName.Load(name)
	if !ok {
		log.Warn("%s service is not registered, send failed", name)
		return
	}

	c := value.(*context)
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

	value, ok := m.serviceByName.Load(name)
	if !ok {
		err := NotRegisterError.Fill(name)
		log.Warn(err.Error())
		return err
	}

	c := value.(*context)
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
