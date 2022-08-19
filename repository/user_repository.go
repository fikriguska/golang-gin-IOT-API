package repository

type UserRepo interface {
	Add(map[string]interface{})
	IsUsernameExist(string) bool
	GetIdByUsername(string) int
	IsTokenExist(string) bool
	IsActivatedCheckByToken(string) bool
	IsActivatedCheckByUsername(string) bool
	Activate(string)
	Auth(string, string) bool
}
