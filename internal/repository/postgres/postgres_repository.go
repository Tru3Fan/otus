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
	query := `INSERT INTO users (username, telegram_user_id, telegram_username) VALUES ($1, NULLIF($2::bigint, 0), NULLIF($3, ''))  RETURNING id, username, telegram_user_id, telegram_username`
	row := db.PostgresDB.QueryRow(query, u.Username, u.TelegramUserID, u.TelegramUsername)

	var created model.User
	var tgID sql.NullInt64
	var tgUsername sql.NullString
	err := row.Scan(&created.UserID, &created.Username, &tgID, &tgUsername)
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
	query := `SELECT id, username, telegram_user_id, telegram_username FROM users WHERE id = $1`
	row := db.PostgresDB.QueryRow(query, id)

	var u model.User
	var tgID sql.NullInt64
	var tgUsername sql.NullString
	err := row.Scan(&u.UserID, &u.Username, &tgID, &tgUsername)
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
	query := `SELECT id, username, telegram_user_id, telegram_username FROM users`
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
		if err := rows.Scan(&u.UserID, &u.Username, &tgUserID, &tgUsername); err != nil {
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
	query := `UPDATE users SET username = $1, telegram_user_id = NULLIF($2::bigint, 0), telegram_username = NULLIF($3, '') WHERE id = $4 RETURNING id, username, telegram_user_id, telegram_username`
	row := db.PostgresDB.QueryRow(query, updated.Username, updated.TelegramUserID, updated.TelegramUsername, id)

	var u model.User
	var tgID sql.NullInt64
	var tgUsername sql.NullString
	err := row.Scan(&u.UserID, &u.Username, &tgID, &tgUsername)
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
	query := `SELECT id, username, telegram_user_id, telegram_username FROM users WHERE telegram_user_id = $1`
	row := db.PostgresDB.QueryRow(query, telegramID)

	var u model.User
	var tgID sql.NullInt64
	var tgUsername sql.NullString
	err := row.Scan(&u.UserID, &u.Username, &tgID, &tgUsername)
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

func (r *UserRepo) AddPendingUser(username string) error {
	_, err := db.PostgresDB.Exec(`INSERT INTO pending_users (telegram_username) VALUES ($1) ON CONFLICT DO NOTHING`, username)
	return err
}

func (r *UserRepo) IsPendingUser(username string) (bool, error) {
	var exists bool
	err := db.PostgresDB.QueryRow(`SELECT EXISTS(SELECT 1 FROM pending_users WHERE telegram_username = $1)`, username).Scan(&exists)
	return exists, err
}

func (r *UserRepo) DeletePendingUser(username string) error {
	_, err := db.PostgresDB.Exec("DELETE FROM pending_users WHERE telegram_username = $1", username)
	return err
}

func (r *UserRepo) GetUserByTelegramUsername(username string) (model.User, error) {
	query := `SELECT id, username, telegram_user_id, telegram_username FROM users WHERE telegram_username = $1`
	row := db.PostgresDB.QueryRow(query, username)

	var u model.User
	var tgID sql.NullInt64
	var tgUsername sql.NullString
	err := row.Scan(&u.UserID, &u.Username, &tgID, &tgUsername)
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
	query := `INSERT INTO tasks (title,  user_id, status, assigned_by, deadline)
	VALUES ($1, NULLIF($2,0), $3, NULLIF($4, 0), $5) 
	RETURNING id, title, user_id, status, assigned_by, deadline`

	row := db.PostgresDB.QueryRow(query, t.Title, t.UserID, t.Status, t.AssignedBy, t.Deadline)

	var created model.Task
	var userID sql.NullInt64
	var assignedBy sql.NullInt64
	var deadline sql.NullTime
	err := row.Scan(&created.TaskID, &created.Title, &userID, &created.Status, &assignedBy, &deadline)
	if err != nil {
		return model.Task{}, err
	}
	if userID.Valid {
		created.UserID = int(userID.Int64)
	}
	if assignedBy.Valid {
		created.AssignedBy = int(assignedBy.Int64)
	}
	if deadline.Valid {
		created.Deadline = &deadline.Time
	}
	return created, nil
}

func (r *TaskRepo) GetTaskByID(id int) (model.Task, error) {
	query := `SELECT id, title, user_id, status, assigned_by, deadline FROM tasks WHERE id = $1`
	row := db.PostgresDB.QueryRow(query, id)

	var t model.Task
	var userID sql.NullInt64
	var assignedBy sql.NullInt64
	var deadline sql.NullTime
	err := row.Scan(&t.TaskID, &t.Title, &userID, &t.Status, &assignedBy, &deadline)
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
	if deadline.Valid {
		t.Deadline = &deadline.Time
	}
	return t, nil
}

func (r *TaskRepo) GetAllTasks() ([]model.Task, error) {
	query := `SELECT id, title,  user_id, status, assigned_by, deadline FROM tasks`
	rows, err := db.PostgresDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		var t model.Task
		var userID sql.NullInt64
		var assignedBy sql.NullInt64
		var deadline sql.NullTime
		if err := rows.Scan(&t.TaskID, &t.Title, &userID, &t.Status, &assignedBy, &deadline); err != nil {
			return nil, err
		}
		if userID.Valid {
			t.UserID = int(userID.Int64)
		}
		if assignedBy.Valid {
			t.AssignedBy = int(assignedBy.Int64)
		}
		if deadline.Valid {
			t.Deadline = &deadline.Time
		}
		tasks = append(tasks, t)

	}
	return tasks, nil
}

func (r *TaskRepo) UpdateTask(id int, updated model.Task) (model.Task, error) {
	query := `UPDATE tasks SET title = $1,  user_id = NULLIF($2, 0), status = $3, assigned_by = NULLIF($4, 0) WHERE id = $5 RETURNING id, title, user_id, status, assigned_by`
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

func (r *TaskRepo) GetTasksByAuthorID(authorID int) ([]model.Task, error) {
	query := `SELECT id, title, user_id, status, assigned_by FROM tasks WHERE assigned_by = $1`
	rows, err := db.PostgresDB.Query(query, authorID)
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
