package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"refyt-backend/bookings/domain"
	"refyt-backend/libs"
	"refyt-backend/libs/uow"
)

type IBookingRepo interface {
	FindBookingsBySellerID(ctx context.Context, sellerID string) (bookings []domain.Booking, err error)
	FindBookingByProductID(ctx context.Context, productID string) (bookings []domain.Booking, err error)
	InsertBooking(ctx context.Context, uow uow.UnitOfWork, booking domain.Booking) (err error)
}

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
			pq.Array(&booking.ImageUrls),
		)

		fmt.Printf("%s\n", booking.BookingID)

		if err != nil {
			return bookings, err
		}

		bookings = append(bookings, booking)
	}

	fmt.Printf("%d\n", len(bookings))

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

func (repo *BookingRepo) InsertBooking(ctx context.Context, uow uow.UnitOfWork, booking domain.Booking) (err error) {

	_, err = uow.GetTx().ExecContext(ctx, insertBooking,
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
		return err
	}

	return nil
}
