package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPathEnvKey = "CONFIG_PATH"
)

type Config struct {
	Env 			string 			`yaml:"env" 			env-default:"development"`
	MigrationPath 	string 			`yaml:"migrationPath" 	env-required:"true"`
	Database 		StorageConfig 	`yaml:"database" 										env-required:"true"`
	HTTPServer 						`yaml:"http_server"										env-required:"true"`
}

type StorageConfig struct {
	Host 			string 			`yaml:"host" 			env-default:"localhost"	`
	Port 			string 			`yaml:"port" 			env-default:"5432" 		`
	DatabaseName 	string 			`yaml:"databaseName" 	env-default:"postgres" 	`
	User 			string 			`yaml:"user" 			env-default:"postgres" 	`
	Password 		string 			`yaml:"password" 		env-default:"1488" 		`
}

type HTTPServer struct {
	Address 		string 			`yaml:"address"			env-default:"localhost:8080`
	Timeout 		time.Duration	`yaml:"timeout" 										env-required:"true"`
	IddleTimeout	time.Duration	`yaml:"iddle_timeout" 									env-required:"true"`
}

func MustLoadConfig () *Config {
	configPath := os.Getenv(configPathEnvKey)
	if configPath == "" {
		log.Fatalf("%s is not set up", configPathEnvKey)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist: %s", configPath, err.Error())
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config file: %s", err.Error())
	}

	return &cfg
}