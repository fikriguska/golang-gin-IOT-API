package repository

type User struct {
	Id_user  int    `json:"id_user"`
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Status   bool   `json:"status"`
	Token    string `json:"token"`
}

type UserRepository interface {
	Create(user User)
	GetByUsername(string) User
	GetPassword(string) string
	GetStatusByUsername(string) bool
}
