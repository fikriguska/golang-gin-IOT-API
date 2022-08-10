package service

import "src/repository"

type UserService interface {
	Create(request repository.User) repository.User
}
