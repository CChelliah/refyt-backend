package libs

import (
	"database/sql"
	postgresdb "refyt-backend/libs/database"
)

type PostgresDatabase struct {
	Db *sql.DB
}

func NewDatabase() (db *PostgresDatabase, err error) {

	postgresDaabase, err := postgresdb.NewPostgresDatabase()

	if err != nil {
		return db, err
	}

	db = &PostgresDatabase{
		Db: postgresDaabase,
	}

	return db, nil
}
