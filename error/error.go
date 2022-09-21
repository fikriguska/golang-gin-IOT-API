package error

import "errors"

// const (
// 	INVALID_PARAMS = 400
// 	INVALID_EMAIL  = 1001
// )

var (
	ErrInvalidParams            = errors.New("parameter is invalid")
	ErrInvalidEmail             = errors.New("email format is incorrect")
	ErrEmailUsernameAlreadyUsed = errors.New("email or Username already used")
	ErrAddUserFail              = errors.New("add user failed")
	// ErrInvalidToken             = errors.New("token is invalid")
	ErrActivateUserFail         = errors.New("activate user failed")
	ErrUsernameOrPassIncorrect  = errors.New("username not found or password incorrect")
	ErrUsernameOrEmailIncorrect = errors.New("username or email incorrect")
	ErrUserNotActive            = errors.New("account is inactive, check email for activation")
	ErrTokenNotFound            = errors.New("token not found")
	ErrUserAlreadyActive        = errors.New("your account has already activated")
	ErrUserIdNotFound           = errors.New("id user not found")
	ErrEditUserNotPermitted     = errors.New("can't edit another user's account")
	ErrDeleteUserNotPermitted   = errors.New("can't delete another user's account")
	ErrUserStillUsingNode       = errors.New("can't delete, you still using node")
	ErrOldPasswordIncorrect     = errors.New("old password is incorrect")
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

	ErrNotAdministrator = errors.New("you are not administrator")
)

func PanicIfNeeded(err interface{}) {
	if err != nil {
		panic(err)
	}
}
