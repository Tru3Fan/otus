package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	MongoDB     *mongo.Database
	RedisClient *redis.Client
	PostgresDB  *sql.DB
)

func Connect() error {
	if err := connectRedis(); err != nil {
		return err
	}
	if err := contextPostgres(); err != nil {
		return err
	}
	return nil
}

func connectRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return err
	}
	return nil
}

func contextPostgres() error {
	dsn := os.Getenv("POSTGRES_DSN")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}

	PostgresDB = db
	fmt.Println("connected to PostgresSQL")

	if err := runMigrations(db); err != nil {
		return err
	}
	return nil
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	fmt.Println("migrations applied successfully")
	return nil

}

func Disconnect() {
	if MongoDB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		MongoDB.Client().Disconnect(ctx)
	}
	if RedisClient != nil {
		RedisClient.Close()
	}
	if PostgresDB != nil {
		PostgresDB.Close()
	}

}
