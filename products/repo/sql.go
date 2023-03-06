package repo

const (
	createProduct = `INSERT INTO products (product_id, title, description, quantity, price, created_at, updated_at) 
								VALUES($1, $2, $3, $4, $5, $6, $7)
								RETURNING *;`

	updateProduct = `UPDATE products set title = $2, description = $3, quantity = $4, price = $5, updated_at = $6
							WHERE product_id = $1 RETURNING *;`

	findProductByID = ` SELECT product_id, title, description, quantity, price, created_at, updated_at 
						FROM products WHERE product_id = $1;`

	deleteProduct = `DELETE FROM products WHERE product_id = $1;`

	findAll = `SELECT product_id, title, description, quantity, price, created_at, updated_at  FROM products;`
)
