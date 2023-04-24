package email

import "refyt-backend/billing/model"

type EmailService interface {
	SendWelcomeEmail(toEmailAddress string) (err error)
	SendOrderConfirmationEmail(toEmailAddress string, productBookings []model.ProductBooking) (err error)
	SendNoticeOfLateFeeEmail() (err error)
	SendShippedOrderEmail() (err error)
	SendPickupReturnReminder() (err error)
	SendOrderPickUpEmail() (err error)
	SendShippingReturnReminder() (err error)
	SendPickupReturnEmail() (err error)
}
