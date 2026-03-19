package repository

import (
	"context"
	"otus/internal/db"
	"otus/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getUserCollection() *mongo.Collection {
	return db.MongoDB.Collection("user")
}

func getTaskCollection() *mongo.Collection {
	return db.MongoDB.Collection("task")
}

//---------------User------------------

func MongoAddUser(u model.User) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u.UserID = nextUserID()
	_, err := getUserCollection().InsertOne(ctx, u)
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}

func MongoGetUserByID(id int) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var u model.User
	err := getUserCollection().FindOne(ctx, bson.M{"userid": id}).Decode(&u)
	if err != nil {
		return model.User{}, ErrNotFound
	}
	return u, nil
}

func MongoGetAllUsers() ([]model.User, error) {
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

func MongoUpdateUser(id int, updated model.User) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updated.UserID = id
	result, err := getUserCollection().UpdateOne(ctx, bson.M{"userid": id}, bson.M{"$set": updated})
	if err != nil {
		return model.User{}, err
	}
	if result.MatchedCount == 0 {
		return model.User{}, ErrNotFound
	}
	return updated, nil
}

func MongoDeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := getUserCollection().DeleteOne(ctx, bson.M{"userid": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}

//----------------Task-------------------------

func MongoAddTask(t model.Task) (model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.TaskID = nextTaskID()
	_, err := getTaskCollection().InsertOne(ctx, t)
	if err != nil {
		return model.Task{}, err
	}
	return t, nil
}

func MongoGetTaskByID(id int) (model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var t model.Task
	err := getTaskCollection().FindOne(ctx, bson.M{"taskid": id}).Decode(&t)
	if err != nil {
		return model.Task{}, ErrNotFound
	}
	return t, nil
}

func MongoGetAllTasks() ([]model.Task, error) {
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

func MongoUpdateTask(id int, updated model.Task) (model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updated.TaskID = id
	result, err := getTaskCollection().UpdateOne(ctx, bson.M{"taskid": id}, bson.M{"$set": updated})
	if err != nil {
		return model.Task{}, err
	}
	if result.MatchedCount == 0 {
		return model.Task{}, ErrNotFound
	}
	return updated, nil
}

func MongoDeleteTask(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := getTaskCollection().DeleteOne(ctx, bson.M{"taskid": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}
