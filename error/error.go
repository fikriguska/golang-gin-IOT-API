package error

import "errors"

// const (
// 	INVALID_PARAMS = 400
// 	INVALID_EMAIL  = 1001
// )

var (
	ErrInvalidParams            = errors.New("parameter is invalid")
	ErrInvalidEmail             = errors.New("email is invalid")
	ErrUserExist                = errors.New("email or Username already exists")
	ErrAddUserFail              = errors.New("add user failed")
	ErrInvalidToken             = errors.New("token is invalid")
	ErrActivateUserFail         = errors.New("activate user failed")
	ErrUsernameOrPassIncorrect  = errors.New("username not found or password incorrect")
	ErrUsernameOrEmailIncorrect = errors.New("username or email incorrect")
	ErrUserNotActive            = errors.New("account is inactive, check email for activation")
	// ErrUserLoginFailed         = errors.New("user is act")

	ErrInvalidHardwareType     = errors.New("type must single-board computer, microcontroller unit, or sensor")
	ErrHardwareNotFound        = errors.New("hardware not found")
	ErrHardwareMustbeSensor    = errors.New("hardware type not match, type should be sensor")
	ErrUseHardwareNotPermitted = errors.New("you can't use another user's node")

	ErrNodeNotFound           = errors.New("node not found")
	ErrDeleteNodeNotPermitted = errors.New("you can't delete another user's node")
	ErrUseNodeNotPermitted    = errors.New("you can't use another user's node")
	ErrSeeNodeNotPermitted    = errors.New("you can't see another user's node")
	ErrEditNodeNotPermitted   = errors.New("you can't edit another user's node")

	ErrSensorNotFound           = errors.New("sensor not found")
	ErrDeleteSensorNotPermitted = errors.New("you can't delete another user's sensor")
	ErrUseSensorNotPermitted    = errors.New("you can't use another user's sensor")
	ErrSeeSensorNotPermitted    = errors.New("you can't see another user's sensor")
	ErrEditSensorNotPermitted   = errors.New("you can't edit another user's sensor")
)

func PanicIfNeeded(err interface{}) {
	if err != nil {
		panic(err)
	}
}
