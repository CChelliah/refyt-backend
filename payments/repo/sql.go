package repo

const (
	insertCheckoutSession = `INSERT INTO checkout_sessions (checkout_session_id, status, booking_id, created_at, updated_at) 
								VALUES($1, $2, $3, $4, $5);`

	updateCheckoutSessions = `UPDATE checkout_sessions SET status = 'paid', updated_at = $1 WHERE checkout_session_id = $2 RETURNING booking_id;`

	updateBookings = `UPDATE bookings SET status = $2, updated_at = $1, shipping_method = $4 WHERE booking_id = ANY($3::TEXT[]);`

	findCustomerByID = `SELECT email FROM users WHERE uid = $1`

	findBookingsWithProductInfo = `SELECT b.booking_id, b.customer_id, b.start_date, b.end_date, b.status, b.shipping_method, p.product_id,
									p.product_name, p.description, p.price, p.rrp, p.fit_notes, p.designer, p.category, p.shipping_price, p.size, p.images
									FROM bookings b LEFT JOIN products p ON b.product_id = p.product_id 
									WHERE b.booking_id = ANY($1::TEXT[]);`

	findBookingByID = `SELECT booking_id, product_id, customer_id, start_date, end_date, status, created_at, updated_at, shipping_price
							FROM booking_service.bookings b INNER JOIN product_service.products p ON 
    						b.product_id = p.product_id WHERE b.booking_id = $1;`
)
