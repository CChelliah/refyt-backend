package repo

import (
	"database/sql"
	"refyt-backend/libs"
)

type ISchedulerRepo interface {
}

type SchedulerRepo struct {
	db *sql.DB
}

func NewSchedulerRepo(db *libs.PostgresDatabase) (productRepo SchedulerRepo) {

	productRepo = SchedulerRepo{
		db: db.Db,
	}

	return productRepo
}
