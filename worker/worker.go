package worker

import (
	"log"
)

type Work interface {
	Do()
}

type Worker struct {
	ID          int
	Work        chan Work
	WorkerQueue chan chan Work
	QuitChan    chan bool
}

func New(id int, queue chan chan Work) Worker {
	return Worker{
		ID:          id,
		Work:        make(chan Work),
		WorkerQueue: queue,
		QuitChan:    make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.Work
			select {
			case work := <-w.Work:
				log.Printf("worker%d recived work\n", w.ID)
				work.Do()
				log.Printf("worker%d finished work\n", w.ID)
			case <-w.QuitChan:
				log.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
