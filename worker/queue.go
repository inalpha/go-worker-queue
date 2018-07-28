package worker

func Queue(nworkers int, work chan Work) {
	queue := make(chan chan Work, nworkers)
	for i := 0; i < nworkers; i++ {
		w := New(i+1, queue)
		w.Start()
	}

	go func() {
		for {
			select {
			case w := <-work:
				go func() {
					worker := <-queue
					worker <- w
				}()
			}
		}
	}()
}