package transport

import (
	"net/http"

	"github.com/SJMaks/financer_back/internal/features/operation_types/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	log *zap.SugaredLogger
	srv *service.Service
}

func NewHandler(log *zap.SugaredLogger, srv *service.Service) *Handler {
	return &Handler{
		log: log,
		srv: srv,
	}
}
func (h *Handler) GetAll(c *gin.Context) {
	types, err := h.srv.GetAll(c.Request.Context())
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, DomainsToResponseDTO(types))
}

func (h *Handler) writeServiceError(c *gin.Context, err error) {
	switch {
	default:
		h.log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
}
