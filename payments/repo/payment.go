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
	"strings"
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

func (repo *PaymentRepository) InsertCheckoutSessions(ctx context.Context, uow uow.UnitOfWork, checkoutSessionID string, status string, booking domain.Booking) (err error) {

	utcNow := time.Now().UTC()

	_, err = uow.GetTx().ExecContext(ctx, insertCheckoutSession,
		&checkoutSessionID,
		&status,
		&booking.BookingID,
		utcNow,
		utcNow,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PaymentRepository) UpdateCheckoutSessionStatus(ctx context.Context, uow uow.UnitOfWork, checkoutSessionID string) (bookingIds []string, err error) {

	utcNow := time.Now().UTC()

	rows, err := uow.GetTx().QueryContext(ctx, updateCheckoutSessions, utcNow, checkoutSessionID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	bookingIds = []string{}

	for rows.Next() {

		var bookingID string

		if err := rows.Scan(&bookingID); err != nil {
			return nil, err
		}

		bookingIds = append(bookingIds, bookingID)
	}

	return bookingIds, nil
}

func (repo *PaymentRepository) UpdateBookings(ctx context.Context, uow uow.UnitOfWork, bookingIds []string, shippingMethod string) (err error) {

	utcNow := time.Now().UTC()

	inClause := fmt.Sprintf("{%s}", strings.Join(bookingIds, ","))

	rows, err := uow.GetTx().QueryContext(ctx, updateBookings, utcNow, "Scheduled", inClause, shippingMethod)

	if err != nil {
		return nil
	}

	defer rows.Close()

	return nil
}

func (repo *PaymentRepository) GetCustomerById(ctx *gin.Context, customerID string) (customer domain.Customer, err error) {

	err = repo.db.QueryRowContext(ctx, findCustomerByID, customerID).Scan(
		&customer.Email,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return customer, ErrCustomerNotFound
	case err != nil:
		return customer, err
	}

	return customer, nil

}

func (repo *PaymentRepository) GetBookingsWithProductInfo(ctx *gin.Context, bookingIDs []string) (productBooking []models.ProductBooking, err error) {

	inClause := fmt.Sprintf("{%s}", strings.Join(bookingIDs, ","))

	rows, err := repo.db.QueryContext(ctx, findBookingsWithProductInfo, inClause)

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
