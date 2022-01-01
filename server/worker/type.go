package worker

import (
	"sync"
)

type Job interface {
	e()
}

type Worker struct {
	*sync.RWMutex
	work []chan *Job
	i int
	b int
}

func NewWorkers(i int) (w *Worker) {
	w = &Worker{
		work: make([]chan *Job, i),
		i: i,
		b: -1,
	}
	for l := 0; l < i; l = l + 1 {
		w.work[l] = make(chan *Job)
		go w.Grap(l)
	}
	return
}

func (w *Worker) Grap(i int) {
	for {
		select {
			case msg1 := <-w.work[i]:
				l := *msg1
				l.e()
		}
	}
}

func (w *Worker) AddTask (e *Job) {
	w.b = (w.b % w.i) + 1
	go func() {
		w.work[w.b] <- e
	} ()
}
