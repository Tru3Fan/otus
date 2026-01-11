package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"otus/internal/repository"
	"otus/internal/service"
	"syscall"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan repository.Storable, 50)

	go service.GenerateAndCreate(ctx, ch)
	go repository.Add(ctx, ch)
	go repository.LogNew(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan

	cancel()

	fmt.Println("Получен сигнал: ", sig)
}
