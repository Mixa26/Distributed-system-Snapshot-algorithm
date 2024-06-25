package helper

import "Snapshot/servent/handler"

// Implemention provided by https://medium.com/@ahmet9417/golang-thread-pool-and-scheduler-434dd094715a

type Job struct {
	Handler handler.MessageHandler
}

type Pool struct {
	JobQueue  chan Job      // Prodcution Line
	readyPool chan chan Job // Worker Channel
}

// Manager in action
func (q *Pool) dispatch() {
	for {
		select {
		case job := <-q.JobQueue:
			workerXChannel := <-q.readyPool //free worker x founded
			workerXChannel <- job           // here is your job worker x
		}
	}
}

type worker struct {
	readyPool chan chan Job //get work from the boss
	job       chan Job
}

func (w *worker) Start() {
	go func() {
		for {
			w.readyPool <- w.job //hey i am ready to work on new job
			select {
			case job := <-w.job: // hey i am waiting for new job
				job.Handler.Run() // ok i am on it
			}
		}
	}()
}

func InitiatePool() *Pool {
	pool := &Pool{
		JobQueue:  make(chan Job),
		readyPool: make(chan chan Job),
	}

	go pool.dispatch()

	// We assign 5 workers.
	for i := 0; i < 5; i++ {
		worker := &worker{
			readyPool: pool.readyPool,
			job:       make(chan Job),
		}
		worker.Start()
	}

	return pool
}
