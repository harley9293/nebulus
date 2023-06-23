package echo

import "github.com/harley9293/nebulus/pkg/def"

type Service struct {
	def.DefaultHandler
}

func (m *Service) Print(req string) string {
	return "echo: " + req
}
