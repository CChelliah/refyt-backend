package repo

const (
	findBookingsByProductID = `SELECT booking_id, product_id, customer_id, start_date, end_date, status, created_at, updated_at
							FROM bookings WHERE product_id = ANY($1::TEXT[])
							ORDER BY product_id, start_date DESC;`

	insertBooking = `INSERT INTO bookings (booking_id, product_id, customer_id, start_date, end_date, status, created_at, updated_at)
 						VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
)
