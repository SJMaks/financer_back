package auth

import (
	"fmt"
	"time"

	"github.com/SJMaks/financer_back/internal/core/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Payload TokenPayload `json:"payload"`
	jwt.RegisteredClaims
}

type JWTTokenManager struct {
	cfg *config.Config
}

func NewJWT(cfg *config.Config) *JWTTokenManager {
	return &JWTTokenManager{
		cfg: cfg,
	}
}

var _ TokenManager = (*JWTTokenManager)(nil)

func (tm *JWTTokenManager) Generate(p TokenPayload) (string, error) {
	regClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "financer",
		Subject:   p.Login,
	}
	claims := JWTClaims{
		Payload:          p,
		RegisteredClaims: regClaims,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.cfg.JWTSecret))
}

func (tm *JWTTokenManager) Parse(t string) (*TokenPayload, error) {
	claims := &JWTClaims{}
	_, err := jwt.ParseWithClaims(t, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(tm.cfg.JWTSecret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, err
	}
	return &claims.Payload, nil
}
