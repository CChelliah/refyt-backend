package domain

type Product struct {
	ProductID     string `json:"productID"`
	ShippingPrice int64  `json:"shippingPrice"`
}
