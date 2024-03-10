package models

type Request struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	Username string `json:"usernames"`
	Email    string `json:"email"`
	JwtToken string `json:"jwttoken"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Points   int    `json:"points"`
}

type UserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Points   int    `json:"points"`
}
