package echo

import "github.com/harley9293/nebulus/service"

type Service struct {
	service.DefaultHandler
}

func (m *Service) Print(req string) string {
	return "echo: " + req
}
