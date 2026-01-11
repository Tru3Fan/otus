package main

import (
	"otus/internal/repository"
	"otus/internal/service"
	"time"
)

func main() {

	ch := make(chan repository.Storable, 50)
	logCh := make(chan struct{})

	go repository.Add(ch)
	go repository.LogNew(logCh)
	go service.GenerateAndCreate(ch)

	time.Sleep(1 * time.Second)

	close(logCh)
}
