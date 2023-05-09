package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"refyt-backend/libs"
	"refyt-backend/libs/email/sendgrid/models"
	"refyt-backend/libs/uow"
	"refyt-backend/payments/domain"
	"time"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrBookingNotFound  = errors.New("booking not found")
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(env *libs.PostgresDatabase) (billingRepo PaymentRepository) {

	billingRepo = PaymentRepository{
		db: env.Db,
	}

	return billingRepo
}

func (repo *PaymentRepository) InsertPayment(ctx context.Context, uow uow.UnitOfWork, payment domain.Payment) (err error) {

	_, err = uow.GetTx().ExecContext(ctx, insertPayment,
		&payment.PaymentID,
		&payment.CheckoutSessionID,
		&payment.Status,
		&payment.BookingID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PaymentRepository) UpdatePaymentStatus(ctx context.Context, uow uow.UnitOfWork, checkoutSessionID string) (bookingID string, err error) {

	fmt.Println("Payment Status 1")
	utcNow := time.Now().UTC()

	err = uow.GetTx().QueryRowContext(ctx, updatePayment, utcNow, checkoutSessionID).Scan(&bookingID)

	if err != nil {
		return bookingID, err
	}

	fmt.Println("Payment Status 2")

	return bookingID, nil
}

func (repo *PaymentRepository) UpdateBooking(ctx context.Context, uow uow.UnitOfWork, bookingID string, shippingMethod string) (err error) {

	fmt.Println("Update bookings 1")

	utcNow := time.Now().UTC()

	fmt.Printf("InClause %s\n", shippingMethod)

	_, err = uow.GetTx().ExecContext(ctx, updateBooking, utcNow, "Scheduled", bookingID, shippingMethod)

	if err != nil {
		return nil
	}

	fmt.Println("Update booking 2")
	return nil
}

func (repo *PaymentRepository) GetCustomerById(ctx *gin.Context, customerID string) (customer domain.Customer, err error) {

	fmt.Println("Customer ID 1")

	err = repo.db.QueryRowContext(ctx, findCustomerByID, customerID).Scan(
		&customer.Email,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return customer, ErrCustomerNotFound
	case err != nil:
		return customer, err
	}

	fmt.Println("Customer ID 2")

	return customer, nil

}

func (repo *PaymentRepository) GetBookingsWithProductInfo(ctx *gin.Context, bookingID string) (productBooking []models.ProductBooking, err error) {

	rows, err := repo.db.QueryContext(ctx, findBookingsWithProductInfo, bookingID)

	if err != nil {
		return productBooking, err
	}

	defer rows.Close()

	productBookings := []models.ProductBooking{}

	for rows.Next() {

		var booking models.ProductBooking
		var product models.Product

		err = rows.Scan(
			&booking.BookingID,
			&booking.CustomerID,
			&booking.StartDate,
			&booking.EndDate,
			&booking.Status,
			&booking.ShippingMethod,
			&product.ProductID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.RRP,
			&product.FitNotes,
			&product.Designer,
			&product.Category,
			&product.ShippingPrice,
			&product.Size,
			pq.Array(&product.ImageUrls),
		)

		if err != nil {
			return productBooking, err
		}

		productBookings = append(productBookings, booking)
	}

	if err = rows.Err(); err != nil {
		return productBookings, err
	}

	return productBookings, nil
}

func (repo *PaymentRepository) FindBookingByID(ctx context.Context, bookingID string) (booking domain.Booking, err error) {

	err = repo.db.QueryRowContext(ctx, findBookingByID, bookingID).Scan(
		&booking.BookingID,
		&booking.ProductID,
		&booking.CustomerID,
		&booking.StartDate,
		&booking.EndDate,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
		&booking.ShippingPrice,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return booking, ErrBookingNotFound
	case err != nil:
		return booking, err
	}

	return booking, nil

}
