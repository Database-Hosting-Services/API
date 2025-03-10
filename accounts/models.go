package accounts

import "time"

type User struct {
	ID        int       `db:"id" json:"id"`
	OID       string    `db:"oid" json:"oid"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password"`
	Image     string    `db:"image" json:"image"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	LastLogin time.Time `db:"last_login" json:"last_login"`
}

type UserSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUnVerified struct {
	User
	Code string `josn:"code"`
}

type UpdatePasswordModel struct {
	CurrentPassword string `json:"current_password"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type ResetPasswordForm struct {
	Code     string `json:"code"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (u *User) GetOId() string {
	return u.OID
}

func (u *User) GetUserName() string {
	return u.Username
}
