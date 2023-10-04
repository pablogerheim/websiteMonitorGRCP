package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	ENV  = os.Getenv("ENVIRONMENT")
	PORT = os.Getenv("PORT")

	// SPANNER
	SPANNER_PROJECT  = os.Getenv("SPANNER_PROJECT")
	SPANNER_INSTANCE = os.Getenv("SPANNER_INSTANCE")
	SPANNER_DATABASE = os.Getenv("SPANNER_DATABASE")
	SPANNER_USER     = os.Getenv("SPANNER_USER")
	SPANNER_HOST     = os.Getenv("SPANNER_HOST")
	SPANNER_PASSWORD = os.Getenv("SPANNER_PASSWORD")
	SPANNER_PORT     = os.Getenv("SPANNER_PORT")
)
