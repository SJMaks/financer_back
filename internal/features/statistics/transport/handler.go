package transport

import (
	"net/http"
	"strconv"

	"github.com/SJMaks/financer_back/internal/core/apperrors"
	"github.com/SJMaks/financer_back/internal/core/middleware"
	"github.com/SJMaks/financer_back/internal/core/utils"
	"github.com/SJMaks/financer_back/internal/features/statistics/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	log *zap.SugaredLogger
	srv *service.Service
	mid *middleware.AuthMiddleware
}

func NewHandler(log *zap.SugaredLogger, srv *service.Service, mid *middleware.AuthMiddleware) *Handler {
	return &Handler{
		log: log,
		srv: srv,
		mid: mid,
	}
}

func (h *Handler) GetByUser(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}
	stcs, err := h.srv.GetByUser(c.Request.Context(), p.ID)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, ByUserDomainsToSliceDTO(stcs))
}

func (h *Handler) GetByAccount(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}
	idStr, ok := c.Params.Get("accountID")
	if !ok {
		utils.WriteHTTPError(c, h.log, service.MissingID)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteHTTPError(c, h.log, service.InvalidID)
		return
	}
	stcs, err := h.srv.GetByAccount(c.Request.Context(), int64(id), p.ID)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, ByAccountDomainsToSliceDTO(stcs))
}
