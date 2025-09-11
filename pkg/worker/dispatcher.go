package worker

type Dispatcher struct {
	WorkerPool []*Worker
	JobQueue   chan Job
}

func NewDispatcher(numWorkers int, jobQueue chan Job) *Dispatcher {
	pool := make([]*Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		pool[i] = NewWorker(i + 1)
	}
	return &Dispatcher{
		WorkerPool: pool,
		JobQueue:   jobQueue,
	}
}

func (d *Dispatcher) Run() {
	for _, w := range d.WorkerPool {
		w.Start()
	}
	for job := range d.JobQueue {
		dispatched := false
		for !dispatched {
			for _, w := range d.WorkerPool {
				select {
				case w.JobChan <- job:
					dispatched = true
				default:
				}
			}
		}
	}
}

func (d *Dispatcher) Stop() {
	for _, w := range d.WorkerPool {
		w.Stop()
	}
}
