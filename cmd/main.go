package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"otus/internal/db"
	"otus/internal/generat"
	grpcserver "otus/internal/grpc"
	"otus/internal/handler"
	"otus/internal/repository/csv"
	mongoRepo "otus/internal/repository/mongo"
	"otus/internal/service"
	"otus/pkg/pb"
	"sync"
	"syscall"
	"time"

	_ "otus/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
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

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file", err)
		os.Exit(1)
	}

	if err := db.Connect(); err != nil {
		fmt.Println("Error connecting to database", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	if err := csv.LoadAllData(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan csv.Storable, 50)

	var wg sync.WaitGroup

	wg.Add(1)
	go generat.GenerateAndCreate(ctx, ch, &wg)

	wg.Add(1)
	go csv.Add(ctx, ch, &wg)

	wg.Add(1)
	go csv.LogNew(ctx, &wg)

	userRepo := mongoRepo.NewUserRepo()
	taskRepo := mongoRepo.NewTaskRepo()

	userSvc := service.NewUserService(userRepo)
	taskSvc := service.NewTaskService(taskRepo)

	userHandler := handler.NewUserHandler(userSvc)
	taskHandler := handler.NewTaskHandler(taskSvc)

	//Server 8080
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.POST("/login", handler.Login)
		api.GET("/users", userHandler.GetUsers)
		api.GET("/user/:id", userHandler.GetUser)
		api.GET("/tasks", taskHandler.GetTasks)
		api.GET("/task/:id", taskHandler.GetTask)

		protected := api.Group("/")
		protected.Use(handler.AuthMiddleware())
		{
			protected.POST("/user", userHandler.CreateUser)
			protected.PUT("/user/:id", userHandler.UpdateUser)
			protected.DELETE("/user/:id", userHandler.DeleteUser)

			protected.POST("/task", taskHandler.CreateTask)
			protected.PUT("/task/:id", taskHandler.UpdateTask)
			protected.DELETE("/task/:id", taskHandler.DeleteTask)
		}
	}

	srv := &http.Server{Addr: ":8080", Handler: r}

	go func() {
		fmt.Println("http server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("http server error:", err)
		}
	}()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcSrv := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcSrv, &grpcserver.UserServer{})
	pb.RegisterTaskServiceServer(grpcSrv, &grpcserver.TaskServer{})

	go func() {
		fmt.Println("grpc server listening on :50051")
		if err := grpcSrv.Serve(lis); err != nil {
			fmt.Println("grpc server error:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	fmt.Println("Получен сигнал: ", sig)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
	grpcSrv.GracefulStop()

	cancel()
	wg.Wait()

	fmt.Println("Горутины завершины")
}
