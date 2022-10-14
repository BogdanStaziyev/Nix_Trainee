package config

import "os"

type Configuration struct {
	DatabaseName      string
	DatabaseHost      string
	DatabaseUser      string
	DatabasePassword  string
	MigrateToVersion  string
	MigrationLocation string
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
		DatabaseName:      `coordinate`,
		DatabaseHost:      `localhost:54322`,
		DatabaseUser:      `postgres`,
		DatabasePassword:  `password`,
		MigrateToVersion:  migrateToVersion,
		MigrationLocation: migrationLocation,
	}
}
