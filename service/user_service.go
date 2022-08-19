package service

import "src/model"

type UserService interface {
	IsExist(model.User) bool
	Add(model.User) error
	Activate(model.User)
	IsTokenValid(model.User) bool
	Auth(model.User) (bool, bool)
}
