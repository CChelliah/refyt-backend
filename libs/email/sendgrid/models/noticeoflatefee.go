package models

type NoticeOfLateFee struct {
	Customer Customer  `json:"customer"`
	Bookings []Booking `json:"bookings"`
	LateFee  LateFee   `json:"lateFee"`
}

type LateFee struct {
	Amount      float64 `json:"amount"`
	OverdueDays int     `json:"overdueDays"`
}
