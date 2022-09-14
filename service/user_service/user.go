package user_service

import (
	"crypto/sha256"
	"log"
	"net/mail"
	"net/smtp"

	"fmt"
	"strconv"

	"encoding/hex"
	"src/models"
	"src/util"
)

type User struct {
	models.User
}

func (u *User) IsEmailValid() bool {
	_, err := mail.ParseAddress(u.Email)
	return err == nil
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

func (u *User) IsTokenValid() bool {
	exist := models.IsUserTokenExist(u.Token)

	activated := models.IsUserActivatedCheckByToken(u.Token)

	valid := exist && !activated
	return valid
}

func genMessageWithheader(from string, to string, subject string, body string) string {
	return "Content-Type: text/html\n" +
		"From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		body
}

func sendEmailActivation(id int, username string, to string, token string) error {

	urlCode := "http://localhost:8080/user/activation?token=" + token

	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpAddr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	body := `
      <h1>Activation Message</h1>
      <h4>Dear ` + username + `</h4>
      <p>We have accepted your registration. Your account is:</p>
      <li>
        <ul> Id User: ` + strconv.Itoa(id) + `</ul>
        <ul> Username: ` + username + `</ul>
      </li>
      <p>Click <a href=` + urlCode + `>here</a> to activate your account</p>
      <p><h5>Thank you</h5></p>     
	`
	from := "skripsibintang@gmail.com"
	password := "vytfhtsgzxsnpbao"
	subject := "Activation Message"
	auth := smtp.PlainAuth("", from, password, smtpHost)

	message := genMessageWithheader(from, to, subject, body)

	err := smtp.SendMail(smtpAddr, auth, from, []string{to}, []byte(message))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Auth() (bool, bool) {
	hashedPassByte := sha256.Sum256([]byte(u.Password))
	hashedPass := hex.EncodeToString(hashedPassByte[:])
	credCorrect := models.AuthUser(u.Username, hashedPass)
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
	log.Println(newPass, hashedNewPass)
	sendEmailForgetPassword(u.Email, u.Username, newPass)
	models.UpdateUserPassword(u.Email, hashedNewPass)

}

func (u *User) IsUsingNode() bool {
	return models.IsNodeExistByUserId(u.Id)
}

func (u *User) Delete() {
	models.DeleteUser(u.Id)
}

func sendEmailForgetPassword(email string, username string, newPass string) error {

	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpAddr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	body := `
	<h3>Dear ` + username + `. </h3>

	<p>We have accepted your forget password request. Use this password for log in.</p>

	<p><h4>` + newPass + `</h4></p>

	<p>Thank You</p>  
	`
	from := "skripsibintang@gmail.com"
	password := "vytfhtsgzxsnpbao"

	subject := "New Password"
	auth := smtp.PlainAuth("", from, password, smtpHost)

	message := genMessageWithheader(from, email, subject, body)

	err := smtp.SendMail(smtpAddr, auth, from, []string{email}, []byte(message))
	if err != nil {
		return err
	}
	return nil
}
