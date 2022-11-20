package config

import (
	"golang.org/x/oauth2"
	"os"
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

	return Configuration{
		DatabaseName:      os.Getenv("DB_NAME"),
		DatabaseHost:      os.Getenv("DB_HOST"),
		DatabaseUser:      os.Getenv("DB_USER"),
		DatabasePassword:  os.Getenv("DB_PASSWORD"),
		MigrateToVersion:  migrateToVersion,
		MigrationLocation: migrationLocation,
		AccessSecret:      os.Getenv("ACCESS_SECRET"),
		RefreshSecret:     os.Getenv("REFRESH_SECRET"),
		OAUTH:	LoadOAUTHConfiguration(),
	}
}
