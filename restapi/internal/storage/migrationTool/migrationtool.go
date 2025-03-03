package migrationtool

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	_ "github.com/lib/pq"
)

type MigrationType int

const (
	UpMigration MigrationType = iota
	DownMigration
	ResetMigration
)

// ExecuteMigrationTool is a tool for migrating the database
func ExecuteMigrationTool(db *sql.DB, migrationPath string, migrationType MigrationType) error {
	sqlScripts, err := parseSQLScriptFromFile(migrationPath, migrationType)
	if err != nil {
		return fmt.Errorf("failed to parse sql script from file: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	for _, script := range sqlScripts {
		if _, err := tx.Exec(script); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute sql script: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func parseSQLScriptFromFile(filePath string, migrationType MigrationType) ([]string, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path is empty")
	}

	// we should read the directory and return the content
	entries, err := os.ReadDir(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("directory is empty")
	}

	var (
		sqlScripts []string
	)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		switch migrationType {
		case UpMigration:
			if isUpMigrationFile(entry.Name()) {
				script, err := readSQLScriptFromFile(path.Join(filePath, entry.Name()))
				if err != nil {
					return nil, fmt.Errorf("failed to read sql script from file: %w", err)
				}

				sqlScripts = append(sqlScripts, script)
			}

		case DownMigration:
			if isDownMigrationFile(entry.Name()) {
				script, err := readSQLScriptFromFile(path.Join(filePath, entry.Name()))
				if err != nil {
					return nil, fmt.Errorf("failed to read sql script from file: %w", err)
				}

				sqlScripts = append(sqlScripts, script)
			}

		case ResetMigration:
			if isResetMigrationFile(entry.Name()) {
				script, err := readSQLScriptFromFile(path.Join(filePath, entry.Name()))
				if err != nil {
					return nil, fmt.Errorf("failed to read sql script from file: %w", err)
				}

				sqlScripts = append(sqlScripts, script)
			}
			
		default:
			return nil, fmt.Errorf("unknown migration type")
		}
	}

	return sqlScripts, nil
}

// to read the content of the file
func readSQLScriptFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}

	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}

// function to check if file is .up.sql format
func isUpMigrationFile(fileName string) bool {
	return strings.HasSuffix(fileName, ".up.sql")
}

// function to check if file is .down.sql format
func isDownMigrationFile(fileName string) bool {
	return strings.HasSuffix(fileName, ".down.sql")
}

// function to check if file is .reset.sql format
func isResetMigrationFile(fileName string) bool {
	return strings.HasSuffix(fileName, ".reset.sql")
}