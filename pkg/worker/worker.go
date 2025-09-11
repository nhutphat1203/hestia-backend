package worker

import "log"

type Worker struct {
	ID      int
	JobChan chan Job
	Quit    chan bool
}

func NewWorker(id int) *Worker {
	return &Worker{
		ID:      id,
		JobChan: make(chan Job),
		Quit:    make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case job := <-w.JobChan:
				err := job.Execute()
				if err != nil {
					log.Printf("Worker %d failed: %v", w.ID, err)
				}
			case <-w.Quit:
				log.Printf("Worker %d stopping", w.ID)
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.Quit <- true
}
