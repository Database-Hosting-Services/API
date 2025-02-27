package accounts

import "time"

type User struct {
	ID        int       `db:"id" json:"id"`
	OID       string    `db:"oid" json:"oid"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password"`
	Image     string    `db:"image" json:"image"`
	Verified  bool      `db:"verified" json:"verified"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	LastLogin time.Time `db:"last_login" json:"last_login"`
}
