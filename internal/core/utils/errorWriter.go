package utils

import (
	"errors"
	"net/http"

	"github.com/SJMaks/financer_back/internal/core/apperrors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func WriteHTTPError(c *gin.Context, log *zap.SugaredLogger, err error) {
	var svcErr *apperrors.Error
	if errors.As(err, &svcErr) {
		c.AbortWithStatusJSON(svcErr.Status, gin.H{
			"error": svcErr.Msg,
			"code":  svcErr.Code,
		})
		return
	}
	log.Error(err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"error": "internal error",
		"code":  "internal",
	})
}
