package util

import (
	"sync"
)

type AutoInc struct {
	start, step int
	queue       chan int
	kill        chan bool
	wg          *sync.WaitGroup
	running     bool
}

func NewAutoInc(start, step int) (ai *AutoInc) {
	ai = &AutoInc{
		start:   start,
		step:    step,
		queue:   make(chan int, 4),
		kill:    make(chan bool, 1),
		wg:      new(sync.WaitGroup),
		running: true,
	}
	ai.wg.Add(1)
	go ai.process()
	return
}

func (ai *AutoInc) process() {
	defer func() { recover() }()
	for i := ai.start; ai.running; i = i + ai.step {
		select {
		case <-ai.kill:
			ai.running = false
		case ai.queue <- i:
		}
	}

	close(ai.queue)
	ai.wg.Done()
}

func (ai *AutoInc) Id() int {
	return <-ai.queue
}

func (ai *AutoInc) Close() {
	ai.kill <- true
	ai.wg.Wait()
}
