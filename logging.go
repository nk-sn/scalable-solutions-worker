package main

import (
	"context"
	"fmt"
)

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
			}
		}
	}()
}
