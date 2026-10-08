package middleware

import (
	"net/http"

	"github.com/SJMaks/financer_back/internal/core/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommonMiddleware struct {
	log *zap.SugaredLogger
	cfg *config.Config
}

func NewCommonMiddleware(log *zap.SugaredLogger, cfg *config.Config) *CommonMiddleware {
	return &CommonMiddleware{
		log: log,
		cfg: cfg,
	}
}

func (m *CommonMiddleware) CORS(c *gin.Context) {
	origin := c.GetHeader("Origin")
	if origin == "" {
		c.Next()
		return
	}

	_, allowed := m.cfg.AllowedOrigins[origin]
	if allowed {
		h := c.Writer.Header()
		h.Set("Access-Control-Allow-Origin", origin)
		h.Set("Vary", "Origin")
	}

	if c.Request.Method == http.MethodOptions {
		if !allowed {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		h := c.Writer.Header()
		h.Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		if reqH := c.GetHeader("Access-Control-Request-Headers"); reqH != "" {
			h.Set("Access-Control-Allow-Headers", reqH)
		}
		h.Set("Access-Control-Max-Age", "86400")
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.Next()
}
