package entity

type User struct {
	UserID   string `db:"userID"`
	Email    string `db:"email"`
	PassHash []byte `db:"passHash"`
}
