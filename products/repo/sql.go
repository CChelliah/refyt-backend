package repo

const (
	createProduct = `INSERT INTO products (product_id, product_name, description, designer, category, fit_notes, size, rrp, price, shipping_price, user_uid, images, created_at, updated_at) 
								VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
								RETURNING product_id, product_name, description, designer, category, fit_notes, size, rrp, price, shipping_price, images, created_at, updated_at;`

	updateProduct = `UPDATE products set product_name = $2, description = $3, quantity = $4, price = $5, updated_at = $6
							WHERE product_id = $1 RETURNING product_id, product_name, description, quantity, price, rrp, designer, fit_notes, created_at, updated_at;`

	findProductByID = `SELECT product_id, product_name, description, designer, category, fit_notes, size, rrp, price, shipping_price, images, created_at, updated_at 
						FROM products WHERE product_id = $1;`

	deleteProduct = `DELETE FROM products WHERE product_id = $1;`

	find = `SELECT product_id, product_name, description, designer, category, fit_notes, size, rrp, price, shipping_price, images, created_at, updated_at FROM products WHERE true`

	findAllByUserID = `SELECT product_id, product_name, description, designer, category, fit_notes, size, rrp, price, shipping_price, images, created_at, updated_at FROM products WHERE user_uid = $1;`
)
