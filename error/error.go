package error

import "errors"

// const (
// 	INVALID_PARAMS = 400
// 	INVALID_EMAIL  = 1001
// )

var (
	ErrInvalidParams = errors.New("parameter is invalid")
	ErrInvalidEmail  = errors.New("email is invalid")
	ErrUserExist     = errors.New("email or Username already exists")
	ErrAddUserFail   = errors.New("add user failed")
)
