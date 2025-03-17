package storage

import "github.com/ayayaakasvin/auth/internal/models/user"

type Storage interface {
	AuthenticateUser 	(username, password string)	(string, error) // should return either jwt token or error
	GetUser				(username string)			(*user.User, error) // helps with fetching data about user

	Close 					()							(error)
	Ping					()							(error)
}