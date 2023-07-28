package nebulus

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/harley9293/blotlog"
	"github.com/harley9293/nebulus/internal/service"
	"github.com/harley9293/nebulus/pkg/def"
)

type server struct {
	running bool
	kill    chan os.Signal
}

var svr *server

func init() {
	log.Info("init nebulus...")
	svr = new(server)
	svr.kill = make(chan os.Signal, 1)

	// Handle abnormal signals
	signal.Notify(svr.kill, syscall.SIGINT, syscall.SIGTERM)
}

func Register(name string, h def.Handler, args ...any) error {
	return service.Register(name, h, args...)
}

func Destroy(name string) {
	service.Destroy(name)
}

func Send(f string, in ...any) {
	service.Send(f, in...)
}

func Call(f string, inout ...any) error {
	return service.Call(f, inout...)
}

func Run() {
	svr.running = true
	log.Info("nebulus running")
	for svr.running {
		select {
		case <-time.After(16 * time.Millisecond):
			service.Tick()
		case <-svr.kill:
			svr.running = false
			log.Info("recv kill signal, nebulus shutdown...")
		}
	}

	svr.stop()
}

func Shutdown() {
	svr.kill <- syscall.SIGTERM
	time.Sleep(1 * time.Second)
}

func (s *server) stop() {
	service.Stop()
	log.Info("nebulus stop")
}
