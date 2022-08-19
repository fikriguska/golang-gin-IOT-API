package service

type UserService interface {
	IsExist(User) bool
	Add(User)
	Activate(User)
	IsTokenValid(User) bool
	Auth(User) (bool, bool)
}
