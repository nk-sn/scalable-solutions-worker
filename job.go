package main

import (
	"github.com/google/uuid"

	"fmt"
	"time"
)

type JobResult struct {
	JobID string
	Err   error
}

func NewSleepJob() Job {
	return &SleepJob{
		JobID: uuid.NewString(),
	}
}

type SleepJob struct {
	JobID string
}

func (j *SleepJob) ID() string {
	return j.JobID
}

func (j *SleepJob) Do() error {
	fmt.Println(fmt.Sprintf("The job %s is running", j.ID()))
	time.Sleep(2 * time.Second)
	fmt.Println(fmt.Sprintf("The job %s completed", j.ID()))
	return nil
}
