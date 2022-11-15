package main

import (
	"context"
	"fmt"
	"time"
)

var ErrorChannelListeningTimeout = 4 * time.Second

func log(msg string) {
	// logging
	fmt.Println(fmt.Sprintf("logging: %s", msg))
}

func ListenErrorChannel(ctx context.Context, errorChannel chan error) {
	go func() {
		fmt.Println("Started listening error channel")
		for {
			select {
			case <-ctx.Done():
				for err := range errorChannel {
					log(err.Error())
				}
				fmt.Println("Finished listening error channel")
				return
			case err := <-errorChannel:
				log(err.Error())
			default:
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()
}
