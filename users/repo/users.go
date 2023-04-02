package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"refyt-backend/libs"
	"refyt-backend/users/domain"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type IUserRepository interface {
	InsertUser() (domain.User, error)
	FindUserByID(uid string) (domain.User, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *libs.PostgresDatabase) (userRepo UserRepository) {

	userRepo = UserRepository{
		db: db.Db,
	}

	return userRepo
}

func (repo *UserRepository) FindUserByID(uid string) (user domain.User, err error) {

	err = repo.db.QueryRow(findUserByID, uid).Scan(
		&user.Uid,
		&user.Email,
		&user.StripeCustomerID,
		&user.CustomerNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	fmt.Println("test")

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return user, ErrUserNotFound
	case err != nil:
		return user, err
	}

	return user, nil

}

func (repo *UserRepository) CreateUser(user domain.User, stripeCustomerID string) (domain.User, error) {

	err := repo.db.QueryRow(insertUser,
		user.Uid,
		user.Email,
		stripeCustomerID,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(
		&user.Uid,
		&user.Email,
		&user.StripeCustomerID,
		&user.CustomerNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
