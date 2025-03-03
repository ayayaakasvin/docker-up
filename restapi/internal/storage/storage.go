package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"restapi/internal/config"
	"restapi/internal/errorset"
	"restapi/internal/lib/hashtool"
	"restapi/internal/models"
	migrationTool "restapi/internal/storage/migrationTool"

	"github.com/lib/pq"
)

// PostgresStorage implements the Storage interface for PostgreSQL
type PostgresStorage struct {
	db *sql.DB
	config *config.Config
}

// NewPostgresStorage creates a new PostgresStorage
func NewPostgresStorage(cfg *config.Config) (*PostgresStorage, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DatabaseName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migration tool
	if err := migrationTool.ExecuteMigrationTool(
		db, cfg.MigrationPath, migrationTool.UpMigration); err != nil {
		return nil, fmt.Errorf("failed to run migration tool: %w", err)
	}

	return &PostgresStorage{db: db, config: cfg}, nil
}

// SaveUser inserts a new user record into the PostgreSQL database
func (ps *PostgresStorage) SaveUser(username, password string) (int64, error) {
	hashedPassword, err := hashtool.BcryptHashing(password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	stmt, err := ps.db.Prepare("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var userID int64
	err = stmt.QueryRow(username, hashedPassword).Scan(&userID)
	if err != nil {
        if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
            return 0, fmt.Errorf("username already exists")
        }
        return 0, fmt.Errorf("failed to execute statement: %w", err)
    }

	return userID, nil
}

// GetUserByID retrieves a record from the PostgreSQL database by key
func (ps *PostgresStorage) GetUserByID(id int64) (*models.User, error) {
	stmt, err := ps.db.Prepare("SELECT id, username, password, created_at FROM users WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRow(id).Scan(&user.ID, &user.UserName, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorset.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	return &user, nil
}

// UsernameExists checks if a record with the given username exists in the PostgreSQL database
func (ps *PostgresStorage) UsernameExists(name string) (bool, error) {
	stmt, err := ps.db.Prepare("SELECT 1 FROM users WHERE username = $1")
	if err != nil {
		return false, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var exists bool
	err = stmt.QueryRow(name).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to execute statement: %w", err)
	}

	return exists, nil
}

// UpdateUser updates a record in the PostgreSQL database
func (ps *PostgresStorage) UpdateUserPassword(id int64, password string) error {
	stmt, err := ps.db.Prepare("UPDATE users SET password = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var hashedPassword string
	if hashedPassword, err = hashtool.BcryptHashing(password); err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	result, err := stmt.Exec(hashedPassword, id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errorset.ErrUserNotFound
	}

	return nil
}

// DeleteUser deletes a record from the PostgreSQL database
func (ps *PostgresStorage) DeleteUser(id int64) error {
	stmt, err := ps.db.Prepare("DELETE FROM users WHERE id = $1")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	} else if rowsAffected == 0 {
		return errorset.ErrUserNotFound
	}

	return nil
}

// SaveTask inserts a new task record into the PostgreSQL database
func (ps *PostgresStorage) SaveTask(userId int64, content string) (int64, error) {
	stmt, err := ps.db.Prepare("INSERT INTO tasks (user_id, task_content) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var taskID int64
	err = stmt.QueryRow(userId, content).Scan(&taskID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == errorset.ErrForeignKeyConstraintViolation {
			return 0, errorset.ErrUserNotFound
		}

		return 0, fmt.Errorf("failed to execute statement: %w", err)
	}

	return taskID, nil
}

// GetTasksByUserID retrieves a record from the PostgreSQL database by key
func (ps *PostgresStorage) GetTasksByUserID(userID int64) ([]*models.Task, error) {
	if _, err := ps.GetUserByID(userID); err != nil {
		if errors.Is(err, errorset.ErrUserNotFound) {
			return nil, errorset.ErrUserNotFound
		}

		return nil, err
	}

	stmt, err := ps.db.Prepare("SELECT id, user_id, task_content, created_at FROM tasks WHERE user_id = $1")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ID, &task.UserID, &task.TaskContent, &task.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

// UpdateTask updates a record in the PostgreSQL database
func (ps *PostgresStorage) UpdateTaskContent(task_id int64, content string) error {
	stmt, err := ps.db.Prepare("UPDATE tasks SET task_content = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}

	_, err = stmt.Exec(content, task_id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

// DeleteTask deletes a record from the PostgreSQL database
func (ps *PostgresStorage) DeleteTask(task_id int64) error {
	stmt, err := ps.db.Prepare("DELETE FROM tasks WHERE id = $1")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(task_id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

// Ping checks the connection to the PostgreSQL database
func (ps *PostgresStorage) Ping() error {
	return ps.db.Ping()
}

// Close closes the connection to the PostgreSQL database
func (ps *PostgresStorage) Close() error {
	return ps.db.Close()
}

// Reset clears all records from the PostgreSQL database
func (ps *PostgresStorage) Reset() error {
	err := migrationTool.ExecuteMigrationTool(ps.db, ps.config.MigrationPath, migrationTool.ResetMigration)
	if err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}

	return nil
}