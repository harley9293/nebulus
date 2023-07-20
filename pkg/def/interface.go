package def

type Handler interface {
	OnInit(args ...any) error
	OnStop()
	OnTick()
}

type DefaultHandler struct {
}

func (h *DefaultHandler) OnInit(args ...any) error {
	return nil
}

func (h *DefaultHandler) OnStop() {
}

func (h *DefaultHandler) OnTick() {
}
