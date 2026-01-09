package main

import (
	"otus/internal/service"
	"time"
)

func main() {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		service.GenerateAndStore()
		continue
	}
}
