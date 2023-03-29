package domain

import (
	"time"
)

type Product struct {
	ProductID   string    `json:"productID"`
	ProductName string    `json:"productName"`
	Description string    `json:"description"`
	Quantity    int64     `json:"quantity"`
	Price       int64     `json:"price"`
	RRP         int64     `json:"rrp"`
	Designer    string    `json:"designer"`
	FitNotes    string    `json:"fitNotes"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func CreateProduct(productID string, productName string, description string, quantity int64, price int64, rrp int64, designer string, fitNotes string) (product Product, err error) {

	utcNow := time.Now().UTC()

	product = Product{
		ProductID:   productID,
		ProductName: productName,
		Description: description,
		Quantity:    quantity,
		Price:       price,
		RRP:         rrp,
		Designer:    designer,
		FitNotes:    fitNotes,
		CreatedAt:   utcNow,
		UpdatedAt:   utcNow,
	}

	return product, nil
}
