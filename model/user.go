package model

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type UserLogin struct {
	Username    string `json:"username"`
	Password string `json:"password"`
}