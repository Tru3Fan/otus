package repository

import (
	"database/sql"
	"otus/internal/db"
	"otus/internal/model"
)

func PgAddUser(u model.User) (model.User, error) {
	query := `INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id, username, email`
	row := db.PostgresDB.QueryRow(query, u.Username, u.Email)

	var created model.User
	err := row.Scan(&created.UserID, &created.Username, &created.Email)
	if err != nil {
		return model.User{}, err
	}
	return created, nil
}

func PgGetUserByID(id int) (model.User, error) {
	query := `SELECT id, username, email FROM users WHERE id = $1`
	row := db.PostgresDB.QueryRow(query, id)

	var u model.User
	err := row.Scan(&u.UserID, &u.Username, &u.Email)
	if err == sql.ErrNoRows {
		return model.User{}, ErrNotFound
	}
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}

func PgGetAllUsers() ([]model.User, error) {
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

func PgUpdateUser(id int, updated model.User) (model.User, error) {
	query := `UPDATE users SET username = $1, email = $2 WHERE id = $3 RETURNING id, username, email`
	row := db.PostgresDB.QueryRow(query, updated.Username, updated.Email, id)

	var u model.User
	err := row.Scan(&u.UserID, &u.Username, &u.Email)
	if err == sql.ErrNoRows {
		return model.User{}, ErrNotFound
	}
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}

func PgDeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := db.PostgresDB.Exec(query, id)
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func PgAddTask(t model.Task) (model.Task, error) {
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

func PgGetTaskByID(id int) (model.Task, error) {
	query := `SELECT id, title, user_id FROM tasks WHERE id = $1`
	row := db.PostgresDB.QueryRow(query, id)

	var t model.Task
	var userID sql.NullInt64
	err := row.Scan(&t.TaskID, &t.Title, &userID)
	if err == sql.ErrNoRows {
		return model.Task{}, ErrNotFound
	}
	if err != nil {
		return model.Task{}, err
	}
	if userID.Valid {
		t.UserID = int(userID.Int64)
	}
	return t, nil
}

func PgGetAllTasks() ([]model.Task, error) {
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

func PgUpdateTask(id int, updated model.Task) (model.Task, error) {
	query := `UPDATE tasks SET title = $1, user_id = NULLIF($2, 0)  WHERE id = $3 RETURNING id, title, user_id`
	row := db.PostgresDB.QueryRow(query, updated.Title, updated.UserID, id)

	var t model.Task
	var userID sql.NullInt64
	err := row.Scan(&t.TaskID, &t.Title, &userID)
	if err == sql.ErrNoRows {
		return model.Task{}, ErrNotFound
	}
	if err != nil {
		return model.Task{}, err
	}
	if userID.Valid {
		t.UserID = int(userID.Int64)
	}
	return t, nil
}

func PgGetTasksByUserID(userID int) ([]model.Task, error) {
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

func PgDeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	result, err := db.PostgresDB.Exec(query, id)
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func PgCreateUserWithTask(u model.User, t model.Task) (model.User, model.Task, error) {
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
