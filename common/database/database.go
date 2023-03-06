package postgresdb

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func NewPostgresDatabase() (database *sql.DB) {

	db, err := sql.Open("postgres", "postgresql://cavinashchelliah:password@localhost:5432/postgres?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	return db
}
