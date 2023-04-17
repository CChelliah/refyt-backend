package repo

const (
	findBookingsBySellerID = `SELECT booking_id, b.product_id, customer_id, start_date, end_date, status, created_at, updated_at, sp.product_name FROM
								(SELECT product_id, product_name FROM products p WHERE user_uid = $1) as sp
								INNER JOIN bookings b ON sp.product_id = b.product_id
								ORDER BY b.start_date DESC;`
)
