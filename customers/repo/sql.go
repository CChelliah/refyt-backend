package repo

const (
	insertCustomer = `INSERT INTO customer_service.customers (customer_id, email, stripe_id, created_at, updated_at) 
								VALUES($1, $2, $3, $4, $5)
								RETURNING customer_id, email, stripe_id, customer_number, created_at, updated_at;`

	findCustomerByID = `SELECT customer_id, email, stripe_id, customer_number, stripe_connect_id, created_at, updated_at FROM customer_service.customers WHERE customer_id = $1`

	updateCustomer = `UPDATE customer_service.customers SET stripe_connect_id = $1, updated_at = $2 WHERE customer_id = $3
						RETURNING customer_id, email, stripe_id, customer_number, stripe_connect_id, created_at, updated_at;`
)
