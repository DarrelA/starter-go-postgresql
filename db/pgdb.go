package pgdb

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Db is the database connection pool.
var Db *pgxpool.Pool

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error is occurred  on .env file please check")
	}

	host := os.Getenv("PGDB_HOST")
	port, _ := strconv.Atoi(os.Getenv("PGDB_PORT")) // Convert from int type
	user := os.Getenv("PGDB_USERNAME")
	dbname := os.Getenv("PGDB_NAME")
	pass := os.Getenv("PGDB_PASSWORD")
	sslmode := os.Getenv("PGDB_SSLMODE")

	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s sslmode=%s dbname=%s",
		host, port, user, pass, sslmode, dbname,
	)

	Db, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Printf("Unable to connect to the database: %v\n", err)
		panic(err)
	}
	log.Println("Successfully connected to the Postgres database!")
}
