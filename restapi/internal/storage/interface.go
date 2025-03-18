package storage

import (
	"restapi/internal/models/user"
	"restapi/internal/models/task"
)

type Storage interface {
	SaveUser(username, password string) (int64, error)
	GetUserByID(id int64) (*user.User, error)
	UsernameExists(name string) (bool, error)
	UpdateUserPassword(id int64, password string) error
	DeleteUser(id int64) error

	SaveTask(userId int64, content string) (int64, error)
	GetTasksByUserID(userID int64) ([]*task.Task, error)
	GetTaskByTaskID(taskID int64) (*task.Task, error)
	UpdateTaskContent(task_id int64, content string) error
	DeleteTask(task_id int64) error

	Ping() error
	Close() error
}
