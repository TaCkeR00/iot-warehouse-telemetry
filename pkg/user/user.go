package user

type User struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is-admin"`
}
