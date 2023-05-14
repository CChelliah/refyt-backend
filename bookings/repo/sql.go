package repo

const (
	findBookingsBySellerID = `SELECT booking_id, b.product_id, customer_id, start_date, end_date, status, created_at, updated_at, sp.product_name, shipping_method, images FROM
								(SELECT product_id, product_name, images FROM product_service.products p WHERE customer_id = $1) as sp
								INNER JOIN booking_service.bookings b ON sp.product_id = b.product_id
								ORDER BY b.start_date DESC;`

	findBookingsByProductID = `SELECT booking_id, product_id, customer_id, start_date, end_date, status, created_at, updated_at FROM
								booking_service.bookings WHERE product_id = $1 and status = $2;`

	insertBooking = `INSERT INTO booking_service.bookings (booking_id, product_id, customer_id, start_date, end_date, status, created_at, updated_at, shipping_method)
 						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
 						ON CONFLICT (booking_id) DO UPDATE
							SET status = $6,
								updated_at = $8,
								shipping_method = $9;`

	findBookingByID = `SELECT booking_id, product_id, customer_id, start_date, end_date, status, created_at, updated_at, shipping_method  FROM
								booking_service.bookings WHERE booking_id = $1;`

	findCustomerByID = `SELECT email FROM customer_service.customers WHERE customer_id = $1`

	findBookingsWithProductInfo = `SELECT b.booking_id, b.customer_id, b.start_date, b.end_date, b.status, b.shipping_method, p.product_id,
									p.product_name, p.description, p.price, p.rrp, p.fit_notes, p.designer, p.category, p.shipping_price, p.size, p.images
									FROM booking_service.bookings b LEFT JOIN product_service.products p ON b.product_id = p.product_id 
									WHERE b.booking_id = $1;`
)
