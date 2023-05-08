package domain

import (
	"refyt-backend/libs/events"
	"refyt-backend/libs/events/evdata"
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

func CreateProduct(productID string, productName string, description string, price int64, rrp int64, designer string, fitNotes string, category string, size int64, shippingPrice int64) (product Product, event events.Event, err error) {

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

	event = product.toEvent(events.ProductCreatedEvent)

	return product, event, nil
}

func (p *Product) AddImageUrls(imageUrls []string) (event events.Event) {

	p.ImageUrls = imageUrls

	event = p.toEvent(events.ProductUpdatedEvent)

	return event
}

func (p *Product) Update(productID *string, productName *string, description *string, price *int64, rrp *int64, designer *string, fitNotes *string, category *string, size *int64, shippingPrice *int64, imageUrls *[]string) (event events.Event) {

	utcNow := time.Now().UTC()

	if productID != nil {
		p.ProductID = *productID
	}

	if productName != nil {
		p.Name = *productName
	}

	if description != nil {
		p.Description = *description
	}

	if designer != nil {
		p.Designer = *designer
	}

	if category != nil {
		p.Category = *category
	}

	if fitNotes != nil {
		p.FitNotes = *fitNotes
	}

	if size != nil {
		p.Size = *size
	}

	if rrp != nil {
		p.RRP = *rrp * 100
	}

	if price != nil {
		p.Price = *price * 100
	}

	if shippingPrice != nil {
		p.ShippingPrice = *shippingPrice * 100
	}

	if imageUrls != nil {
		p.ImageUrls = *imageUrls
	}

	p.UpdatedAt = utcNow

	event = p.toEvent(events.ProductUpdatedEvent)

	return event
}

func (p *Product) toEvent(eventName events.EventName) (event events.Event) {

	event = events.Event{
		EventName: eventName,
		Data: evdata.ProductEvent{
			ProductID:     p.ProductID,
			Name:          p.Name,
			Description:   p.Description,
			Designer:      p.Designer,
			Category:      p.Category,
			FitNotes:      p.FitNotes,
			Size:          p.Size,
			RRP:           p.RRP,
			Price:         p.Price,
			ShippingPrice: p.ShippingPrice,
			ImageUrls:     p.ImageUrls,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
		},
	}

	return event
}
