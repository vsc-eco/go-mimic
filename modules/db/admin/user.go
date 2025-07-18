package admin

type UserQueries interface {
	CreateUser(UserCredential) error
}

type UserCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
