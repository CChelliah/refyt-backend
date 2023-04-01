package postgresdb

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func NewPostgresDatabase() (database *sql.DB, err error) {

	user, exists := os.LookupEnv("DB_USER")

	if !exists {
		return database, fmt.Errorf("unable to find database user")
	}

	password, exists := os.LookupEnv("DB_PASSWORD")

	if !exists {
		return database, fmt.Errorf("unable to find database password")
	}

	host, exists := os.LookupEnv("DB_HOST")

	if !exists {
		return database, fmt.Errorf("unable to find database host")
	}

	port, exists := os.LookupEnv("DB_PORT")

	if !exists {
		return database, fmt.Errorf("unable to find port")
	}

	db, err := sql.Open("postgres", fmt.Sprintf("postgresql://%s:%s@%s:%s/postgres?sslmode=disable", user, password, host, port))

	if err != nil {
		return db, err
	}

	return db, nil
}
