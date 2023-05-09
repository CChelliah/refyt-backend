package repo

const (
	createProduct = `INSERT INTO product_service.products (product_id, product_name, description, designer, category, fit_notes, size, rrp, price, shipping_price, customer_id, images, created_at, updated_at) 
								VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
								RETURNING product_id, product_name, description, designer, category, fit_notes, size, rrp, price, shipping_price, images, created_at, updated_at;`

	updateProduct = `UPDATE product_service.products p set product_name = $2, description = $3, designer = $4, category = $5, fit_notes = $6, size = $7, rrp = $8, price = $9, shipping_price = $10, images = $11, updated_at = $12
							WHERE p.product_id = $1 RETURNING product_id, product_name, description, designer, category, fit_notes, size, rrp, price, shipping_price, images, updated_at, created_at;`

	findProductByID = `SELECT product_id, product_name, description, designer, category, fit_notes, size, rrp, price, shipping_price, images, created_at, updated_at 
						FROM product_service.products WHERE product_id = $1;`

	deleteProduct = `DELETE FROM product_service.products WHERE product_id = $1;`

	find = `SELECT product_id, product_name, description, designer, category, fit_notes, size, rrp, price, shipping_price, images, created_at, updated_at FROM product_service.products WHERE true`

	findAllByUserID = `SELECT product_id, product_name, description, designer, category, fit_notes, size, rrp, price, shipping_price, images, created_at, updated_at FROM product_service.products WHERE customer_id = $1;`
)
