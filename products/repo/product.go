package repo

import (
	"database/sql"
	"errors"
	"time"
	"trading-card-app-backend/common"
	"trading-card-app-backend/products/domain"
)

type IProductRepository interface {
	CreateProduct() (domain.Product, error)
	UpdateProduct(productID string, title string, description string, quantity int64, price int64) (product domain.Product, err error)
	FindByID(productID string) (product domain.Product, err error)
	DeleteProduct(productID string) (err error)
	FindAll() (product []domain.Product, err error)
}

type ProductRepository struct {
	db *sql.DB
}

var (
	ErrProductNotFound = errors.New("product not found")
)

func NewProductRepository(env *common.Env) (productRepo ProductRepository) {

	productRepo = ProductRepository{
		db: env.Db,
	}

	return productRepo
}

func (repo *ProductRepository) CreateProduct(product domain.Product) (domain.Product, error) {

	err := repo.db.QueryRow(createProduct,
		product.ProductID,
		product.Title,
		product.Description,
		product.Quantity,
		product.Price,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(
		&product.ProductID,
		&product.Title,
		&product.Description,
		&product.Quantity,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

func (repo *ProductRepository) UpdateProduct(productID string, title string, description string, quantity int64, price int64) (product domain.Product, err error) {

	utcNow := time.Now().UTC()

	err = repo.db.QueryRow(updateProduct,
		productID,
		title,
		description,
		quantity,
		price,
		utcNow,
	).Scan(
		&product.ProductID,
		&product.Title,
		&product.Description,
		&product.Quantity,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Product{}, ErrProductNotFound
	case err != nil:
		return domain.Product{}, err
	}

	return product, nil
}

func (repo *ProductRepository) FindByID(productID string) (product domain.Product, err error) {

	err = repo.db.QueryRow(findProductByID,
		productID,
	).Scan(
		&product.ProductID,
		&product.Title,
		&product.Description,
		&product.Quantity,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Product{}, ErrProductNotFound
	case err != nil:
		return domain.Product{}, err
	}

	return product, nil
}

func (repo *ProductRepository) DeleteProduct(productID string) (err error) {

	res, err := repo.db.Exec(deleteProduct, productID)

	if err != nil {
		return err
	}

	numDeleted, err := res.RowsAffected()

	switch {
	case numDeleted == 0:
		return ErrProductNotFound
	case err != nil:
		return err
	}

	return nil
}

func (repo *ProductRepository) FindAll() (products []domain.Product, err error) {

	products = []domain.Product{}

	rows, err := repo.db.Query(findAll)

	if err != nil {
		return []domain.Product{}, err
	}

	defer rows.Close()

	for rows.Next() {

		var product domain.Product

		err = rows.Scan(
			&product.ProductID,
			&product.Title,
			&product.Description,
			&product.Quantity,
			&product.Price,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return []domain.Product{}, err
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return []domain.Product{}, err
	}

	return products, nil
}
