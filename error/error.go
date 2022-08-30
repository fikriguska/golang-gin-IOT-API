package error

import "errors"

// const (
// 	INVALID_PARAMS = 400
// 	INVALID_EMAIL  = 1001
// )

var (
	ErrInvalidParams           = errors.New("parameter is invalid")
	ErrInvalidEmail            = errors.New("email is invalid")
	ErrUserExist               = errors.New("email or Username already exists")
	ErrAddUserFail             = errors.New("add user failed")
	ErrInvalidToken            = errors.New("token is invalid")
	ErrActivateUserFail        = errors.New("activate user failed")
	ErrUsernameOrPassIncorrect = errors.New("username not found or password incorrect")
	ErrUserNotActive           = errors.New("account is inactive, check email for activation")
	// ErrUserLoginFailed         = errors.New("user is act")

	ErrInvalidHardwareType = errors.New("type must single-board computer, microcontroller unit, or sensor")
)

func PanicIfNeeded(err interface{}) {
	if err != nil {
		panic(err)
	}
}
