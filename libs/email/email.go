package email

import (
	"refyt-backend/libs/email/sendgrid/models"
)

type EmailService interface {
	SendWelcomeEmail(toEmailAddress string) (err error)
	SendOrderConfirmationEmail(toEmailAddress string, productBookings []models.ProductBooking) (err error)
	SendNoticeOfLateFeeEmail() (err error)
	SendShippedOrderEmail() (err error)
	SendPickupReturnReminder() (err error)
	SendOrderPickUpEmail() (err error)
	SendShippingReturnReminder() (err error)
	SendPickupReturnEmail() (err error)
}
