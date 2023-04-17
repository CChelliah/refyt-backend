package repo

const (
	insertUser = `INSERT INTO users (uid, email, stripe_id, created_at, updated_at) 
								VALUES($1, $2, $3, $4, $5)
								RETURNING uid, email, stripe_id, customer_number, created_at, updated_at;`

	findUserByID = `SELECT uid, email, stripe_id, customer_number, stripe_connect_id, created_at, updated_at FROM users WHERE uid = $1`
)
