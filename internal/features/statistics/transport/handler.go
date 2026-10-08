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

// GetByUser     godoc
// @Summary      Get by user
// @Description  Get summary statistics of all user's accounts
// @Tags         statistics
// @Security     BearerAuth
// @Success      200 {object} UserStatisticsSliceDTO "Summary user statistics"
// @Failure      404 {object} utils.HTTPError "User not found"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /statistics [get]
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

// GetByAccount  godoc
// @Summary      Get by account
// @Description  Get statistics of the particular account
// @Tags         statistics
// @Security     BearerAuth
// @Param        accountID path int true "Account ID"
// @Success      200 {object} AccountStatisticsSliceDTO "Account statistics"
// @Failure      404 {object} utils.HTTPError "Account not found"
// @Failure      404 {object} utils.HTTPError "User not found"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /accounts/{accountID}/statistics [get]
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
