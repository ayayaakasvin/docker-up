package user

import "time"

type User struct {
	UserID    int64     `json:"userId"`
	UserName  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}