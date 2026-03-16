package grpcserver

import (
	"context"
	"otus/internal/model"
	"otus/internal/repository"
	"otus/pkg/pb"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
}

type TaskServer struct {
	pb.UnimplementedTaskServiceServer
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	u, err := repository.AddUser(model.User{Username: req.Username})
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{Id: int32(u.UserID), Username: u.Username}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	u, err := repository.GetUserByID(int(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{Id: int32(u.UserID), Username: u.Username}, nil
}

func (s *UserServer) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	all := repository.GetAllUsers()
	resp := make([]*pb.UserResponse, len(all))
	for i, u := range all {
		resp[i] = &pb.UserResponse{Id: int32(u.UserID), Username: u.Username}
	}
	return &pb.GetUsersResponse{Users: resp}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	u, err := repository.UpdateUser(int(req.Id), model.User{Username: req.Username})
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{Id: int32(u.UserID), Username: u.Username}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteResponse, error) {
	if err := repository.DeleteUser(int(req.Id)); err != nil {
		return nil, err
	}
	return &pb.DeleteResponse{Message: "user deleted"}, nil
}

func (s *TaskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.TaskResponse, error) {
	t, err := repository.AddTask(model.Task{Title: req.Title})
	if err != nil {
		return nil, err
	}
	return &pb.TaskResponse{Id: int32(t.TaskID), Title: t.Title}, nil
}
func (s *TaskServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.TaskResponse, error) {
	t, err := repository.GetTaskByID(int(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.TaskResponse{Id: int32(t.TaskID), Title: t.Title}, nil
}

func (s *TaskServer) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	all := repository.GetAllTasks()
	resp := make([]*pb.TaskResponse, len(all))
	for i, u := range all {
		resp[i] = &pb.TaskResponse{Id: int32(u.TaskID), Title: u.Title}
	}
	return &pb.GetTasksResponse{Tasks: resp}, nil
}

func (s *TaskServer) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.TaskResponse, error) {
	t, err := repository.UpdateTask(int(req.Id), model.Task{Title: req.Title})
	if err != nil {
		return nil, err
	}
	return &pb.TaskResponse{Id: int32(t.TaskID), Title: t.Title}, nil
}

func (s *TaskServer) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteResponse, error) {
	if err := repository.DeleteTask(int(req.Id)); err != nil {
		return nil, err
	}
	return &pb.DeleteResponse{Message: "task deleted"}, nil
}
