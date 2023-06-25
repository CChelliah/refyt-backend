package model

import (
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
	SellerID      string    `json:"sellerID"`
	UserName      *string   `json:"userName"`
}
