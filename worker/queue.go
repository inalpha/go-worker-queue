package worker

type Queue struct {
	queue chan Work
}
func NewQueue(nworkers int) Queue {
	workQueue := make(chan Work, 100)
	pool := make(chan chan Work, nworkers)
	for i := 0; i < nworkers; i++ {
		w := New(i+1, pool)
		w.Start()
	}

	go func() {
		for {
			select {
			case w := <-workQueue:
				go func() {
					worker := <-pool
					worker <- w
				}()
			}
		}
	}()

	return Queue{
		queue: workQueue,
	}
}

func (q *Queue) Submit(w Work) {
	q.queue <- w
}


