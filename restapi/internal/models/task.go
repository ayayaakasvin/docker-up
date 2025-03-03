package models

import "time"

type Task struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	TaskContent string    `json:"task_content"`
	CreatedAt   time.Time `json:"created_at"`
}