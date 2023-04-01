package repo

import (
	"context"
	"database/sql"
	"errors"
	"refyt-backend/common"
	"refyt-backend/common/uow"
	"refyt-backend/products/domain"
	"time"
)

type IProductRepository interface {
	CreateProduct() (domain.Product, error)
	UpdateProduct(productID string, title string, description string, quantity int64, price int64) (product domain.Product, err error)
	FindByID(productID string) (product domain.Product, err error)
	DeleteProduct(productID string) (err error)
	FindAll() (product []domain.Product, err error)
	FindByUserID(userID string) (product []domain.Product, err error)
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

func (repo *ProductRepository) InsertProduct(ctx context.Context, uow uow.UnitOfWork, product domain.Product, uid string) (domain.Product, error) {

	err := uow.GetTx().QueryRowContext(ctx, createProduct,
		product.ProductID,
		product.ProductName,
		product.Description,
		product.Quantity,
		product.Price,
		product.RRP,
		product.Designer,
		product.FitNotes,
		uid,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(
		&product.ProductID,
		&product.ProductName,
		&product.Description,
		&product.Quantity,
		&product.Price,
		&product.RRP,
		&product.Designer,
		&product.FitNotes,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

func (repo *ProductRepository) UpdateProduct(ctx context.Context, uow uow.UnitOfWork, productID string, title string, description string, quantity int64, price int64) (product domain.Product, err error) {

	utcNow := time.Now().UTC()

	err = uow.GetTx().QueryRowContext(ctx, updateProduct,
		productID,
		title,
		description,
		quantity,
		price,
		utcNow,
	).Scan(
		&product.ProductID,
		&product.ProductName,
		&product.Description,
		&product.Quantity,
		&product.Price,
		&product.RRP,
		&product.Designer,
		&product.FitNotes,
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
		&product.ProductName,
		&product.Description,
		&product.Quantity,
		&product.Price,
		&product.RRP,
		&product.Designer,
		&product.FitNotes,
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

func (repo *ProductRepository) DeleteProduct(ctx context.Context, uow uow.UnitOfWork, productID string) (err error) {

	res, err := uow.GetTx().ExecContext(ctx, deleteProduct, productID)

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
			&product.ProductName,
			&product.Description,
			&product.Quantity,
			&product.Price,
			&product.RRP,
			&product.Designer,
			&product.FitNotes,
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

func (repo *ProductRepository) FindByUserID(userID string) (products []domain.Product, err error) {
	products = []domain.Product{}

	rows, err := repo.db.Query(findAllByUserID, userID)

	if err != nil {
		return []domain.Product{}, err
	}

	defer rows.Close()

	for rows.Next() {

		var product domain.Product

		err = rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.Description,
			&product.Quantity,
			&product.Price,
			&product.RRP,
			&product.Designer,
			&product.FitNotes,
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
