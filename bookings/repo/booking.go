package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"refyt-backend/bookings/domain"
	"refyt-backend/libs"
	"refyt-backend/libs/email/sendgrid/models"
	"refyt-backend/libs/uow"
)

type IBookingRepo interface {
	FindBookingsBySellerID(ctx context.Context, sellerID string) (bookings []domain.Booking, err error)
	FindBookingByProductID(ctx context.Context, productID string) (bookings []domain.Booking, err error)
	Store(ctx context.Context, uow uow.UnitOfWork, booking domain.Booking) (err error)
	FindBookingByID(ctx context.Context, uow uow.UnitOfWork, bookingID string) (booking domain.Booking, err error)
}

var (
	ErrCustomerNotFound = errors.New("payment not found")
)

type BookingRepo struct {
	db *sql.DB
}

func NewBookingRepo(env *libs.PostgresDatabase) (bookingRepo BookingRepo) {

	bookingRepo = BookingRepo{
		db: env.Db,
	}

	return bookingRepo
}

func (repo *BookingRepo) FindBookingsBySellerID(ctx context.Context, sellerID string) (bookings []domain.Booking, err error) {

	rows, err := repo.db.QueryContext(ctx, findBookingsBySellerID, sellerID)

	if err != nil {
		return bookings, err
	}

	defer rows.Close()

	bookings = []domain.Booking{}

	for rows.Next() {

		var booking domain.Booking

		err = rows.Scan(
			&booking.BookingID,
			&booking.ProductID,
			&booking.CustomerID,
			&booking.StartDate,
			&booking.EndDate,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
			&booking.ProductName,
			&booking.ShippingMethod,
			pq.Array(&booking.ImageUrls),
		)

		if err != nil {
			return bookings, err
		}

		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return bookings, err
	}

	return bookings, nil
}

func (repo *BookingRepo) FindBookingByProductID(ctx context.Context, productID string) (bookings []domain.Booking, err error) {

	rows, err := repo.db.QueryContext(ctx, findBookingsByProductID, productID, "Scheduled")

	if err != nil {
		return bookings, err
	}

	defer rows.Close()

	bookings = []domain.Booking{}

	for rows.Next() {

		var booking domain.Booking

		err = rows.Scan(
			&booking.BookingID,
			&booking.ProductID,
			&booking.CustomerID,
			&booking.StartDate,
			&booking.EndDate,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)

		if err != nil {
			return bookings, err
		}

		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return bookings, err
	}

	return bookings, nil

}

func (repo *BookingRepo) Store(ctx context.Context, uow uow.UnitOfWork, booking domain.Booking) (err error) {

	_, err = uow.GetTx().ExecContext(ctx, insertBooking,
		&booking.BookingID,
		&booking.ProductID,
		&booking.CustomerID,
		&booking.StartDate,
		&booking.EndDate,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
		&booking.ShippingMethod,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repo *BookingRepo) FindBookingByID(ctx context.Context, uow uow.UnitOfWork, bookingID string) (booking domain.Booking, err error) {

	err = uow.GetTx().QueryRowContext(ctx, findBookingByID, bookingID).Scan(
		&booking.BookingID,
		&booking.ProductID,
		&booking.CustomerID,
		&booking.StartDate,
		&booking.EndDate,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
		&booking.ShippingMethod,
	)

	if err != nil {
		return booking, err
	}

	return booking, nil
}

func (repo *BookingRepo) GetCustomerById(ctx context.Context, customerID string) (customer domain.Customer, err error) {

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

func (repo *BookingRepo) FindOrderConfirmationPayload(ctx context.Context, bookingID string) (booking models.ProductBooking, err error) {

	err = repo.db.QueryRowContext(ctx, findBookingsWithProductInfo, bookingID).Scan(
		&booking.BookingID,
		&booking.CustomerID,
		&booking.StartDate,
		&booking.EndDate,
		&booking.Status,
		&booking.ShippingMethod,
		&booking.Product.ProductID,
		&booking.Product.Name,
		&booking.Product.Description,
		&booking.Product.Price,
		&booking.Product.RRP,
		&booking.Product.FitNotes,
		&booking.Product.Designer,
		&booking.Product.Category,
		&booking.Product.ShippingPrice,
		&booking.Product.Size,
		pq.Array(&booking.Product.ImageUrls))

	if err != nil {
		return booking, err
	}

	return booking, nil
}
