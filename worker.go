package main

import (
	"github.com/google/uuid"

	"context"
	"fmt"
)

type Work struct {
	id            string
	cxt           context.Context
	cancel        context.CancelFunc
	jobsChannel   chan Job
	errorChannel  chan error
	resultChannel chan JobResult
}

func NewWork(cxt context.Context, jobsChannel chan Job, errorChannel chan error, resultChannel chan JobResult) *Work {
	workerContext, cancel := context.WithCancel(cxt)
	return &Work{
		id:            uuid.NewString(),
		cxt:           workerContext,
		cancel:        cancel,
		jobsChannel:   jobsChannel,
		errorChannel:  errorChannel,
		resultChannel: resultChannel,
	}
}

func (w *Work) Run() {
	fmt.Println(fmt.Sprintf("The worker %s turned on", w.id))
	go func() {
		for {
			select {
			case <-w.cxt.Done():
				fmt.Println(fmt.Sprintf("The worker %s turned off", w.id))
				return
			case job := <-w.jobsChannel:
				err := job.Do()
				if err != nil {
					w.errorChannel <- err
				}
				w.resultChannel <- JobResult{JobID: job.ID(), Err: err}
			}
		}
	}()
}

func (w *Work) Shutdown() {
	w.cancel()
}
