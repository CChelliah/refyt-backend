package repo

import (
	"database/sql"
	"errors"
	"refyt-backend/customers/domain"
	"refyt-backend/libs"
	"time"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
)

type ICustomerRepository interface {
	FindCustomerByID(customerID string) (customer domain.Customer, err error)
	InsertCustomer(customer domain.Customer, stripeCustomerID string) (domain.Customer, error)
	UpdateCustomer(customer domain.Customer) (domain.Customer, error)
}

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *libs.PostgresDatabase) (customerRepo *CustomerRepository) {

	customerRepo = &CustomerRepository{
		db: db.Db,
	}

	return customerRepo
}

func (repo *CustomerRepository) FindCustomerByID(customerID string) (customer domain.Customer, err error) {

	err = repo.db.QueryRow(findCustomerByID, customerID).Scan(
		&customer.CustomerID,
		&customer.Email,
		&customer.StripeCustomerID,
		&customer.CustomerNumber,
		&customer.StripeConnectID,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return customer, ErrCustomerNotFound
	case err != nil:
		return customer, err
	}
	return customer, nil

}

func (repo *CustomerRepository) InsertCustomer(customer domain.Customer, stripeCustomerID string) (domain.Customer, error) {

	err := repo.db.QueryRow(insertCustomer,
		customer.CustomerID,
		customer.Email,
		stripeCustomerID,
		customer.CreatedAt,
		customer.UpdatedAt,
	).Scan(
		&customer.CustomerID,
		&customer.Email,
		&customer.StripeCustomerID,
		&customer.CustomerNumber,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		return domain.Customer{}, err
	}

	return customer, nil
}

func (repo *CustomerRepository) UpdateCustomer(customer domain.Customer) (domain.Customer, error) {

	utcNow := time.Now().UTC()

	err := repo.db.QueryRow(updateCustomer,
		&customer.StripeConnectID,
		&utcNow,
		&customer.CustomerID,
	).Scan(
		&customer.CustomerID,
		&customer.Email,
		&customer.StripeCustomerID,
		&customer.CustomerNumber,
		&customer.StripeConnectID,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Customer{}, ErrCustomerNotFound
	case err != nil:
		return domain.Customer{}, err
	}

	return customer, nil
}
