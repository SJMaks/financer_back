package transport

import (
	"net/http"

	"github.com/SJMaks/financer_back/internal/features/users/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	srv *service.Service
}

func NewHandler(srv *service.Service) *Handler {
	return &Handler{
		srv: srv,
	}
}

func (h Handler) Ping(c *gin.Context) {
	answer, err := h.srv.Ping()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusCreated, answer)
}
