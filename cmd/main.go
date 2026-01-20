package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"otus/internal/repository"
	"otus/internal/service"
	"sync"
	"syscall"
)

func main() {

	if err := repository.LoadAllData(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan repository.Storable, 50)

	var wg sync.WaitGroup

	wg.Add(1)
	go service.GenerateAndCreate(ctx, ch, &wg)

	wg.Add(1)
	go repository.Add(ctx, ch, &wg)

	wg.Add(1)
	go repository.LogNew(ctx, &wg)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan

	fmt.Println("Получен сигнал: ", sig)

	cancel()
	wg.Wait()

	fmt.Println("Горутины завершины")
}
