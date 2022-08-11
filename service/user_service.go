package service

import "src/repository"

type UserService interface {
	Create(user repository.User) error
	IsExist(user repository.User) (bool, error)
}
