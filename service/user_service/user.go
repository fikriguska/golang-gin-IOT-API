package user_service

import (
	"crypto/sha256"

	"encoding/hex"
	"src/models"
	"src/util"
)

type User struct {
	models.User
}

func (u *User) IsExist() bool {
	exist := false

	if u.Username != "" {
		exist = exist || models.IsUserUsernameExist(u.Username)
	}

	if u.Email != "" {
		exist = exist || models.IsUserEmailExist(u.Email)
	}

	if u.Id > 0 {
		exist = exist || models.IsUserExistById(u.Id)
	}

	return exist
}

func (u *User) Add() {
	hashedTokenByte := sha256.Sum256([]byte(u.Username + u.Email + u.Password))
	hashedToken := hex.EncodeToString(hashedTokenByte[:])
	hashedPassByte := sha256.Sum256([]byte(u.Password))
	hashedPass := hex.EncodeToString(hashedPassByte[:])

	u.Status = false
	u.Is_admin = false
	u.Password = hashedPass
	u.Token = hashedToken

	models.AddUser(u.User)
	id := models.GetUserIdByUsername(u.Username)
	sendEmailActivation(id, u.Username, u.Email, hashedToken)
}

func (u *User) Activate() {
	models.ActivateUser(u.Token)
}

func (u *User) TokenValidation() (bool, bool) {
	exist := models.IsUserTokenExist(u.Token)

	activated := models.IsUserActivatedCheckByToken(u.Token)

	return exist, activated
}

func (u *User) Auth() (bool, bool) {
	hashedPass := util.Sha256String(u.Password)
	credCorrect := models.IsUsernameAndPasswordExist(u.Username, hashedPass)
	activated := false
	if credCorrect {
		activated = models.IsUserActivatedCheckByUsername(u.Username)
	}
	return credCorrect, activated
}

func (u *User) Get() (int, string, string, bool) {
	res := models.GetUserByUsername(u.User)
	return res.Id, res.Username, res.Password, res.Is_admin
}

func (u *User) IsEmailAndUsernameMatched() (bool, bool) {
	match := models.IsEmailAndUsernameExist(u.Email, u.Username)
	activated := models.IsUserActivatedCheckByUsername(u.Username)
	return match, activated
}

func (u *User) SetRandomPassword() {
	newPass := util.RandomString(10)
	hashedNewPass := util.Sha256String(newPass)
	sendEmailForgetPassword(u.Email, u.Username, newPass)
	models.UpdateUserPasswordByEmail(u.Email, hashedNewPass)
}

func (u *User) IsUsingNode() bool {
	return models.IsNodeExistByUserId(u.Id)
}

func (u *User) Delete() {
	models.DeleteUser(u.Id)
}

func (u *User) SetPassword() {
	hashedPasswd := util.Sha256String(u.Password)
	models.UpdateUserPasswordById(u.Id, hashedPasswd)
}
