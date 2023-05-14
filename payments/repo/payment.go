package repo

import (
	"context"
	"database/sql"
	"errors"
	"refyt-backend/libs"
	"refyt-backend/libs/uow"
	"refyt-backend/payments/domain"
	"time"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrBookingNotFound = errors.New("booking not found")
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

func (repo *PaymentRepository) Store(ctx context.Context, uow uow.UnitOfWork, payment domain.Payment) (err error) {

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

	utcNow := time.Now().UTC()

	err = uow.GetTx().QueryRowContext(ctx, updatePayment, utcNow, checkoutSessionID).Scan(&bookingID)

	if err != nil {
		return bookingID, err
	}

	return bookingID, nil
}

func (repo *PaymentRepository) UpdateBooking(ctx context.Context, uow uow.UnitOfWork, bookingID string, shippingMethod string) (err error) {

	utcNow := time.Now().UTC()

	_, err = uow.GetTx().ExecContext(ctx, updateBooking, utcNow, "Scheduled", bookingID, shippingMethod)

	if err != nil {
		return nil
	}

	return nil
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

func (repo *PaymentRepository) FindPaymentByCheckoutSessionID(ctx context.Context, uow uow.UnitOfWork, paymentID string) (payment domain.Payment, err error) {

	err = uow.GetTx().QueryRowContext(ctx, findPaymentByCheckoutSessionID, paymentID).Scan(
		&payment.PaymentID,
		&payment.CheckoutSessionID,
		&payment.Status,
		&payment.BookingID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return payment, ErrPaymentNotFound
	case err != nil:
		return payment, err
	}

	return payment, nil

}
