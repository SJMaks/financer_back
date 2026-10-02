package service

type User struct {
	ID           int64
	Version      uint64
	Login        string
	PasswordHash string
	Fullname     string
}
