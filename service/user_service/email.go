package user_service

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"strconv"
)

func (u *User) IsEmailValid() bool {
	_, err := mail.ParseAddress(u.Email)
	return err == nil
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
