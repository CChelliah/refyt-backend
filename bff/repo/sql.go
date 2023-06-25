package repo

const (
	find = `SELECT product_id, product_name, description, designer, category, fit_notes, size, rrp, price, shipping_price, images, p.created_at, p.updated_at, c.customer_id, c.user_name FROM product_service.products p LEFT JOIN customer_service.customers c ON c.customer_id = p.customer_id WHERE true`

	findAllByUserID = `SELECT product_id, product_name, description, designer, category, fit_notes, size, rrp, price, shipping_price, images, p.created_at, p.updated_at,  c.customer_id, c.user_name FROM product_service.products p LEFT JOIN customer_service.customers c ON c.customer_id = p.customer_id WHERE customer_id = $1;`

	findProductByID = `SELECT product_id, product_name, description, designer, category, fit_notes, size, rrp, price, shipping_price, images, p.created_at, p.updated_at,  c.customer_id, c.user_name
						FROM product_service.products p LEFT JOIN customer_service.customers c ON c.customer_id = p.customer_id WHERE product_id = $1;`
)
