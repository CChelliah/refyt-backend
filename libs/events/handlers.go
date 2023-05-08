package events

type HandlerType string

var (
	CustomerHandler = HandlerType("customer")
	ProductHandler  = HandlerType("product")
)
