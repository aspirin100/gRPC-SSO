package entity

type App struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}
