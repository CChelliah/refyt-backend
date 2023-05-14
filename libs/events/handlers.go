package events

type HandlerType string

var (
	CustomerHandler = HandlerType("customer")
	ProductHandler  = HandlerType("product")
	PaymentHandler  = HandlerType("payment")

	BookingHandler        = HandlerType("booking")
	BookingPaymentHandler = HandlerType("booking.payment")
)
