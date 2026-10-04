package domain

type User struct {
	ID       int64
	Version  int64
	Login    string
	Password string
	Fullname string
}
