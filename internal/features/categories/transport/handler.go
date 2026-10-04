package transport

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SJMaks/financer_back/internal/core/middleware"
	"github.com/SJMaks/financer_back/internal/features/categories/service"
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
	dto := &CreateCategoryDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctg := CreateCategoryDTOToDomain(dto)
	ctg.UserID = p.ID
	res, err := h.srv.Create(c.Request.Context(), ctg)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, DomainToResponseDTO(res))
}

func (h *Handler) GetAllOfUser(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ctgs, err := h.srv.GetAllOfUser(c.Request.Context(), p.ID)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, DomainsToResponseDTO(ctgs))
}

func (h *Handler) GetByID(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	idStr, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.MissingID.Error()})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.InvalidID.Error()})
		return
	}
	ctg, err := h.srv.GetByID(c.Request.Context(), int64(id), p.ID)
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
	idStr, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.MissingID.Error()})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.InvalidID.Error()})
		return
	}
	dto := &UpdateCategoryDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctg := UpdateCategoryDTOToDomain(dto)
	ctg.ID = int64(id)
	ctg.UserID = p.ID
	res, err := h.srv.UpdateByID(c.Request.Context(), ctg)
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
	idStr, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.MissingID.Error()})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": service.InvalidID.Error()})
		return
	}
	err = h.srv.DeleteByID(c.Request.Context(), int64(id), p.ID)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.CategoryNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	case errors.Is(err, service.EmptyTitle) ||
		errors.Is(err, service.MissingOperationTypeID):
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
