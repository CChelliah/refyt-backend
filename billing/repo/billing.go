package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"refyt-backend/billing/domain"
	"refyt-backend/billing/model"
	"refyt-backend/libs"
	"refyt-backend/libs/uow"
	"strings"
	"time"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
)

type IBillingRepository interface {
	GetBookingsByProductID(ctx context.Context, uow uow.UnitOfWork, productIDs []string) (bookings map[string][]domain.Booking, err error)
	InsertBookings(ctx context.Context, uow uow.UnitOfWork, bookings []domain.Booking) (err error)
	FindBookingsBySellerID(ctx context.Context, sellerID string) (bookings []domain.Booking, err error)
	InsertCheckoutSessions(ctx context.Context, work uow.UnitOfWork, checkoutSessionID string, status string, bookings []domain.Booking) (err error)
	UpdateCheckoutSessionStatus(ctx context.Context, uow uow.UnitOfWork, checkoutSessionID string) (bookingIds []string, err error)
	UpdatePaidBookings(ctx context.Context, uow uow.UnitOfWork, bookingIds []string, shippingMethod string)
	GetCustomerById(ctx *gin.Context, customerID string) (customer domain.Customer, err error)
}

type BillingRepository struct {
	db *sql.DB
}

func NewBillingRepository(env *libs.PostgresDatabase) (billingRepo BillingRepository) {

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

func (repo *BillingRepository) FindBookingsBySellerID(ctx context.Context, sellerID string) (bookings []domain.Booking, err error) {

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

func (repo *BillingRepository) InsertCheckoutSessions(ctx context.Context, uow uow.UnitOfWork, checkoutSessionID string, status string, bookings []domain.Booking) (err error) {

	utcNow := time.Now().UTC()

	for _, booking := range bookings {

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
	}
	return nil

	return nil
}

func (repo *BillingRepository) UpdateCheckoutSessionStatus(ctx context.Context, uow uow.UnitOfWork, checkoutSessionID string) (bookingIds []string, err error) {

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

	return bookingIds, nil
}

func (repo *BillingRepository) UpdateBookings(ctx context.Context, uow uow.UnitOfWork, bookingIds []string, shippingMethod string) (err error) {

	utcNow := time.Now().UTC()

	inClause := fmt.Sprintf("{%s}", strings.Join(bookingIds, ","))

	rows, err := uow.GetTx().QueryContext(ctx, updateBookings, utcNow, "Scheduled", inClause, shippingMethod)

	if err != nil {
		return nil
	}

	defer rows.Close()

	return nil
}

func (repo *BillingRepository) GetCustomerById(ctx *gin.Context, customerID string) (customer domain.Customer, err error) {

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

func (repo *BillingRepository) GetBookingsWithProductInfo(ctx *gin.Context, bookingIDs []string) (productBooking []model.ProductBooking, err error) {

	inClause := fmt.Sprintf("{%s}", strings.Join(bookingIDs, ","))

	rows, err := repo.db.QueryContext(ctx, findBookingsWithProductInfo, inClause)

	if err != nil {
		return productBooking, err
	}

	defer rows.Close()

	productBookings := []model.ProductBooking{}

	for rows.Next() {

		var booking model.ProductBooking
		var product model.Product

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
			&product.Designer,
			&product.Category,
			&product.FitNotes,
			&product.Size,
			&product.RRP,
			&product.Price,
			&product.ShippingPrice,
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
