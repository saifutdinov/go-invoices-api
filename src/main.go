package main

import (
	"context"
	"log"
	"time"

	"github.com/saifutdinov/go-invoices-api/pkg/chronos"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler := chronos.New(ctx)

	scheduler.Schedule(
		chronos.NewTask(
			time.Now().Add(3*time.Second),
			func(ctx context.Context) {
				log.Println(
					"TASK EXECUTED:",
					time.Now().Format(time.RFC3339),
				)
			},
		),
	)

	log.Println("task scheduled at:", time.Now().Format(time.RFC3339))

	time.Sleep(5 * time.Second)

	scheduler.Stop()
}
