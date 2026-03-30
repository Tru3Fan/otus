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
	query := `INSERT INTO users (username, email, telegram_user_id, telegram_username) VALUES ($1, $2, NULLIF($3, 0), NULLIF($4, ''))  RETURNING id, username, email, telegram_user_id, telegram_username`
	row := db.PostgresDB.QueryRow(query, u.Username, u.Email, u.TelegramUserID, u.TelegramUsername)

	var created model.User
	var tgID sql.NullInt64
	var tgUsername sql.NullString
	err := row.Scan(&created.UserID, &created.Username, &created.Email, &tgID, &tgUsername)
	if err != nil {
		return model.User{}, err
	}
	if tgID.Valid {
		created.TelegramUserID = tgID.Int64
	}
	if tgUsername.Valid {
		created.TelegramUsername = tgUsername.String
	}
	return created, nil
}

func (r *UserRepo) GetUserByID(id int) (model.User, error) {
	query := `SELECT id, username, email, telegram_user_id, telegram_username FROM users WHERE id = $1`
	row := db.PostgresDB.QueryRow(query, id)

	var u model.User
	var tgID sql.NullInt64
	var tgUsername sql.NullString
	err := row.Scan(&u.UserID, &u.Username, &u.Email, &tgID, &tgUsername)
	if err == sql.ErrNoRows {
		return model.User{}, repository.ErrNotFound
	}
	if err != nil {
		return model.User{}, err
	}
	if tgID.Valid {
		u.TelegramUserID = tgID.Int64
	}
	if tgUsername.Valid {
		u.TelegramUsername = tgUsername.String
	}
	return u, nil
}

