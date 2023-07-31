package util

import "sync/atomic"

type AutoInc struct {
	start, step int
	queue       chan int
	running     atomic.Value
}

func NewAutoInc(start, step int) (ai *AutoInc) {
	ai = &AutoInc{
		start: start,
		step:  step,
		queue: make(chan int, 4),
	}
	ai.running.Store(true)
	go ai.process()
	return
}

func (ai *AutoInc) process() {
	defer func() { recover() }()
	for i := ai.start; ai.running.Load().(bool); i = i + ai.step {
		ai.queue <- i
	}
}

func (ai *AutoInc) Id() int {
	return <-ai.queue
}

func (ai *AutoInc) Close() {
	ai.running.Store(false)
	close(ai.queue)
}
