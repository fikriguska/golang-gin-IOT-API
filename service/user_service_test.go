package service_test

import (
	"testing"

	"src/model"
	"src/repository/mocks"
	"src/service"

	"github.com/stretchr/testify/mock"
)

func TestAddUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepo)
	mockUser := model.User{
		Username: "fikriguska",
		Password: "wkwk",
		Email:    "bintangf00code@gmail.com",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Add", mock.AnythingOfType("map[string]interface{}")).Return()
		mockUserRepo.On("GetIdByUsername", "fikriguska").Return(1)

		userService := service.NewUserService(mockUserRepo)
		userService.Add(mockUser)
	})
}
