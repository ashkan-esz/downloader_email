package main

import (
	"context"
	"downloader_email/configs"
	"downloader_email/internal"
	"downloader_email/rabbitmq"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	configs.LoadEnvVariables()

	err := sentry.Init(sentry.ClientOptions{
		Dsn: configs.GetConfigs().SentryDns,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 0.1,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	time.Sleep(time.Duration(configs.GetConfigs().InitialWaitForMailServer) * time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	rabbit := rabbitmq.Start(ctx)
	defer cancel()

	_ = internal.NewEmailService(rabbit)
	fmt.Println("Ready to handle emails")
	wg.Wait() // dont exit program
}
