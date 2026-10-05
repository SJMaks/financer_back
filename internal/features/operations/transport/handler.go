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

func (h *Handler) GetAllOfAccount(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
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

func (h *Handler) GetByID(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
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

func (h *Handler) UpdateByID(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
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

func (h *Handler) DeleteByID(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
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
