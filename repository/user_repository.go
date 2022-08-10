package repository

type User struct {
	User_id  int
	Username string
}

type UserRepository interface {
	Create(user User) User
}
