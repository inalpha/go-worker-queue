package main

import (
	"net/http"
	"log"
	"time"

	"worker-queue/worker"
)

type WorkRequest struct {
	Name  string
	Delay time.Duration
}

func (w *WorkRequest) Do() {
	log.Printf("Wating %fs for %s\n", w.Delay.Seconds(), w.Name)
	time.Sleep(w.Delay)
}

var WorkQueue = make(chan worker.Work, 100)

func main() {
	worker.Queue(2000, WorkQueue)
	http.HandleFunc("/work", func (w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Set("Allow", "POST")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		name := r.FormValue("name")
		if name == "" {
			http.Error(w, "You must specify a name.", http.StatusBadRequest)
			return
		}

		delay, err := time.ParseDuration(r.FormValue("delay"))
		if err != nil {
			http.Error(w, "Bad delay value: "+err.Error(), http.StatusBadRequest)
			return
		}

		if delay.Seconds() < 1 || delay.Seconds() > 10 {
			http.Error(w, "The delay must be between 1 and 10 seconds, inclusively.", http.StatusBadRequest)
    		return
		}

		work := WorkRequest{Name: name, Delay: delay}
  
  		WorkQueue <- &work

		w.WriteHeader(http.StatusCreated)
	})
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err.Error())
	}
}
