package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func New() Repository {
	url, err := getConnectionString()
	if err != nil {
		// This will not be a connection error, but a DSN parse error or another initialization error.
		log.Fatal(err)
		panic(err)
	}

	conn, err := sqlx.Connect("postgres", url)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or another initialization error.
		log.Fatal(err)
		panic(err)
	}
	return Repository{db: conn}
}

func (s *Repository) IsHealthy(ctx context.Context) bool {
	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
		return false
	}

	return true
}

func getConnectionString() (string, error) {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		return "", err
	}

	// Get PostgreSQL connection details from environment variables
	host := os.Getenv("DATABASE_HOST")
	port := "5432"
	user := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	// Construct PostgreSQL connection string
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname), nil
}
