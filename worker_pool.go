package main

import (
	"context"
	"fmt"
	"time"
)

var ChannelsCloseTimeout = 5 * time.Second

type Worker interface {
	Run()
	Shutdown()
}

type WorkerPool struct {
	cxt            context.Context
	cancel         context.CancelFunc
	workersRunning []Worker
	workersQueue   []Worker
	jobs           []Job
	jobsChannel    chan Job
	resultChannel  chan JobResult
	errorChannel   chan error
}

func NewWorkerPool(cxt context.Context) *WorkerPool {
	workerPoolContext, cancel := context.WithCancel(cxt)

	return &WorkerPool{
		cxt:           workerPoolContext,
		cancel:        cancel,
		jobsChannel:   make(chan Job),
		resultChannel: make(chan JobResult),
		errorChannel:  make(chan error),
	}
}

func (wp *WorkerPool) Start() {
	fmt.Println("Worker Pool turned on")
	go func() {
		for {
			select {
			case <-wp.cxt.Done():
				// waiting to graceful close channels
				time.Sleep(ChannelsCloseTimeout)

				close(wp.jobsChannel)
				fmt.Println("Worker Pool jobs channel is closed")
				close(wp.resultChannel)
				fmt.Println("Worker Pool result channel is closed")
				close(wp.errorChannel)
				fmt.Println("Worker Pool error channel is closed")

				fmt.Println("Worker Pool turned off")
				return
			default:
				if len(wp.workersQueue) != 0 {
					for _, w := range wp.workersQueue {
						w.Run()
						wp.workersRunning = append(wp.workersRunning, w)
					}
					wp.workersQueue = make([]Worker, 0)
				}
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()
}

func (wp *WorkerPool) Shutdown() {
	wp.cancel()
}

func (wp *WorkerPool) AddWorkers(count uint8) {
	for i := uint8(0); i < count; i++ {
		worker := NewWork(wp.cxt, wp.jobsChannel, wp.errorChannel, wp.resultChannel)
		wp.workersQueue = append(wp.workersQueue, worker)
	}
	fmt.Println(fmt.Sprintf("Added %d workers to worker pool", count))
}

func (wp *WorkerPool) RemoveWorkers(count uint8) {
	removed := 0
	for count != 0 && len(wp.workersRunning) != 0 {
		worker := wp.workersRunning[len(wp.workersRunning)-1]
		worker.Shutdown()
		wp.workersRunning = wp.workersRunning[:len(wp.workersRunning)-1]
		count--
		removed++
	}
	fmt.Println(fmt.Sprintf("Removed %d workers from worker pool", removed))
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.jobs = append(wp.jobs, job)
	go func() {
		for {
			select {
			case <-wp.cxt.Done():
				return
			default:
				wp.jobsChannel <- job
			}
		}
	}()
	fmt.Println(fmt.Sprintf("Job %s added to worker pool", job.ID()))
}

func (wp *WorkerPool) SubscribeResults() chan JobResult {
	return wp.resultChannel
}

func (wp *WorkerPool) SubscribeErrors() chan error {
	return wp.errorChannel
}
