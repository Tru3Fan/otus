package postgres_test

import (
	"database/sql"
	"os"
	"otus/internal/db"
	"otus/internal/model"
	"otus/internal/repository"
	"otus/internal/repository/postgres"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		t.Skip("POSTGRES_DSN not set, skipping integration tests")
	}
	conn, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	require.NoError(t, conn.Ping())

	db.PostgresDB = conn

	_, err = conn.Exec("DROP TABLE IF EXISTS tasks; DROP TABLE IF EXISTS users;")
	require.NoError(t, err)

	_, err = conn.Exec(`
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    telegram_user_id BIGINT,
    telegram_username TEXT
);
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    assigned_by INT REFERENCES users(id) ON DELETE SET NULL,
    deadline TIMESTAMP
    );
`)
	require.NoError(t, err)

	t.Cleanup(func() {
		conn.Exec(`DROP TABLE IF EXISTS tasks; DROP TABLE IF EXISTS users;`)
		conn.Close()
	})
}

func TestTaskRepo_AddAndGet(t *testing.T) {
	setupTestDB(t)
	repo := postgres.NewTaskRepo()

	task, err := repo.AddTask(model.Task{Title: "test task", Status: "pending"})
	require.NoError(t, err)
	assert.Equal(t, "test task", task.Title)
	assert.NotZero(t, task.TaskID)

	got, err := repo.GetTaskByID(task.TaskID)
	require.NoError(t, err)
	assert.Equal(t, task.Title, got.Title)
}

func TestTaskRepo_GetByID_NotFound(t *testing.T) {
	setupTestDB(t)
	repo := postgres.NewTaskRepo()

	_, err := repo.GetTaskByID(999)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestTaskRepo_UpdateTask(t *testing.T) {
	setupTestDB(t)
	repo := postgres.NewTaskRepo()

	task, _ := repo.AddTask(model.Task{Title: "old", Status: "pending"})
	updated, err := repo.UpdateTask(task.TaskID, model.Task{Title: "new", Status: "pending"})
	require.NoError(t, err)
	assert.Equal(t, "new", updated.Title)
}

func TestTaskRepo_DeleteTask(t *testing.T) {
	setupTestDB(t)
	repo := postgres.NewTaskRepo()

	task, _ := repo.AddTask(model.Task{Title: "to delete", Status: "pending"})
	err := repo.DeleteTask(task.TaskID)
	assert.NoError(t, err)

	err = repo.DeleteTask(task.TaskID)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestTaskRepo_GetTasksByStatus(t *testing.T) {
	setupTestDB(t)
	repo := postgres.NewTaskRepo()

	repo.AddTask(model.Task{Title: "t1", Status: "pending"})
	repo.AddTask(model.Task{Title: "t2", Status: "done"})
	repo.AddTask(model.Task{Title: "t3", Status: "pending"})

	tasks, err := repo.GetTasksByStatus("pending")
	require.NoError(t, err)
	assert.Len(t, tasks, 2)
}
