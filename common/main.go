package common

import (
	"database/sql"
	postgresdb "refyt-backend/common/database"
)

type IEnv interface {
	*sql.DB
}

type Env struct {
	Db *sql.DB
}

func NewEnv() (env *Env) {

	env = &Env{
		Db: postgresdb.NewPostgresDatabase(),
	}

	return env
}
