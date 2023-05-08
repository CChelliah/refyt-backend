package sendgrid

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"os"
	"refyt-backend/libs/email/sendgrid/models"
)

type Sender struct {
	client *sendgrid.Client
}

const (
	FromEmailAddress                 = "hello@therefyt.com"
	FromName                         = "The Refyt"
	WelcomeEmailTemplateID           = "d-6fe09c3d409d49da97594b1f8e5f0e7c"
	OrderConfirmationEmailTemplateID = "d-360d50494313465fa79c8988a10058a2"
)

func NewSender() (emailService *Sender) {

	sendGridKey, exists := os.LookupEnv("SENDGRID_API_KEY")

	if !exists {
		panic("unable to find sendgrid API Key")
	}

	client := sendgrid.NewSendClient(sendGridKey)

	return &Sender{client: client}
}

func (s *Sender) SendWelcomeEmail(toEmailAddress string) (err error) {

	from := mail.NewEmail(FromName, FromEmailAddress)
	to := mail.NewEmail("", toEmailAddress)
	subject := ""

	p := mail.NewPersonalization()
	p.AddTos(to)

	email := mail.NewSingleEmail(from, subject, to, "", "") // empty body and plain text

	email.SetTemplateID(WelcomeEmailTemplateID)

	_, err = s.client.Send(email)

	if err != nil {
		return err
	}

	return err
}

func (s *Sender) SendOrderConfirmationEmail(toEmailAddress string, productBookings []models.ProductBooking) (err error) {

	from := mail.NewEmail(FromName, FromEmailAddress)
	to := mail.NewEmail("", toEmailAddress)
	subject := ""

	p := mail.NewPersonalization()
	p.AddTos(to)

	p.SetDynamicTemplateData("productBookings", productBookings)

	email := mail.NewSingleEmail(from, subject, to, "", "") // empty body and plain text

	email.SetTemplateID(OrderConfirmationEmailTemplateID)

	res, err := s.client.Send(email)

	if err != nil {
		fmt.Printf("Error %s\n", res.Body)
		fmt.Printf("Error %d\n", res.StatusCode)
		fmt.Printf("Error %s\n", err.Error())
		return err
	}
	fmt.Printf("Order confirmation email sent with code %d %s\n", res.StatusCode, res.Body)
	return nil
}

func (s *Sender) SendShippedOrderEmail() (err error) {

	return nil
}

func (s *Sender) SendOrderPickUpEmail() (err error) {

	return nil
}
func (s *Sender) SendShippingReturnReminder() (err error) {

	return nil
}

func (s *Sender) SendPickupReturnReminder() (err error) {

	return nil
}

func (s *Sender) SendNoticeOfLateFeeEmail() error {

	return nil
}

func (s *Sender) SendPickupReturnEmail() (err error) {
	//TODO implement me
	panic("implement me")
}
