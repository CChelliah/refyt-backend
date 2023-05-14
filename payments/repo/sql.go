package repo

const (
	insertPayment = `INSERT INTO payment_service.payments (payment_id, checkout_session_id, status, booking_id, created_at, updated_at) 
								VALUES($1, $2, $3, $4, $5, $6)
								ON CONFLICT (payment_id) DO UPDATE
									SET status = $3,
								    	created_at = $5;`

	updatePayment = `UPDATE payment_service.payments SET status = 'Paid', updated_at = $1 WHERE checkout_session_id = $2 RETURNING booking_id;`

	updateBooking = `UPDATE booking_service.bookings SET status = $2, updated_at = $1, shipping_method = $4 WHERE booking_id = $3;`

	findBookingByID = `SELECT booking_id, b.product_id, b.customer_id, start_date, end_date, status, b.created_at, b.updated_at, shipping_price
							FROM booking_service.bookings b INNER JOIN product_service.products p ON 
    						b.product_id = p.product_id WHERE b.booking_id = $1;`

	findPaymentByCheckoutSessionID = `SELECT payment_id, checkout_session_id, status, booking_id, created_at, updated_at 
							FROM payment_service.payments WHERE checkout_session_id = $1`
)
