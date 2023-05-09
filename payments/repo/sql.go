package repo

const (
	insertPayment = `INSERT INTO payment_service.payments (payment_id, checkout_session_id, status, booking_id, created_at, updated_at) 
								VALUES($1, $2, $3, $4, $5, $6);`

	updatePayment = `UPDATE payment_service.payments SET status = 'Paid', updated_at = $1 WHERE checkout_session_id = $2 RETURNING booking_id;`

	updateBooking = `UPDATE booking_service.bookings SET status = $2, updated_at = $1, shipping_method = $4 WHERE booking_id = $3;`

	findCustomerByID = `SELECT email FROM customer_service.customers WHERE customer_id = $1`

	findBookingsWithProductInfo = `SELECT b.booking_id, b.customer_id, b.start_date, b.end_date, b.status, b.shipping_method, p.product_id,
									p.product_name, p.description, p.price, p.rrp, p.fit_notes, p.designer, p.category, p.shipping_price, p.size, p.images
									FROM booking_service.bookings b LEFT JOIN product_service.products p ON b.product_id = p.product_id 
									WHERE b.booking_id = $1;`

	findBookingByID = `SELECT booking_id, b.product_id, b.customer_id, start_date, end_date, status, b.created_at, b.updated_at, shipping_price
							FROM booking_service.bookings b INNER JOIN product_service.products p ON 
    						b.product_id = p.product_id WHERE b.booking_id = $1;`
)
