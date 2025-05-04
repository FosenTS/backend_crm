package model

type Role int8

const (
	Director Role = iota
	Employee
)

type User struct {
	UserId   string
	Role     Role
	Username string
	PassHash string
}
