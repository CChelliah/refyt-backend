package repo

const (
	findBookingsByProductID = `SELECT booking_id, product_id, customer_id, start_date, end_date, status, created_at, updated_at
							FROM bookings WHERE product_id = ANY($1::TEXT[])
							ORDER BY product_id, start_date DESC;`

	insertBooking = `INSERT INTO bookings (booking_id, product_id, customer_id, start_date, end_date, status, created_at, updated_at)
 						VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	findBookingsBySellerID = `SELECT booking_id, b.product_id, customer_id, start_date, end_date, status, created_at, updated_at, sp.product_name FROM
								(SELECT product_id, product_name FROM products p WHERE user_uid = $1) as sp
								INNER JOIN bookings b ON sp.product_id = b.product_id
								ORDER BY b.start_date DESC;`

	insertCheckoutSession = `INSERT INTO checkout_sessions (checkout_session_id, status, booking_id, created_at, updated_at) 
								VALUES($1, $2, $3, $4, $5);`

	updateCheckoutSessions = `UPDATE checkout_sessions SET status = 'paid', updated_at = $1 WHERE checkout_session_id = $2 RETURNING booking_id;`

	updateBookings = `UPDATE bookings SET status = $2, updated_at = $1, shipping_method = $4 WHERE booking_id = ANY($3::TEXT[]);`
)
