package main

import (
	"context"
	"fmt"
	"time"
)

var ResultChannelListeningTimeout = 4 * time.Second

func ListenResultChannel(ctx context.Context, resultChannel chan JobResult) {
	go func() {
		fmt.Println("Started listening result channel")
		for {
			select {
			case <-ctx.Done():
				for res := range resultChannel {
					fmt.Println(fmt.Sprintf("result: %s", res.JobID))
				}
				fmt.Println("Finished listening result channel")
				return
			case res := <-resultChannel:
				fmt.Println(fmt.Sprintf("result: %s", res.JobID))
			default:
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()
}
