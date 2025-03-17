package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ayayaakasvin/auth/internal/config"
	"github.com/ayayaakasvin/auth/internal/errorset"
	"github.com/ayayaakasvin/auth/internal/lib/hashtool"
	"github.com/ayayaakasvin/auth/internal/lib/jwtutil"
	"github.com/ayayaakasvin/auth/internal/models/user"
	"github.com/ayayaakasvin/auth/internal/storage"

	_ "github.com/lib/pq"
)

const (
	postgresDriver = "postgres"
)

type PostgreSQL struct {
	connection *sql.DB
	cfg        config.StorageConfig
}

func NewPostgresStorage(cfg config.StorageConfig) storage.Storage {
	psqlObj := new(PostgreSQL)
	psqlObj.cfg = cfg

	db, err := sql.Open(postgresDriver, connectionString(cfg))
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	psqlObj.connection = db

	return psqlObj
}

// AuthenticateUser implements storage.Storage.
func (p *PostgreSQL) AuthenticateUser(username string, password string) (string, error) {
	userObject, err := p.GetUser(username)
	if err != nil {
		return "", err
	}

	err = hashtool.BcryptCompare(userObject.Password, password)
	if err != nil {
		return "", err
	}

	token, err := jwtutil.GenerateJWT(userObject.UserID)
	if err != nil {
		return "", err
	}

	return token, nil
}


// GetUser implements storage.Storage.
func (p *PostgreSQL) GetUser(username string) (*user.User, error) {
	stmt, err := p.connection.Prepare(`SELECT user_id, username, password, created_at from users WHERE username = $1`)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var userReturnObject *user.User = &user.User{}

	err = stmt.QueryRow(username).Scan(&userReturnObject.UserID, 
		&userReturnObject.UserName, 
		&userReturnObject.Password, 
		&userReturnObject.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorset.ErrUserNotFound
		}
		
		return nil, fmt.Errorf("failed to scan: %v", err)
	}

	return userReturnObject, nil
}

// Close implements storage.Storage.
func (p *PostgreSQL) Close() error {
	return p.connection.Close()
}

// Ping implements storage.Storage.
func (p *PostgreSQL) Ping() error {
	return p.connection.Ping()
}

func connectionString(cfg config.StorageConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DatabaseName,
	)
}
