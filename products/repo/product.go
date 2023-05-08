package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"refyt-backend/libs"
	"refyt-backend/libs/uow"
	"refyt-backend/products/domain"
	"time"
)

type IProductRepository interface {
	CreateProduct() (domain.Product, error)
	UpdateProduct(productID string, title string, description string, quantity int64, price int64) (product domain.Product, err error)
	FindByID(productID string) (product domain.Product, err error)
	DeleteProduct(productID string) (err error)
	Find(categories []string, sizes []int64) (product []domain.Product, err error)
	FindByUserID(userID string) (product []domain.Product, err error)
}

type ProductRepository struct {
	db *sql.DB
}

var (
	ErrProductNotFound = errors.New("product not found")
)

func NewProductRepository(db *libs.PostgresDatabase) (productRepo ProductRepository) {

	productRepo = ProductRepository{
		db: db.Db,
	}

	return productRepo
}

func (repo *ProductRepository) InsertProduct(ctx context.Context, uow uow.UnitOfWork, product domain.Product, uid string) (domain.Product, error) {

	err := uow.GetTx().QueryRowContext(ctx, createProduct,
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
		&uid,
		pq.Array(&product.ImageUrls),
		&product.CreatedAt,
		&product.UpdatedAt,
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
	)

	if err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

func (repo *ProductRepository) UpdateProduct(ctx context.Context, uow uow.UnitOfWork, productID string, product domain.Product) (newProduct domain.Product, err error) {

	utcNow := time.Now().UTC()

	err = uow.GetTx().QueryRowContext(ctx, updateProduct,
		&productID,
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
		utcNow,
	).Scan(
		&newProduct.ProductID,
		&newProduct.Name,
		&newProduct.Description,
		&newProduct.Designer,
		&newProduct.Category,
		&newProduct.FitNotes,
		&newProduct.Size,
		&newProduct.RRP,
		&newProduct.Price,
		&newProduct.ShippingPrice,
		pq.Array(&newProduct.ImageUrls),
		&newProduct.UpdatedAt,
		&newProduct.CreatedAt,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Product{}, ErrProductNotFound
	case err != nil:
		return domain.Product{}, err
	}

	return newProduct, nil
}

func (repo *ProductRepository) FindProductByID(productID string) (product domain.Product, err error) {

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

func (repo *ProductRepository) Find(categories []string, sizes []int64) (products []domain.Product, err error) {

	products = []domain.Product{}
	var parameters []interface{}

	paramCount := 1

	categoryClause := ""
	if len(categories) > 0 {
		categoryClause = fmt.Sprintf("AND category = ANY($%d)", paramCount)
		parameters = append(parameters, pq.Array(categories))
		paramCount = paramCount + 1
	}

	sizeClause := ""
	if len(sizes) > 0 {
		sizeClause = fmt.Sprintf("AND size = ANY($%d)", paramCount)
		parameters = append(parameters, pq.Array(sizes))
		paramCount = paramCount + 1
	}

	query := find + " " + categoryClause + " " + sizeClause + ";"

	rows, err := repo.db.Query(query, parameters...)

	if err != nil {
		return []domain.Product{}, err
	}

	defer rows.Close()

	for rows.Next() {

		var product domain.Product

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
			return []domain.Product{}, err
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return []domain.Product{}, err
	}

	return products, nil
}
