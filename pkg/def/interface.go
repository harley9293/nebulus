package def

type Handler interface {
	OnInit(args ...any) error
	OnStop()
	OnTick()
	OnPanic()
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

func (h *DefaultHandler) OnPanic() {
}

type Session interface {
	New(key string) Session
	ID() string
	Get(key string) any
	Set(key string, value any)
	UpdateExpire()
	IsExpired() bool
}
