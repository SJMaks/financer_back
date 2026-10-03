package auth

type TokenPayload struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

type TokenManager interface {
	Generate(p TokenPayload) (string, error)
	Parse(t string) (*TokenPayload, error)
}
