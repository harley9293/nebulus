package nebulus

import (
	log "github.com/harley9293/blotlog"
	"github.com/harley9293/nebulus/internal/service"
	"github.com/harley9293/nebulus/pkg/def"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	kill chan os.Signal
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

func Shutdown() {
	svr.kill <- syscall.SIGTERM
}

func Run() {
	log.Info("nebulus running")
	<-svr.kill
	log.Info("recv kill signal, nebulus shutdown...")
	service.Stop()
	log.Info("nebulus stop")
}
