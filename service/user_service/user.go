package user_service

import (
	"crypto/sha256"
	"net/mail"
	"net/smtp"

	"fmt"
	"strconv"

	"encoding/hex"
	"src/models"
)

type User struct {
	models.User
}

func (u *User) IsEmailValid() bool {
	_, err := mail.ParseAddress(u.Email)
	return err == nil
}

func (u *User) IsExist() bool {
	return models.IsUserUsernameExist(u.Username) || models.IsUserEmailExist(u.Email)
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
	sendEmail(id, u.Username, u.Email, hashedToken)
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

func sendEmail(id int, username string, to string, token string) error {

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
