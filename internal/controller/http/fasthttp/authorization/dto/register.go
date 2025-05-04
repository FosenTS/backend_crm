package dto

type Register struct {
	RoleId   int    `json:"role_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
