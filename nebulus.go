package nebulus

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/harley9293/blotlog"
	"github.com/harley9293/nebulus/service"
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

func Register(name string, h service.Handler, args ...any) error {
	return service.Register(name, h, args...)
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

func (s *server) stop() {
	service.Stop()
	log.Info("nebulus stop")
}
