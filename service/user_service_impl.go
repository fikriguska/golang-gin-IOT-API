package service

import (
	"errors"
	"fmt"
	e "src/error"

	"src/repository"
)

type userService struct {
	UserRepository repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) UserService {
	// var user repository.User
	// user.Email = "bintangf00code@gmail.com"
	// user.Username = "fikriguska"

	// user := repository.User{
	// 	Email:    "bintangf00code@gmail.com",
	// 	Username: "fikriguskax",
	// }
	// fmt.Println(user.)
	return &userService{
		UserRepository: *userRepository,
	}
}

func (u userService) Create(user repository.User) (err error) {
	fmt.Println(user)

	u.UserRepository.Create(repository.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Token:    "sssss",
	})
	return errors.Unwrap(fmt.Errorf("sss"))
}

func (u userService) IsExist(user repository.User) (bool, error) {
	if false {
		return true, e.ErrUserExist
	}
	fmt.Println("xxxx")
	fmt.Println(u.UserRepository.GetByUsername(user.Username))

	return false, nil

}
