package repo

import (
	"database/sql"
	"errors"
	"time"
	"trading-card-app-backend/common"
	"trading-card-app-backend/sellers/domain"
)

var (
	ErrSellerNotFound = errors.New("seller not found")
)

type ISellerRepository interface {
	FindSeller(uid string) (domain.Seller, error)
	UpdateSeller(sellerUpdate domain.Seller) (seller domain.Seller, err error)
}

type SellerRepository struct {
	db *sql.DB
}

func NewSellerRepository(env *common.Env) (sellerRepo SellerRepository) {

	sellerRepo = SellerRepository{
		db: env.Db,
	}

	return sellerRepo
}

func (repo *SellerRepository) FindSeller(uid string) (seller domain.Seller, err error) {

	err = repo.db.QueryRow(findSellerById, uid).Scan(
		&seller.Uid,
		&seller.Email,
		&seller.StripeCustomerID,
		&seller.CustomerNumber,
		&seller.CreatedAt,
		&seller.UpdatedAt,
	)

	if err != nil {
		return seller, err
	}

	return seller, nil
}

func (repo *SellerRepository) UpdateSeller(sellerUpdate domain.Seller) (seller domain.Seller, err error) {

	utcNow := time.Now().UTC()

	err = repo.db.QueryRow(updateSeller,
		&sellerUpdate.ConnectAccountID,
		&utcNow,
		&sellerUpdate.Uid,
	).Scan(
		&seller.Uid,
		&seller.Email,
		&seller.StripeCustomerID,
		&seller.ConnectAccountID,
		&seller.CustomerNumber,
		&seller.UpdatedAt,
		&seller.CreatedAt,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Seller{}, ErrSellerNotFound
	case err != nil:
		return domain.Seller{}, err
	}

	return seller, nil
}
