package email

type EmailService interface {
	SendWelcomeEmail() (err error)
	SendOrderConfirmationEmail() (err error)
	SendNoticeOfLateFeeEmail() (err error)
	SendShippedOrderEmail() (err error)
	SendPickupReturnReminder() (err error)
	SendOrderPickUpEmail() (err error)
	SendShippingReturnReminder() (err error)
	SendPickupReturnEmail() (err error)
}
