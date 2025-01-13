package entity

type User struct {
	UserID   string
	Email    string
	PassHash []byte
}
