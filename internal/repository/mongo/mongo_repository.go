package mongo

import (
	"context"
	"otus/internal/db"
	"otus/internal/model"
	"otus/internal/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct{}
type TaskRepo struct{}

func NewUserRepo() repository.UserRepository {
	return &UserRepo{}
}

func NewTaskRepo() repository.TaskRepository {
	return &TaskRepo{}
}

func getUserCollection() *mongo.Collection {
	return db.MongoDB.Collection("user")
}

func getTaskCollection() *mongo.Collection {
	return db.MongoDB.Collection("task")
}

//---------------User------------------

func (r *UserRepo) AddUser(u model.User) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u.UserID = nextUserID()
	_, err := getUserCollection().InsertOne(ctx, u)
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}

func (r *UserRepo) GetUserByID(id int) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var u model.User
	err := getUserCollection().FindOne(ctx, bson.M{"userid": id}).Decode(&u)
	if err != nil {
		return model.User{}, repository.ErrNotFound
	}
	return u, nil
}

func (r *UserRepo) GetAllUsers() ([]model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := getUserCollection().Find(ctx, bson.M{})
	if err != nil {
		return []model.User{}, err
	}
	defer cursor.Close(ctx)

	var users []model.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) UpdateUser(id int, updated model.User) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updated.UserID = id
	result, err := getUserCollection().UpdateOne(ctx, bson.M{"userid": id}, bson.M{"$set": updated})
	if err != nil {
		return model.User{}, err
	}
	if result.MatchedCount == 0 {
		return model.User{}, repository.ErrNotFound
	}
	return updated, nil
}

func (r *UserRepo) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := getUserCollection().DeleteOne(ctx, bson.M{"userid": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *UserRepo) GetUserByTelegramID(telegramID int64) (model.User, error) {
	return model.User{}, nil
}

//----------------Task-------------------------

func (r *TaskRepo) AddTask(t model.Task) (model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.TaskID = nextTaskID()
	_, err := getTaskCollection().InsertOne(ctx, t)
	if err != nil {
		return model.Task{}, err
	}
	return t, nil
}

func (r *TaskRepo) GetTaskByID(id int) (model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var t model.Task
	err := getTaskCollection().FindOne(ctx, bson.M{"taskid": id}).Decode(&t)
	if err != nil {
		return model.Task{}, repository.ErrNotFound
	}
	return t, nil
}

func (r *TaskRepo) GetAllTasks() ([]model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := getTaskCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []model.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepo) UpdateTask(id int, updated model.Task) (model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updated.TaskID = id
	result, err := getTaskCollection().UpdateOne(ctx, bson.M{"taskid": id}, bson.M{"$set": updated})
	if err != nil {
		return model.Task{}, err
	}
	if result.MatchedCount == 0 {
		return model.Task{}, repository.ErrNotFound
	}
	return updated, nil
}

func (r *TaskRepo) DeleteTask(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := getTaskCollection().DeleteOne(ctx, bson.M{"taskid": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *TaskRepo) GetTasksByUserID(userID int) ([]model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := getTaskCollection().Find(ctx, bson.M{"userid": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []model.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepo) GetTasksByStatus(status string) ([]model.Task, error) {
	return nil, nil
}

func nextUserID() int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := getUserCollection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return 1
	}
	return int(count) + 1
}

func nextTaskID() int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := getTaskCollection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return 1
	}
	return int(count) + 1
}
