package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"otus/internal/handler"
	"otus/internal/repository"
	"otus/internal/service"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
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

	//Server 8080
	r := gin.Default()

	api := r.Group("/api")

	{
		api.POST("/user", handler.CreateUser)
		api.GET("/users", handler.GetUsers)
		api.GET("/user/:id", handler.GetUser)
		api.PUT("/user/:id", handler.UpdateUser)
		api.DELETE("/user/:id", handler.DeleteUser)

		api.POST("/task", handler.CreateTask)
		api.GET("/tasks", handler.GetTasks)
		api.GET("/task/:id", handler.GetTask)
		api.PUT("/task/:id", handler.UpdateTask)
		api.DELETE("/task/:id", handler.DeleteTask)
	}

	srv := &http.Server{Addr: ":8080", Handler: r}

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("http server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("http server error:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	fmt.Println("Получен сигнал: ", sig)

	shutdawnCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdawnCtx)

	cancel()
	wg.Wait()

	fmt.Println("Горутины завершины")
}
