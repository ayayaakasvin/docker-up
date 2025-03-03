package errorset

import (
	"errors"

	"github.com/lib/pq"
)

var (
	ErrUserNotFound        							= errors.New("user not found")
	ErrTaskNotFound        							= errors.New("task not found")
	ErrDuplicateUser       							= errors.New("duplicate user")
	ErrInvalidCredentials  							= errors.New("invalid credentials")
	ErrInvalidPassword								= errors.New("invalid password")
	ErrForeignKeyConstraintViolation pq.ErrorCode 	= "23503"
)