package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"otus/internal/bot"
	"otus/internal/db"
	"otus/internal/handler"
	postgresRepo "otus/internal/repository/postgres"
	"otus/internal/service"
	"strconv"
	"syscall"
	"time"

	_ "otus/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Otus API
// @version         1.0
// @description     API для управления пользователями и задачами
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {

	godotenv.Load()

	if err := db.Connect(); err != nil {
		fmt.Println("Error connecting to database", err)
		os.Exit(1)
	}
	fmt.Println("all database connection established")
	defer db.Disconnect()

	userRepo := postgresRepo.NewUserRepo()
	taskRepo := postgresRepo.NewTaskRepo()

	userSvc := service.NewUserService(userRepo)
	taskSvc := service.NewTaskService(taskRepo)

	userHandler := handler.NewUserHandler(userSvc)
	taskHandler := handler.NewTaskHandler(taskSvc)

	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramToken != "" {
		adminID, _ := strconv.ParseInt(os.Getenv("ADMIN_TELEGRAM_ID"), 10, 64)
		tgBot, err := bot.NewBot(telegramToken, taskSvc, userSvc, adminID)
		if err != nil {
			fmt.Println("Error creating telegram bot", err)
		} else {
			go func() {
				fmt.Println("Telegram bot started")
				if err := tgBot.Start(); err != nil {
					fmt.Println("telegram bot error:", err)
				}
			}()
		}
	}
	//Server 8080
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.POST("/login", handler.Login)
		api.GET("/users", userHandler.GetUsers)
		api.GET("/user/:id", userHandler.GetUser)
		api.GET("/user/:id/tasks", taskHandler.GetTasksByUser)
		api.GET("/tasks", taskHandler.GetTasks)
		api.GET("/task/:id", taskHandler.GetTask)
		api.GET("/tasks/status", taskHandler.GetTasksByStatus)

		protected := api.Group("/")
		protected.Use(handler.AuthMiddleware())
		{
			protected.POST("/user", userHandler.CreateUser)
			protected.PUT("/user/:id", userHandler.UpdateUser)
			protected.DELETE("/user/:id", userHandler.DeleteUser)

			protected.POST("/task", taskHandler.CreateTask)
			protected.PUT("/task/:id", taskHandler.UpdateTask)
			protected.DELETE("/task/:id", taskHandler.DeleteTask)
			protected.PUT("/task/:id/status", taskHandler.UpdateTaskStatus)
		}
	}

	srv := &http.Server{Addr: ":8080", Handler: r}

	go func() {
		fmt.Println("http server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("http server error:", err)
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	fmt.Println("Получен сигнал: ", sig)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
}
