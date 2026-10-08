package transport

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SJMaks/financer_back/internal/core/apperrors"
	"github.com/SJMaks/financer_back/internal/core/middleware"
	"github.com/SJMaks/financer_back/internal/core/utils"
	"github.com/SJMaks/financer_back/internal/features/operations/service"
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
// @Summary      Create operation
// @Description  Create a new operation
// @Tags         operations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        accountID path int true "Account ID"
// @Param        request body CreateOperationDTO true "Create operation request body"
// @Success      200 {object} ResponseOperationDTO "Created operation information"
// @Failure      404 {object} utils.HTTPError "Account not found"
// @Failure      400 {object} utils.HTTPError "Bad request"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /accounts/{accountID}/operations [post]
func (h *Handler) Create(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}
	dto := &CreateOperationDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		appErr := apperrors.BadJSON
		appErr.Msg = err.Error()
		utils.WriteHTTPError(c, h.log, appErr)
		return
	}

	accIDStr, ok := c.Params.Get("accountID")
	if !ok {
		utils.WriteHTTPError(c, h.log, service.MissingAccountID)
		return
	}
	accID, err := strconv.Atoi(accIDStr)
	if err != nil {
		utils.WriteHTTPError(c, h.log, service.InvalidAccountID)
		return
	}

	op := CreateOperationDTOToDomain(dto)
	op.UserID = p.ID
	op.AccountID = int64(accID)

	res, err := h.srv.Create(c.Request.Context(), op)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.JSON(http.StatusCreated, DomainToResponseDTO(res))
}

// GetAllOfUser  godoc
// @Summary      Get operations
// @Description  Get all operations of the current user
// @Tags         operations
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Pagination: Limit"
// @Param        offset query int false "Pagination: Offset"
// @Param        type query string false "Filter: Type of operation. Can be either 'income' or 'expense'"
// @Param        category_id query int false "Filter: ID of the operation category"
// @Param        min_date query string false "Filter: min date of the date range"
// @Param        max_date query string false "Filter: max date of the date range. Cannot be before the min_date"
// @Param        sort query string false "Sort: can be either 'exec_at' or 'sum'"
// @Param        order query string false "Sort: Order of the sort. Can be either 'asc' or 'desc'"
// @Param        accountID path int true "Account ID"
// @Success      200 {object} ResponseOperationsDTO "All user's operations filtered and sorted"
// @Failure      404 {object} utils.HTTPError "Account not found"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /accounts/{accountID}/operations [get]
func (h *Handler) GetAllOfAccount(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}

	accIDStr, ok := c.Params.Get("accountID")
	if !ok {
		utils.WriteHTTPError(c, h.log, service.MissingAccountID)
		return
	}
	accID, err := strconv.Atoi(accIDStr)
	if err != nil {
		utils.WriteHTTPError(c, h.log, service.InvalidAccountID)
		return
	}

	params, err := h.parseListOperationsParams(c, p.ID, int64(accID))
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}

	ops, total, err := h.srv.GetAllOfAccount(c.Request.Context(), *params)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, DomainsToResponseDTO(ops, total, params.Limit, params.Offset))
}

// GetByID       godoc
// @Summary      Get operation
// @Description  Get one operation by ID
// @Tags         operations
// @Produce      json
// @Security     BearerAuth
// @Param        accountID path int true "Account ID"
// @Param        operationID path int true "Operation ID"
// @Success      200 {object} ResponseOperationDTO "Operation information"
// @Failure      404 {object} utils.HTTPError "Account not found"
// @Failure      404 {object} utils.HTTPError "Operation not found"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /accounts/{accountID}/operations/{operationID} [get]
func (h *Handler) GetByID(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}

	accIDStr, ok := c.Params.Get("accountID")
	if !ok {
		utils.WriteHTTPError(c, h.log, service.MissingAccountID)
		return
	}
	accID, err := strconv.Atoi(accIDStr)
	if err != nil {
		utils.WriteHTTPError(c, h.log, service.InvalidAccountID)
		return
	}

	idStr, ok := c.Params.Get("operationID")
	if !ok {
		utils.WriteHTTPError(c, h.log, service.MissingOperationID)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteHTTPError(c, h.log, service.InvalidOperationID)
		return
	}

	ctg, err := h.srv.GetByID(c.Request.Context(), int64(id), p.ID, int64(accID))
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, DomainToResponseDTO(ctg))
}

