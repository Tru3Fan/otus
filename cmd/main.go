package main

import (
	"context"
	"os"
	"os/signal"
	"otus/internal/repository"
	"otus/internal/service"
	"syscall"
	"time"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ch := make(chan repository.Storable, 50)

	go service.GenerateAndCreate(ctx, ch)
	go repository.Add(ctx, ch)
	go repository.LogNew(ctx)

	time.Sleep(1 * time.Second)

	<-ctx.Done()
}
