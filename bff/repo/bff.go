package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"refyt-backend/bff/model"
	"refyt-backend/libs"
)

type IBffRepo interface {
	Find(categories []string, sizes []int64) (product []model.Product, err error)
	FindByUserID(userID string) (product []model.Product, err error)
	FindProductByID(productID string) (product model.Product, err error)
}

type BffRepository struct {
	db *sql.DB
}

func NewBffRepository(db *libs.PostgresDatabase) (bffRepo BffRepository) {

	bffRepo = BffRepository{
		db: db.Db,
	}

	return bffRepo
}

var (
	ErrProductNotFound = errors.New("product not found")
)

func (repo *BffRepository) FindProductByID(productID string) (product model.Product, err error) {

	err = repo.db.QueryRow(findProductByID,
		productID,
	).Scan(
		&product.ProductID,
		&product.Name,
		&product.Description,
		&product.Designer,
		&product.Category,
		&product.FitNotes,
		&product.Size,
		&product.RRP,
		&product.Price,
		&product.ShippingPrice,
		pq.Array(&product.ImageUrls),
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.SellerID,
		&product.UserName,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return model.Product{}, ErrProductNotFound
	case err != nil:
		return model.Product{}, err
	}

	return product, nil
}

func (repo *BffRepository) Find(categories []string, sizes []int64) (products []model.Product, err error) {

	products = []model.Product{}
	var parameters []interface{}

	paramCount := 1

	categoryClause := ""
	if len(categories) > 0 {
		categoryClause = fmt.Sprintf("AND category = ANY($%d)", paramCount)
		parameters = append(parameters, pq.Array(categories))
		paramCount++ //nolint
	}

	sizeClause := ""
	if len(sizes) > 0 {
		sizeClause = fmt.Sprintf("AND size = ANY($%d)", paramCount)
		parameters = append(parameters, pq.Array(sizes))
		paramCount++ //nolint
	}

	query := find + " " + categoryClause + " " + sizeClause + ";"

	rows, err := repo.db.Query(query, parameters...)

	if err != nil {
		return []model.Product{}, err
	}

	defer rows.Close()

	for rows.Next() {

		var product model.Product

		err = rows.Scan(
			&product.ProductID,
			&product.Name,
			&product.Description,
			&product.Designer,
			&product.Category,
			&product.FitNotes,
			&product.Size,
			&product.RRP,
			&product.Price,
			&product.ShippingPrice,
			pq.Array(&product.ImageUrls),
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.SellerID,
			&product.UserName,
		)

		if err != nil {
			return []model.Product{}, err
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return []model.Product{}, err
	}

	return products, nil
}

func (repo *BffRepository) FindByUserID(userID string) (products []model.Product, err error) {
	products = []model.Product{}

	rows, err := repo.db.Query(findAllByUserID, userID)

	if err != nil {
		return []model.Product{}, err
	}

	defer rows.Close()

	for rows.Next() {

		var product model.Product

		err = rows.Scan(
			&product.ProductID,
			&product.Name,
			&product.Description,
			&product.Designer,
			&product.Category,
			&product.FitNotes,
			&product.Size,
			&product.RRP,
			&product.Price,
			&product.ShippingPrice,
			pq.Array(&product.ImageUrls),
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return []model.Product{}, err
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return []model.Product{}, err
	}

	return products, nil
}
