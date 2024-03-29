package main

import (
	"context"
	"downloader_email/configs"
	"downloader_email/internal"
	"downloader_email/rabbitmq"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	configs.LoadEnvVariables()

	ctx, cancel := context.WithCancel(context.Background())
	rabbit := rabbitmq.Start(ctx)
	defer cancel()

	_ = internal.NewEmailService(rabbit)
	wg.Wait() // dont exit program
}
