package repo

import (
	"context"
	"database/sql"
	"fmt"
	"refyt-backend/billing/domain"
	"refyt-backend/common"
	"refyt-backend/common/uow"
	"strings"
)

type IBillingRepository interface {
	GetBookingsByProductID(ctx context.Context, uow uow.UnitOfWork, productIDs []string) (bookings map[string][]domain.Booking, err error)
	InsertBookings(ctx context.Context, uow uow.UnitOfWork, bookings []domain.Booking) (err error)
}

type BillingRepository struct {
	db *sql.DB
}

func NewBillingRepository(env *common.Env) (billingRepo BillingRepository) {

	billingRepo = BillingRepository{
		db: env.Db,
	}

	return billingRepo
}

func (repo *BillingRepository) GetExistingBookingsByProductID(ctx context.Context, uow uow.UnitOfWork, productIDs []string) (bookings map[string][]domain.Booking, err error) {

	inClause := fmt.Sprintf("{%s}", strings.Join(productIDs, ","))

	bookings = map[string][]domain.Booking{}

	rows, err := repo.db.QueryContext(ctx, findBookingsByProductID, inClause)

	if err != nil {
		return bookings, err
	}

	defer rows.Close()

	var productID string
	var booking domain.Booking

	for rows.Next() {

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

		if booking.BookingID != productID {
			productID = booking.ProductID
		}

		bookings[productID] = append(bookings[productID], booking)
	}

	bookings[productID] = append(bookings[productID], booking)

	if err = rows.Err(); err != nil {
		return bookings, err
	}

	return bookings, nil
}

func (repo *BillingRepository) InsertBookings(ctx context.Context, uow uow.UnitOfWork, bookings []domain.Booking) (err error) {

	fmt.Println("Insert Bookings")

	for _, booking := range bookings {

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
	}
	return nil
}
