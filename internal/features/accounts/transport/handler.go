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
	mid *middleware.AuthMiddleware
}

func NewHandler(log *zap.SugaredLogger, srv *service.Service, mid *middleware.AuthMiddleware) *Handler {
	return &Handler{
		log: log,
		srv: srv,
		mid: mid,
	}
}

// Create        godoc
// @Summary      Create account
// @Description  Create a new account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateAccountDTO true "Create account request body"
// @Success      200 {object} ResponseAccountDTO "Created account information"
// @Failure      400 {object} utils.HTTPError "Bad request"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /accounts [post]
func (h *Handler) Create(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
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

// GetAllOfUser  godoc
// @Summary      Get accounts
// @Description  Get all accounts of the current user
// @Tags         accounts
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} ResponseAccountsDTO "All user's accounts"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /accounts [get]
func (h *Handler) GetAllOfUser(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
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

// GetByID       godoc
// @Summary      Get account
// @Description  Get one account by ID
// @Tags         accounts
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Account ID"
// @Success      200 {object} ResponseAccountDTO "Account information"
// @Failure      404 {object} utils.HTTPError "Account not found"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /accounts/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
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
	acc, err := h.srv.GetByID(c.Request.Context(), int64(id), p.ID)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, DomainToResponseDTO(acc))
}

// UpdateByID    godoc
// @Summary      Update account
// @Description  Update an account by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Account ID"
// @Param        request body UpdateAccountDTO true "Update account body"
// @Success      200 {object} ResponseAccountDTO "Updated account information"
// @Failure      400 {object} utils.HTTPError "Bad request"
// @Failure      404 {object} utils.HTTPError "Account not found"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /accounts/{id} [patch]
func (h *Handler) UpdateByID(c *gin.Context) {
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

// DeleteByID    godoc
// @Summary      Delete account
// @Description  Delete an account by ID
// @Tags         accounts
// @Security     BearerAuth
// @Param        id path int true "Account ID"
// @Success      204
// @Failure      404 {object} utils.HTTPError "Account not found"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /accounts/{id} [delete]
func (h *Handler) DeleteByID(c *gin.Context) {
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
	err = h.srv.DeleteByID(c.Request.Context(), int64(id), p.ID)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.Status(http.StatusNoContent)
}
