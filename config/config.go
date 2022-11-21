package config

import (
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"log"
	"os"
	"path/filepath"
)

type Configuration struct {
	DatabaseName      string
	DatabaseHost      string
	DatabaseUser      string
	DatabasePassword  string
	MigrateToVersion  string
	MigrationLocation string
	AccessSecret      string
	RefreshSecret     string
	OAUTH             oauth2.Config
	RedisHost         string
	RedisPort         string
}

func GetConfiguration() Configuration {
	migrationLocation, set := os.LookupEnv("MIGRATION_LOCATION")
	if !set {
		migrationLocation = "migrations"
	}

	migrateToVersion, set := os.LookupEnv("MIGRATE")
	if !set {
		migrateToVersion = "latest"
	}

	err := godotenv.Load(filepath.Join(".env"))
	if err != nil {
		log.Print(err)
	}

	return Configuration{
		DatabaseName:      os.Getenv("NAME_DB"),
		DatabaseHost:      os.Getenv("HOST_DB"),
		DatabaseUser:      os.Getenv("USER_DB"),
		DatabasePassword:  os.Getenv("PASSWORD_DB"),
		MigrateToVersion:  migrateToVersion,
		MigrationLocation: migrationLocation,
		AccessSecret:      os.Getenv("ACCESS_SECRET_"),
		RefreshSecret:     os.Getenv("REFRESH_SECRET_"),
		OAUTH:             LoadOAUTHConfiguration(),
		RedisPort:         "6379",
		RedisHost:         "127.0.0.1",
	}
}
