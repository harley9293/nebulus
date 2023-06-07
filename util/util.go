package util

import (
	"bytes"
	"runtime"

	log "github.com/harley9293/blotlog"
)

func TryE() {
	errs := recover()
	if errs == nil {
		return
	}

	log.Error("%v\n", errs)
	log.Error(panicTrace(8))
}

func panicTrace(kb int) string {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10) //4KB
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return string(stack)
}