func (r *UserRepo) GetAllUsers() ([]model.User, error) {
	query := `SELECT id, username, email, telegram_user_id, telegram_username FROM users`
	rows, err := db.PostgresDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var tgUserID sql.NullInt64
		var tgUsername sql.NullString
		if err := rows.Scan(&u.UserID, &u.Username, &u.Email, &tgUserID, &tgUsername); err != nil {
			return nil, err
		}
		if tgUserID.Valid {
			u.TelegramUserID = tgUserID.Int64
		}
		if tgUsername.Valid {
			u.TelegramUsername = tgUsername.String
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) UpdateUser(id int, updated model.User) (model.User, error) {
	query := `UPDATE users SET username = $1, email = $2, telegram_user_id = NULLIF($3, 0), telegram_username = NULLIF($4, '') WHERE id = $5 RETURNING id, username, email, telegram_user_id, telegram_username`
	row := db.PostgresDB.QueryRow(query, updated.Username, updated.Email, updated.TelegramUserID, updated.TelegramUsername, id)

	var u model.User
	var tgID sql.NullInt64
	var tgUsername sql.NullString
	err := row.Scan(&u.UserID, &u.Username, &u.Email, &tgID, &tgUsername)
	if err == sql.ErrNoRows {
		return model.User{}, repository.ErrNotFound
	}
	if err != nil {
		return model.User{}, err
	}
	if tgID.Valid {
		u.TelegramUserID = tgID.Int64
	}
	if tgUsername.Valid {
		u.TelegramUsername = tgUsername.String
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

func (r *UserRepo) GetUserByTelegramID(telegramID int64) (model.User, error) {
	query := `SELECT id, username, email, telegram_user_id, telegram_username FROM users WHERE telegram_user_id = $1`
	row := db.PostgresDB.QueryRow(query, telegramID)

	var u model.User
	var tgID sql.NullInt64
	var tgUsername sql.NullString
	err := row.Scan(&u.UserID, &u.Username, &u.Email, &tgID, &tgUsername)
	if err == sql.ErrNoRows {
		return model.User{}, repository.ErrNotFound
	}
	if err != nil {
		return model.User{}, err
	}
	if tgID.Valid {
		u.TelegramUserID = tgID.Int64
	}
	if tgUsername.Valid {
		u.TelegramUsername = tgUsername.String
	}
	return u, nil
}

func (r *TaskRepo) AddTask(t model.Task) (model.Task, error) {
	query := `INSERT INTO tasks (title, user_id, status, assigned_by) VALUES ($1, NULLIF($2, 0), $3, NULLIF($4, 0)) RETURNING id, title, user_id, status, assigned_by`
	row := db.PostgresDB.QueryRow(query, t.Title, t.UserID, t.Status, t.AssignedBy)

	var created model.Task
	var userID sql.NullInt64
	var assignedBy sql.NullInt64
	err := row.Scan(&created.TaskID, &created.Title, &userID, &created.Status, &assignedBy)
	if err != nil {
		return model.Task{}, err
	}
	if userID.Valid {
		created.UserID = int(userID.Int64)
	}
	if assignedBy.Valid {
		created.AssignedBy = int(assignedBy.Int64)
	}
	return created, nil
}

func (r *TaskRepo) GetTaskByID(id int) (model.Task, error) {
	query := `SELECT id, title, user_id, status, assigned_by FROM tasks WHERE id = $1`
	row := db.PostgresDB.QueryRow(query, id)

	var t model.Task
	var userID sql.NullInt64
	var assignedBy sql.NullInt64
	err := row.Scan(&t.TaskID, &t.Title, &userID, &t.Status, &assignedBy)
	if err == sql.ErrNoRows {
		return model.Task{}, repository.ErrNotFound
	}
	if err != nil {
		return model.Task{}, err
	}
	if userID.Valid {
		t.UserID = int(userID.Int64)
	}
	if assignedBy.Valid {
		t.AssignedBy = int(assignedBy.Int64)
	}
	return t, nil
}

func (r *TaskRepo) GetAllTasks() ([]model.Task, error) {
	query := `SELECT id, title, user_id, status, assigned_by FROM tasks`
	rows, err := db.PostgresDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	var userID sql.NullInt64
	var assignedBy sql.NullInt64
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.TaskID, &t.Title, &userID, &t.Status, &assignedBy); err != nil {
			return nil, err
		}
		if userID.Valid {
			t.UserID = int(userID.Int64)
		}
		if assignedBy.Valid {
			t.AssignedBy = int(assignedBy.Int64)
		}
		tasks = append(tasks, t)

	}
	return tasks, nil
}

func (r *TaskRepo) UpdateTask(id int, updated model.Task) (model.Task, error) {
	query := `UPDATE tasks SET title = $1, user_id = NULLIF($2, 0), status = $3, assigned_by = NULLIF($4, 0)  WHERE id = $5 RETURNING id, title, user_id, status, assigned_by`
	row := db.PostgresDB.QueryRow(query, updated.Title, updated.UserID, updated.Status, updated.AssignedBy, id)

	var t model.Task
	var userID sql.NullInt64
	var assignedBy sql.NullInt64
	err := row.Scan(&t.TaskID, &t.Title, &userID, &t.Status, &assignedBy)
	if err == sql.ErrNoRows {
		return model.Task{}, repository.ErrNotFound
	}
	if err != nil {
		return model.Task{}, err
	}
	if userID.Valid {
		t.UserID = int(userID.Int64)
	}
	if assignedBy.Valid {
		t.AssignedBy = int(assignedBy.Int64)
	}
	return t, nil
}

func (r *TaskRepo) GetTasksByUserID(userID int) ([]model.Task, error) {
	query := `SELECT id, title, user_id, status, assigned_by FROM tasks WHERE  user_id = $1`
	rows, err := db.PostgresDB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	var assignedBy sql.NullInt64
	for rows.Next() {
		var t model.Task
		var uid sql.NullInt64
		if err := rows.Scan(&t.TaskID, &t.Title, &uid, &t.Status, &assignedBy); err != nil {
			return nil, err
		}
		if uid.Valid {
			t.UserID = int(uid.Int64)
		}
		if assignedBy.Valid {
			t.AssignedBy = int(assignedBy.Int64)
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

	taskQuery := `INSERT INTO tasks (title, user_id, status) VALUES ($1, $2, $3) RETURNING id, title, user_id, status`
	row = tx.QueryRow(taskQuery, t.Title, u.UserID, t.Status)
	if err := row.Scan(&t.TaskID, &t.Title, &t.UserID, &t.Status); err != nil {
		return model.User{}, model.Task{}, err
	}

	if err := tx.Commit(); err != nil {
		return model.User{}, model.Task{}, err
	}
	return u, t, nil
}

func (r *TaskRepo) GetTasksByStatus(status string) ([]model.Task, error) {
	query := `SELECT id, title, user_id, status, assigned_by FROM tasks WHERE status = $1`
	rows, err := db.PostgresDB.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		var uid sql.NullInt64
		var assignedBy sql.NullInt64
		if err := rows.Scan(&t.TaskID, &t.Title, &uid, &t.Status, &assignedBy); err != nil {
			return nil, err
		}
		if uid.Valid {
			t.UserID = int(uid.Int64)
		}
		if assignedBy.Valid {
			t.AssignedBy = int(assignedBy.Int64)
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
