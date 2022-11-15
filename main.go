package main

/*
Необходимо реализовать пулл воркеров для обработки задач.
Добавление задачи должно быть неблокирующим.
Пулл должен иметь возможностью динамически изменять количество воркеров.
Результат обработки должен возвращаться асинхронно.
В качестве дополнительного задания добавить возможность делать graceful shutdown.
*/

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"
)

type Job interface {
	ID() string
	Do() error
}

type WorkerPooler interface {
	Start()
	Shutdown()
	AddWorkers(count uint8)
	RemoveWorkers(count uint8)
	AddJob(job Job)
	SubscribeResults() chan JobResult
	SubscribeErrors() chan error
}

var MainShutdownTimeout = 10 * time.Second

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	job1 := NewSleepJob()
	job2 := NewSleepJob()

	workerPool := NewWorkerPool(ctx)
	workerPool.AddJob(job1)
	workerPool.AddJob(job2)
	workerPool.AddWorkers(5)
	workerPool.Start()

	ListenResultChannel(ctx, workerPool.SubscribeResults())
	ListenErrorChannel(ctx, workerPool.SubscribeErrors())

	time.Sleep(10 * time.Second)

	workerPool.RemoveWorkers(2)

	<-ctx.Done()

	fmt.Println("Got signal to turn worker pool off")

	// Waiting to let graceful shutdown
	time.Sleep(MainShutdownTimeout)

	fmt.Println("The application exited")
}
