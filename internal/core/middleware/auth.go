package middleware

import (
	"net/http"
	"strings"

	"github.com/SJMaks/financer_back/internal/core/auth"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tm auth.TokenManager
}

func NewAuthMiddleware(tm auth.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		tm: tm,
	}
}

type ctxKey string

const userCtxKey ctxKey = "auth.user"

func (a *AuthMiddleware) CheckAuth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no authorization header"})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
		return
	}

	payload, err := a.tm.Parse(headerParts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization token"})
		return
	}

	c.Set(string(userCtxKey), payload)
	c.Next()
}

func UserFromContext(c *gin.Context) (*auth.TokenPayload, bool) {
	v, ok := c.Get(string(userCtxKey))
	if !ok {
		return nil, false
	}
	p, ok := v.(*auth.TokenPayload)
	return p, ok
}
