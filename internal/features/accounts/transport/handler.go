package transport

import (
	"net/http"
	"strconv"

	"github.com/SJMaks/financer_back/internal/core/apperrors"
	"github.com/SJMaks/financer_back/internal/core/middleware"
	"github.com/SJMaks/financer_back/internal/core/utils"
	"github.com/SJMaks/financer_back/internal/features/accounts/service"
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

func (h *Handler) Create(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}
	dto := &CreateAccountDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		appErr := apperrors.BadJSON
		appErr.Msg = err.Error()
		utils.WriteHTTPError(c, h.log, appErr)
		return
	}

	account := CreateAccountDTOToDomain(dto)
	account.UserID = p.ID
	res, err := h.srv.Create(c.Request.Context(), account)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.JSON(http.StatusCreated, DomainToResponseDTO(res))
}

func (h *Handler) GetAllOfUser(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}
	accs, err := h.srv.GetAllOfUser(c.Request.Context(), p.ID)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, DomainsToResponseDTO(accs))
}

func (h *Handler) GetByID(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
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
	acc, err := h.srv.GetByID(c.Request.Context(), int64(id), p.ID)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, DomainToResponseDTO(acc))
}

func (h *Handler) UpdateByID(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
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
	dto := &UpdateAccountDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		appErr := apperrors.BadJSON
		appErr.Msg = err.Error()
		utils.WriteHTTPError(c, h.log, appErr)
		return
	}

	account := UpdateAccountDTOToDomain(dto)
	account.ID = int64(id)
	account.UserID = p.ID
	acc, err := h.srv.UpdateByID(c.Request.Context(), account)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, DomainToResponseDTO(acc))
}

func (h *Handler) DeleteByID(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
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
	err = h.srv.DeleteByID(c.Request.Context(), int64(id), p.ID)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.Status(http.StatusNoContent)
}
