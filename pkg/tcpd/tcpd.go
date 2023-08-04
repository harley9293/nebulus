package tcpd

import (
	log "github.com/harley9293/blotlog"
	"github.com/harley9293/nebulus/pkg/def"
	"github.com/harley9293/nebulus/pkg/errors"
	"github.com/harley9293/nebulus/pkg/util"
	"net"
	"reflect"
)

var InitArgsSizeError = errors.New("http init args size error, got:%d")
var InitArgsTypeError = errors.New("http init args type error, got:%T")

type Service struct {
	def.DefaultHandler

	srv net.Listener

	ai      *util.AutoInc
	connMap map[int]conn
	connCh  chan net.Conn
}

func NewService() *Service {
	return &Service{connMap: make(map[int]conn), connCh: make(chan net.Conn), ai: util.NewAutoInc(1, 1)}
}

func (m *Service) OnInit(args ...any) error {
	if len(args) != 1 {
		return InitArgsSizeError.Fill(len(args))
	}

	address, ok := args[0].(string)
	if !ok {
		return InitArgsTypeError.Fill(reflect.TypeOf(args[0]))
	}

	var err error = nil
	m.srv, err = net.Listen("tcp", address)
	if err != nil {
		return err
	}
	go m.OnAccept()

	return nil
}

// OnTick is called by service goroutine.
func (m *Service) OnTick() {
	select {
	case c := <-m.connCh:
		m.connMap[m.ai.Id()] = conn{c}
	default:
	}
}

// OnAccept is called by stand-alone goroutine.
func (m *Service) OnAccept() {
	for {
		c, err := m.srv.Accept()
		if err != nil {
			log.Error("Error accepting, err: %s", err)
			continue
		}
		m.connCh <- c
	}
}
