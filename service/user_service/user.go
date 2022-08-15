package user_service

import (
	"crypto/sha256"

	"encoding/hex"
	"src/models"
)

type User struct {
	Id       int
	Email    string
	Username string
	Password string
	Status   bool
	Token    string
}

func (u *User) IsExist() (bool, error) {
	check, err := models.IsUserUsernameExist(u.Username)
	if err != nil {
		return false, err
	}
	if check {
		return true, nil
	}
	return false, nil
}

func (u *User) Add() error {
	hashedTokenByte := sha256.Sum256([]byte(u.Username + u.Email + u.Password))
	hashedToken := hex.EncodeToString(hashedTokenByte[:])
	hashedPassByte := sha256.Sum256([]byte(u.Password))
	hashedPass := hex.EncodeToString(hashedPassByte[:])
	u.Status = false
	user := map[string]interface{}{
		"email":    u.Email,
		"username": u.Username,
		"password": hashedPass,
		"status":   u.Status,
		"token":    hashedToken,
	}
	err := models.AddUser(user)
	return err
}
