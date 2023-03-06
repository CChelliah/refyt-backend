package repo

const (
	updateSeller = `UPDATE users set stripe_connect_id = $1, updated_at = $2
							WHERE uid = $3 RETURNING uid, email, stripe_id, stripe_connect_id, customer_number, created_at, updated_at;`

	findSellerById = ` SELECT uid, email, stripe_id, customer_number, created_at, updated_at FROM users WHERE uid = $1;`
)
