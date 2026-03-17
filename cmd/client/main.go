package main

import (
	"context"
	"fmt"
	"log"
	"otus/pkg/pb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect:", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	userClient := pb.NewUserServiceClient(conn)
	taskClient := pb.NewTaskServiceClient(conn)

	// CreateUser
	createdUser, err := userClient.CreateUser(ctx, &pb.CreateUserRequest{Username: "gRPC_Ivan"})
	if err != nil {
		log.Fatal("CreateUser error:", err)
	}
	fmt.Printf("Created user: id=%d username=%s\n", createdUser.Id, createdUser.Username)

	// GetUser
	user, err := userClient.GetUser(ctx, &pb.GetUserRequest{Id: createdUser.Id})
	if err != nil {
		log.Fatal("GetUser error:", err)
	}
	fmt.Printf("User: id=%d username=%s\n", user.Id, user.Username)

	// GetUsers
	users, err := userClient.GetUsers(ctx, &pb.GetUsersRequest{})
	if err != nil {
		log.Fatal("GetUsers error:", err)
	}
	fmt.Printf("All users (%d):\n", len(users.Users))
	for _, u := range users.Users {
		fmt.Printf("id=%d username=%s\n", u.Id, u.Username)
	}

	// UpdateUser
	updatedUser, err := userClient.UpdateUser(ctx, &pb.UpdateUserRequest{Id: createdUser.Id, Username: "gRPC_Ivan_Updated"})
	if err != nil {
		log.Fatal("UpdateUser error:", err)

	}
	fmt.Printf("Updated user: id=%d username=%s\n", updatedUser.Id, updatedUser.Username)

	// DeleteUser
	deleteResp, err := userClient.DeleteUser(ctx, &pb.DeleteUserRequest{Id: createdUser.Id})
	if err != nil {
		log.Fatal("DeleteUser error:", err)
	}
	fmt.Printf("Delete user: %s\n", deleteResp.Message)

	// CreateTask
	createdTask, err := taskClient.CreateTask(ctx, &pb.CreateTaskRequest{Title: "gRPC_Task"})
	if err != nil {
		log.Fatal("CreateTask error:", err)
	}
	fmt.Printf("Created task: id=%d title=%s\n", createdTask.Id, createdTask.Title)

	// GetTask
	task, err := taskClient.GetTask(ctx, &pb.GetTaskRequest{Id: createdTask.Id})
	if err != nil {
		log.Fatal("GetTask error:", err)
	}
	fmt.Printf("Task: id=%d title=%s\n", task.Id, task.Title)

	// GetTasks
	tasks, err := taskClient.GetTasks(ctx, &pb.GetTasksRequest{})
	if err != nil {
		log.Fatal("GetTasks error:", err)
	}
	fmt.Printf("All tasks (%d):\n", len(tasks.Tasks))
	for _, t := range tasks.Tasks {
		fmt.Printf("id=%d title=%s\n", t.Id, t.Title)
	}

	// UpdateTask
	updatedTask, err := taskClient.UpdateTask(ctx, &pb.UpdateTaskRequest{Id: createdTask.Id, Title: "gRPC_Task_Updated"})
	if err != nil {
		log.Fatal("UpdateTask error:", err)
	}
	fmt.Printf("Updated task: id=%d title=%s\n", updatedTask.Id, updatedTask.Title)

	// DeleteTask
	deleteTaskResp, err := taskClient.DeleteTask(ctx, &pb.DeleteTaskRequest{Id: createdTask.Id})
	if err != nil {
		log.Fatal("DeleteTask error:", err)
	}
	fmt.Printf("Delete task: %s\n", deleteTaskResp.Message)
}
