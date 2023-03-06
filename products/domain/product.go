package domain

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ProductID   string    `json:"productID"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Quantity    int64     `json:"quantity"`
	Price       int64     `json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func CreateProduct(title string, description string, quantity int64, price int64) (product Product, err error) {

	productID := uuid.New()

	utcNow := time.Now().UTC()

	product = Product{
		ProductID:   productID.String(),
		Title:       title,
		Description: description,
		Quantity:    quantity,
		Price:       price,
		CreatedAt:   utcNow,
		UpdatedAt:   utcNow,
	}

	return product, nil
}
