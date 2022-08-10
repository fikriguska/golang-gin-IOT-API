package service

import (
	"src/repository"
)

type userService struct {
	UserRepository repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) UserService {
	return &userService{
		UserRepository: *userRepository,
	}
}

func (u userService) Create(request repository.User) repository.User {
	return u.UserRepository.Create(repository.User{})
}
