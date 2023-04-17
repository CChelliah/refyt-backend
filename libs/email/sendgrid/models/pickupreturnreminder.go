package models

type PickupReturnReminder struct {
	OrderNumber string     `json:"orderNumber"`
	Customer    Customer   `json:"customer"`
	Bookings    []Booking  `json:"bookings"`
	PickupInfo  PickupInfo `json:"pickupInfo"`
}

type PickupInfo struct {
	Address     Address `json:"address"`
	PhoneNumber string  `json:"phoneNumber"`
}

type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}
