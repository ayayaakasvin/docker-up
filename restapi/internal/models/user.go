package models

import "time"

type User struct {
	ID 			int64  		`json:"user_id"`
	UserName 	string 		`json:"username"`
	Password 	string 		`json:"password"`
	CreatedAt 	time.Time 	`json:"created_at"`
}