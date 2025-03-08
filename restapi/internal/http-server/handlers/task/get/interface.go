package get

import "restapi/internal/models"

type TasksGetter interface {
	GetTasksByUserID(userId int64) ([]*models.Task, error)
	GetTaskByTaskID(taskId int64) (*models.Task, error)
}