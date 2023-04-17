package sendgrid

import (
	"github.com/sendgrid/sendgrid-go"
	"os"
)

type Sender struct {
	client *sendgrid.Client
}

func NewSender() (emailService Sender) {

	sendGridKey, exists := os.LookupEnv("SENDGRID_API_KEY")

	if !exists {
		panic("unable to find sendgrid API Key")
	}

	client := sendgrid.NewSendClient(sendGridKey)

	return Sender{client: client}
}

func (e *Sender) SendWelcomeEmail() (err error) {

	return nil
}

func (e *Sender) SendOrderConfirmationEmail() (err error) {

	return nil
}

func (e *Sender) SendShippedOrderEmail() (err error) {

	return nil
}

func (e *Sender) SendOrderPickUpEmail() (err error) {

	return nil
}
func (e *Sender) SendShippedReturnEmail() (err error) {

	return nil
}

func (e *Sender) SendPickupReturnEmail() (err error) {

	return nil
}

func (e *Sender) SendNoticeOfLateFeeEmail() error {

	return nil
}
