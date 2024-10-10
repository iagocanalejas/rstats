package db

import (
	"fmt"
	"os"

	"github.com/iagocanalejas/rstats/internal/utils/assert"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func New() Repository {
	conn, err := sqlx.Connect("postgres", getConnectionString())
	assert.NoError(err, "connecting to database")
	return Repository{db: conn}
}

func getConnectionString() string {
	err := godotenv.Load(".env")
	assert.NoError(err, "loading .env file")

	host := os.Getenv("DATABASE_HOST")
	port := "5432"
	user := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)
}
