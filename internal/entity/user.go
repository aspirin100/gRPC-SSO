package entity

type User struct {
	UserID   string `db:"id"`
	Email    string `db:"email"`
	PassHash []byte `db:"passHash"`
}
