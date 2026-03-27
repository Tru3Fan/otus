package postgres

import (
	"database/sql"
	"otus/internal/db"
	"otus/internal/model"
	"otus/internal/repository"
)

type UserRepo struct{}
type TaskRepo struct{}

func NewUserRepo() repository.UserRepository {
	return &UserRepo{}
}

func NewTaskRepo() repository.TaskRepository {
	return &TaskRepo{}
}

func (r *UserRepo) AddUser(u model.User) (model.User, error) {
	query := `INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id, username, email`
	row := db.PostgresDB.QueryRow(query, u.Username, u.Email)

	var created model.User
	err := row.Scan(&created.UserID, &created.Username, &created.Email)
	if err != nil {
		return model.User{}, err
	}
	return created, nil
}

func (r *UserRepo) GetUserByID(id int) (model.User, error) {
	query := `SELECT id, username, email FROM users WHERE id = $1`
	row := db.PostgresDB.QueryRow(query, id)

	var u model.User
	err := row.Scan(&u.UserID, &u.Username, &u.Email)
	if err == sql.ErrNoRows {
		return model.User{}, repository.ErrNotFound
	}
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}

func (r *UserRepo) GetAllUsers() ([]model.User, error) {
	query := `SELECT id, username, email FROM users`
	rows, err := db.PostgresDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.UserID, &u.Username, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) UpdateUser(id int, updated model.User) (model.User, error) {
	query := `UPDATE users SET username = $1, email = $2 WHERE id = $3 RETURNING id, username, email`
	row := db.PostgresDB.QueryRow(query, updated.Username, updated.Email, id)

	var u model.User
	err := row.Scan(&u.UserID, &u.Username, &u.Email)
	if err == sql.ErrNoRows {
		return model.User{}, repository.ErrNotFound
	}
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}

func (r *UserRepo) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := db.PostgresDB.Exec(query, id)
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *TaskRepo) AddTask(t model.Task) (model.Task, error) {
	query := `INSERT INTO tasks (title, user_id) VALUES ($1, NULLIF($2, 0)) RETURNING id, title, user_id`
	row := db.PostgresDB.QueryRow(query, t.Title, t.UserID)

	var created model.Task
	var userID sql.NullInt64
	err := row.Scan(&created.TaskID, &created.Title, &userID)
	if err != nil {
		return model.Task{}, err
	}
	if userID.Valid {
		created.UserID = int(userID.Int64)
	}
	return created, nil
}

func (r *TaskRepo) GetTaskByID(id int) (model.Task, error) {
	query := `SELECT id, title, user_id FROM tasks WHERE id = $1`
	row := db.PostgresDB.QueryRow(query, id)

	var t model.Task
	var userID sql.NullInt64
	err := row.Scan(&t.TaskID, &t.Title, &userID)
	if err == sql.ErrNoRows {
		return model.Task{}, repository.ErrNotFound
	}
	if err != nil {
		return model.Task{}, err
	}
	if userID.Valid {
		t.UserID = int(userID.Int64)
	}
	return t, nil
}

func (r *TaskRepo) GetAllTasks() ([]model.Task, error) {
	query := `SELECT id, title, user_id FROM tasks`
	rows, err := db.PostgresDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	var userID sql.NullInt64
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.TaskID, &t.Title, &userID); err != nil {
			return nil, err
		}
		if userID.Valid {
			t.UserID = int(userID.Int64)
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *TaskRepo) UpdateTask(id int, updated model.Task) (model.Task, error) {
	query := `UPDATE tasks SET title = $1, user_id = NULLIF($2, 0)  WHERE id = $3 RETURNING id, title, user_id`
	row := db.PostgresDB.QueryRow(query, updated.Title, updated.UserID, id)

	var t model.Task
	var userID sql.NullInt64
	err := row.Scan(&t.TaskID, &t.Title, &userID)
	if err == sql.ErrNoRows {
		return model.Task{}, repository.ErrNotFound
	}
	if err != nil {
		return model.Task{}, err
	}
	if userID.Valid {
		t.UserID = int(userID.Int64)
	}
	return t, nil
}

func (r *TaskRepo) GetTasksByUserID(userID int) ([]model.Task, error) {
	query := `SELECT id, title, user_id FROM tasks WHERE user_id = $1`
	rows, err := db.PostgresDB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		var uid sql.NullInt64
		if err := rows.Scan(&t.TaskID, &t.Title, &uid); err != nil {
			return nil, err
		}
		if uid.Valid {
			t.UserID = int(uid.Int64)
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *TaskRepo) DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	result, err := db.PostgresDB.Exec(query, id)
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *TaskRepo) CreateUserWithTask(u model.User, t model.Task) (model.User, model.Task, error) {
	tx, err := db.PostgresDB.Begin()
	if err != nil {
		return model.User{}, model.Task{}, err
	}
	defer tx.Rollback()

	userQuery := `INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id, username, email`
	row := tx.QueryRow(userQuery, u.Username, u.Email)
	if err := row.Scan(&u.UserID, &u.Username, &u.Email); err != nil {
		return model.User{}, model.Task{}, err
	}

	taskQuery := `INSERT INTO tasks (title, user_id) VALUES ($1, $2) RETURNING id, title, user_id`
	row = tx.QueryRow(taskQuery, t.Title, u.UserID)
	if err := row.Scan(&t.TaskID, &t.Title, &t.UserID); err != nil {
		return model.User{}, model.Task{}, err
	}

	if err := tx.Commit(); err != nil {
		return model.User{}, model.Task{}, err
	}
	return u, t, nil
}
