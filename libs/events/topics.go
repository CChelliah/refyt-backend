package events

type Topic string

var (
	CustomerTopic = Topic("customer")
	BookingTopic  = Topic("booking")
	PaymentTopic  = Topic("payment")
	ShippingTopic = Topic("shipping")
	ProductTopic  = Topic("product")
)

type EventName string

var (
	CustomerCreatedEvent = EventName("customer.created")
	CustomerUpdatedEvent = EventName("customer.updated")

	ProductCreatedEvent = EventName("product.created")
	ProductUpdatedEvent = EventName("product.updated")
	ProductDeletedEvent = EventName("product.deleted")

	BookingCreatedEvent   = EventName("booking.created")
	BookingScheduledEvent = EventName("booking.scheduled")
	BookingUpdatedEvent   = EventName("booking.updated")

	PaymentCreatedEvent   = EventName("payment.created")
	PaymentSucceededEvent = EventName("payment.succeeded")
)
