package def

import "testing"

func TestDefaultHandler(t *testing.T) {
	h := &DefaultHandler{}
	err := h.OnInit()
	if err != nil {
		t.Fatal(err)
	}

	h.OnTick()
	h.OnStop()
}
