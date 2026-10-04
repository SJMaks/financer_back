package transport

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SJMaks/financer_back/internal/core/middleware"
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
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	dto := &CreateOperationDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accIDStr, ok := c.Params.Get("accountID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.MissingID.Error()})
		return
	}
	accID, err := strconv.Atoi(accIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.InvalidID.Error()})
		return
	}

	op := CreateOperationDTOToDomain(dto)
	op.UserID = p.ID
	op.AccountID = int64(accID)

	res, err := h.srv.Create(c.Request.Context(), op)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, DomainToResponseDTO(res))
}

func (h *Handler) GetAllOfAccount(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	accIDStr, ok := c.Params.Get("accountID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.MissingID.Error()})
		return
	}
	accID, err := strconv.Atoi(accIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.InvalidID.Error()})
		return
	}

	params, err := h.parseListOperationsParams(c, p.ID, int64(accID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ops, total, err := h.srv.GetAllOfAccount(c.Request.Context(), *params)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, DomainsToResponseDTO(ops, total, params.Limit, params.Offset))
}

func (h *Handler) GetByID(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	accIDStr, ok := c.Params.Get("accountID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.MissingID.Error()})
		return
	}
	accID, err := strconv.Atoi(accIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.InvalidID.Error()})
		return
	}

	idStr, ok := c.Params.Get("operationID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.MissingID.Error()})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.InvalidID.Error()})
		return
	}

	ctg, err := h.srv.GetByID(c.Request.Context(), int64(id), p.ID, int64(accID))
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, DomainToResponseDTO(ctg))
}

func (h *Handler) UpdateByID(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	accIDStr, ok := c.Params.Get("accountID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.MissingID.Error()})
		return
	}
	accID, err := strconv.Atoi(accIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.InvalidID.Error()})
		return
	}

	idStr, ok := c.Params.Get("operationID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.MissingID.Error()})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.InvalidID.Error()})
		return
	}

	dto := &UpdateOperationDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	op := UpdateOperationDTOToDomain(dto)
	op.ID = int64(id)
	op.UserID = p.ID
	op.AccountID = int64(accID)
	res, err := h.srv.UpdateByID(c.Request.Context(), op)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, DomainToResponseDTO(res))
}

func (h *Handler) DeleteByID(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	accIDStr, ok := c.Params.Get("accountID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.MissingID.Error()})
		return
	}
	accID, err := strconv.Atoi(accIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.InvalidID.Error()})
		return
	}

	idStr, ok := c.Params.Get("operationID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.MissingID.Error()})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.InvalidID.Error()})
		return
	}

	err = h.srv.DeleteByID(c.Request.Context(), int64(id), p.ID, int64(accID))
	if err != nil {
		h.writeServiceError(c, err)
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
			return nil, fmt.Errorf("invalid limit")
		}
		limit = v
	}
	params.Limit = limit

	offset := 0
	if s := c.Query("offset"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v < 0 {
			return nil, fmt.Errorf("invalid offset")
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
		return nil, fmt.Errorf("invalid type: %q", t)
	}

	// --- category_id ---
	if s := c.Query("category_id"); s != "" {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil || v <= 0 {
			return nil, fmt.Errorf("invalid category_id")
		}
		params.Filter.CategoryID = &v
	}

	// --- min_date / max_date ---
	if s := c.Query("min_date"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return nil, fmt.Errorf("invalid min_date")
		}
		params.Filter.MinExecAt = &t
	}
	if s := c.Query("max_date"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return nil, fmt.Errorf("invalid max_date")
		}
		params.Filter.MaxExecAt = &t
	}
	if params.Filter.MinExecAt != nil && params.Filter.MaxExecAt != nil &&
		params.Filter.MaxExecAt.Before(*params.Filter.MinExecAt) {
		return nil, fmt.Errorf("max_date before min_date")
	}

	// --- sort / order ---
	params.SortBy = c.DefaultQuery("sort", "exec_at")
	params.SortOrd = c.DefaultQuery("order", "desc")

	return params, nil
}

func (h *Handler) writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.OperationNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	case errors.Is(err, service.EmptyTitle) ||
		errors.Is(err, service.MissingAccountID) ||
		errors.Is(err, service.MissingCategoryID) ||
		errors.Is(err, service.NegativeSum):
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case errors.Is(err, service.UpdateVersionConflict):
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	default:
		h.log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
}
