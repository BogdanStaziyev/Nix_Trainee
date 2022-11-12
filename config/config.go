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
		DatabaseName:      `nix_db`,
		DatabaseHost:      `localhost:8081`,
		DatabaseUser:      `admin`,
		DatabasePassword:  `password`,
		MigrateToVersion:  migrateToVersion,
		MigrationLocation: migrationLocation,
		AccessSecret:      "access",
		RefreshSecret:     "refresh",
		OAUTH:             LoadOAUTHConfiguration(),
	}
}
