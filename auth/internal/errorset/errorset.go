package errorset

import (
	"errors"
)

var (
	ErrEmptyUserModel 		error = errors.New("Empty User Model")
	ErrBindRequest 			error = errors.New("failed to bind request")
	ErrRequestSend			error = errors.New("failed to do proxy request")
	ErrRequestCreate		error = errors.New("failed to create proxy request")
	ErrAuthentificateUser	error = errors.New("failed to authentificate user")
	ErrUserNotFound        	error = errors.New("user not found")
)