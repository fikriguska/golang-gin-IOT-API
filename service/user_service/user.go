package user_service

import (
	"crypto/sha256"
	"net/smtp"

	"fmt"
	"strconv"

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
	if err != nil {
		return err
	}
	id, err := models.GetUserIdByUsername(u.Username)
	if err != nil {

		return err
	}
	err = sendEmail(id, u.Username, u.Email, hashedToken)
	return err
}

func genMessageWithheader(from string, to string, subject string, body string) string {
	return "Content-Type: text/html\n" +
		"From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		body
}

func sendEmail(id int, username string, to string, token string) error {

	urlCode := "http://192.168.100.105:5000/user/activation?token=" + token

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
