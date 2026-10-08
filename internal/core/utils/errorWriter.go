package utils

import (
	"errors"
	"net/http"

	"github.com/SJMaks/financer_back/internal/core/apperrors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HTTPError struct {
	Error string `json:"error" example:"error message"`
	Code  string `json:"code" example:"error_code"`
}

func WriteHTTPError(c *gin.Context, log *zap.SugaredLogger, err error) {
	var svcErr *apperrors.Error
	if errors.As(err, &svcErr) {
		c.AbortWithStatusJSON(svcErr.Status, HTTPError{
			Error: svcErr.Msg,
			Code:  svcErr.Code,
		})
		return
	}
	log.Error(err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, HTTPError{
		Error: "internal error",
		Code:  "internal",
	})
}