// UpdateByID    godoc
// @Summary      Update operation
// @Description  Update an operation by ID
// @Tags         operations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        accountID path int true "Account ID"
// @Param        operationID path int true "Operation ID"
// @Param        request body UpdateOperationDTO true "Update operation body"
// @Success      200 {object} ResponseOperationDTO "Updated operation information"
// @Failure      404 {object} utils.HTTPError "Account not found"
// @Failure      400 {object} utils.HTTPError "Bad request"
// @Failure      404 {object} utils.HTTPError "Operation not found"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /accounts/{accountID}/operations/{operationID} [patch]
func (h *Handler) UpdateByID(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}

	accIDStr, ok := c.Params.Get("accountID")
	if !ok {
		utils.WriteHTTPError(c, h.log, service.MissingAccountID)
		return
	}
	accID, err := strconv.Atoi(accIDStr)
	if err != nil {
		utils.WriteHTTPError(c, h.log, service.InvalidAccountID)
		return
	}

	idStr, ok := c.Params.Get("operationID")
	if !ok {
		utils.WriteHTTPError(c, h.log, service.MissingOperationID)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteHTTPError(c, h.log, service.InvalidOperationID)
		return
	}

	dto := &UpdateOperationDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		appErr := apperrors.BadJSON
		appErr.Msg = err.Error()
		utils.WriteHTTPError(c, h.log, appErr)
		return
	}

	op := UpdateOperationDTOToDomain(dto)
	op.ID = int64(id)
	op.UserID = p.ID
	op.AccountID = int64(accID)
	res, err := h.srv.UpdateByID(c.Request.Context(), op)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, DomainToResponseDTO(res))
}

// DeleteByID    godoc
// @Summary      Delete operation
// @Description  Delete an operation by ID
// @Tags         operations
// @Security     BearerAuth
// @Param        accountID path int true "Account ID"
// @Param        operationID path int true "Operation ID"
// @Success      204
// @Failure      404 {object} utils.HTTPError "Account not found"
// @Failure      404 {object} utils.HTTPError "Operation not found"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /accounts/{accountID}/operations/{operationID} [delete]
func (h *Handler) DeleteByID(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}

	accIDStr, ok := c.Params.Get("accountID")
	if !ok {
		utils.WriteHTTPError(c, h.log, service.MissingAccountID)
		return
	}
	accID, err := strconv.Atoi(accIDStr)
	if err != nil {
		utils.WriteHTTPError(c, h.log, service.InvalidAccountID)
		return
	}

	idStr, ok := c.Params.Get("operationID")
	if !ok {
		utils.WriteHTTPError(c, h.log, service.MissingOperationID)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteHTTPError(c, h.log, service.InvalidOperationID)
		return
	}

	err = h.srv.DeleteByID(c.Request.Context(), int64(id), p.ID, int64(accID))
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) parseListOperationsParams(c *gin.Context, userID, accountID int64) (*service.ListOperationsParams, error) {
	params := &service.ListOperationsParams{
		Filter: service.ListOperationsFilter{
			UserID:    userID,
			AccountID: accountID,
		},
	}

	// --- limit / offset ---
	limit := 50
	if s := c.Query("limit"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v <= 0 || v > 100 {
			appErr := apperrors.IvalidQuery
			appErr.Msg = "ivalid limit"
			return nil, appErr
		}
		limit = v
	}
	params.Limit = limit

	offset := 0
	if s := c.Query("offset"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v < 0 {
			appErr := apperrors.IvalidQuery
			appErr.Msg = "ivalid offset"
			return nil, appErr
		}
		offset = v
	}
	params.Offset = offset

	// --- type: income | expense ---
	switch t := c.Query("type"); t {
	case "":
		// не задано — не фильтруем
	case "income":
		id := int64(1)
		params.Filter.OperationTypeID = &id
	case "expense":
		id := int64(2)
		params.Filter.OperationTypeID = &id
	default:
		appErr := apperrors.IvalidQuery
		appErr.Msg = fmt.Sprintf("invalid type: %q", t)
		return nil, appErr
	}

	// --- category_id ---
	if s := c.Query("category_id"); s != "" {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil || v <= 0 {
			appErr := apperrors.IvalidQuery
			appErr.Msg = "invalid category_id"
			return nil, appErr
		}
		params.Filter.CategoryID = &v
	}

	// --- min_date / max_date ---
	if s := c.Query("min_date"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			appErr := apperrors.IvalidQuery
			appErr.Msg = "invalid min_date"
			return nil, appErr
		}
		params.Filter.MinExecAt = &t
	}
	if s := c.Query("max_date"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			appErr := apperrors.IvalidQuery
			appErr.Msg = "invalid max_date"
			return nil, appErr
		}
		params.Filter.MaxExecAt = &t
	}
	if params.Filter.MinExecAt != nil && params.Filter.MaxExecAt != nil &&
		params.Filter.MaxExecAt.Before(*params.Filter.MinExecAt) {
		appErr := apperrors.IvalidQuery
		appErr.Msg = "max_date before min_date"
		return nil, appErr
	}

	// --- sort / order ---
	params.SortBy = c.DefaultQuery("sort", "exec_at")
	params.SortOrd = c.DefaultQuery("order", "asc")

	return params, nil
}
