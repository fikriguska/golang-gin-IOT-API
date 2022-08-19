package service

import (
	e "src/error"

	"crypto/sha256"
	"net/smtp"

	"fmt"
	"strconv"

	"encoding/hex"
	"src/model"
	"src/repository"
)

type userServiceImpl struct {
	UserRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) UserService {
	return &userServiceImpl{
		UserRepo: userRepo,
	}
}

func (service *userServiceImpl) IsExist(request model.User) bool {
	return service.UserRepo.IsUsernameExist(request.Username)
}

func (service *userServiceImpl) Add(request model.User) error {

	exist := service.IsExist(request)

	if exist {
		return e.ErrUserExist
	}

	hashedTokenByte := sha256.Sum256([]byte(request.Username + request.Email + request.Password))
	hashedToken := hex.EncodeToString(hashedTokenByte[:])
	hashedPassByte := sha256.Sum256([]byte(request.Password))
	hashedPass := hex.EncodeToString(hashedPassByte[:])
	user := map[string]interface{}{
		"email":    request.Email,
		"username": request.Username,
		"password": hashedPass,
		"status":   false,
		"token":    hashedToken,
	}
	service.UserRepo.Add(user)
	id := service.UserRepo.GetIdByUsername(request.Username)
	sendEmail(id, request.Username, request.Email, hashedToken)
	return nil
}

func (service *userServiceImpl) Activate(request model.User) {
	service.UserRepo.Activate(request.Token)
}

func (service *userServiceImpl) IsTokenValid(request model.User) bool {
	exist := service.UserRepo.IsTokenExist(request.Token)

	activated := service.UserRepo.IsActivatedCheckByToken(request.Token)

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

func (service *userServiceImpl) Auth(request model.User) (bool, bool) {
	hashedPassByte := sha256.Sum256([]byte(request.Password))
	hashedPass := hex.EncodeToString(hashedPassByte[:])
	credCorrect := service.UserRepo.Auth(request.Username, hashedPass)

	activated := service.UserRepo.IsActivatedCheckByUsername(request.Username)
	fmt.Println(activated)
	return credCorrect, activated

}
