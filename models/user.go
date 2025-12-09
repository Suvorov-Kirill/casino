package models

type User struct {
	ID       int
	Email    string
	Password string
	Balance  int
	Role     string
}
