package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/iagocanalejas/regatas/internal/db"
	races "github.com/iagocanalejas/regatas/pkg/races"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	connStr, err := getConnectionString()
	if err != nil {
		log.Fatal(fmt.Errorf("error loading .env file: %w", err))
		return
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal(fmt.Errorf("error oppening database: %w", err))
		return
	}

	queries := db.New(conn)
	entities := races.NewService(queries)

	fetchFlag, err := entities.GetFlags(ctx)
	if err != nil {
		return
	}

	// prints true
	log.Println(fetchFlag)
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
