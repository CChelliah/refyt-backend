package domain

import (
	"fmt"
	"time"
)

type Product struct {
	ProductID     string    `json:"productID"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Designer      string    `json:"designer"`
	Category      string    `json:"category"`
	FitNotes      string    `json:"fitNotes"`
	Size          int64     `json:"size"`
	RRP           int64     `json:"rrp"`
	Price         int64     `json:"price"`
	ShippingPrice int64     `json:"shippingPrice"`
	ImageUrls     []string  `json:"imageUrls"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func CreateProduct(productID string, productName string, description string, price int64, rrp int64, designer string, fitNotes string, category string, size int64, shippingPrice int64) (product Product, err error) {

	utcNow := time.Now().UTC()

	product = Product{
		ProductID:     productID,
		Name:          productName,
		Description:   description,
		Designer:      designer,
		Category:      category,
		FitNotes:      fitNotes,
		Size:          size,
		RRP:           rrp * 100,
		Price:         price * 100,
		ShippingPrice: shippingPrice * 100,
		CreatedAt:     utcNow,
		UpdatedAt:     utcNow,
	}

	fmt.Println("")

	return product, nil
}

func (p *Product) AddImageUrls(imageUrls []string) {

	p.ImageUrls = imageUrls
}
