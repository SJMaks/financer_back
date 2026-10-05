package middleware

import (
	"strings"

	"github.com/SJMaks/financer_back/internal/core/apperrors"
	"github.com/SJMaks/financer_back/internal/core/auth"
	"github.com/SJMaks/financer_back/internal/core/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	log *zap.SugaredLogger
	tm  auth.TokenManager
}

func NewAuthMiddleware(log *zap.SugaredLogger, tm auth.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		log: log,
		tm:  tm,
	}
}

type ctxKey string

const userCtxKey ctxKey = "auth.user"

func (a *AuthMiddleware) CheckAuth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		utils.WriteHTTPError(c, a.log, apperrors.Unauthorized)
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		utils.WriteHTTPError(c, a.log, apperrors.Unauthorized)
		return
	}

	payload, err := a.tm.Parse(headerParts[1])
	if err != nil {
		utils.WriteHTTPError(c, a.log, apperrors.Unauthorized)
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
