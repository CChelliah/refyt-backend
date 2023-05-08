package repo

const (
	findBookingsBySellerID = `SELECT booking_id, b.product_id, customer_id, start_date, end_date, status, created_at, updated_at, sp.product_name, sp.images FROM
								(SELECT product_id, product_name, images FROM product_service.products p WHERE user_uid = $1) as sp
								INNER JOIN booking_service.bookings b ON sp.product_id = b.product_id
								ORDER BY b.start_date DESC;`

	findBookingsByProductID = `SELECT booking_id, product_id, customer_id, start_date, end_date, status, created_at, updated_at FROM
								booking_service.bookings WHERE product_id = $1 and status = $2;`

	insertBooking = `INSERT INTO booking_service.bookings (booking_id, product_id, customer_id, start_date, end_date, status, created_at, updated_at)
 						VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
)
