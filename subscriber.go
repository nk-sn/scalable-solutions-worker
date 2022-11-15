package main

import (
	"context"
	"fmt"
)

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
			}
		}
	}()
}
