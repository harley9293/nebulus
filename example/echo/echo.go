package echo

import (
	"github.com/harley9293/nebulus"
	"github.com/harley9293/nebulus/pkg/def"
)

type Service struct {
	def.DefaultHandler
}

func (m *Service) Print(req string) string {
	return "echo: " + req
}

func Example() {
	err := nebulus.Register("echo", new(Service))
	if err != nil {
		panic(err)
	}
	nebulus.Run()
}
