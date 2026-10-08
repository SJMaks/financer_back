package transport

import (
	"net/http"
	"strconv"

	"github.com/SJMaks/financer_back/internal/core/apperrors"
	"github.com/SJMaks/financer_back/internal/core/middleware"
	"github.com/SJMaks/financer_back/internal/core/utils"
	"github.com/SJMaks/financer_back/internal/features/categories/service"
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
// @Summary      Create category
// @Description  Create a new category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateCategoryDTO true "Create category request body"
// @Success      200 {object} ResponseCategoryDTO "Created category information"
// @Failure      400 {object} utils.HTTPError "Bad request"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /categories [post]
func (h *Handler) Create(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}
	dto := &CreateCategoryDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		appErr := apperrors.BadJSON
		appErr.Msg = err.Error()
		utils.WriteHTTPError(c, h.log, appErr)
		return
	}

	ctg := CreateCategoryDTOToDomain(dto)
	ctg.UserID = p.ID
	res, err := h.srv.Create(c.Request.Context(), ctg)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.JSON(http.StatusCreated, DomainToResponseDTO(res))
}

// GetAllOfUser  godoc
// @Summary      Get categories
// @Description  Get all categories of the current user
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} ResponseCategoriesDTO "All user's categories"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /categories [get]
func (h *Handler) GetAllOfUser(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}
	ctgs, err := h.srv.GetAllOfUser(c.Request.Context(), p.ID)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, DomainsToResponseDTO(ctgs))
}

// GetByID       godoc
// @Summary      Get category
// @Description  Get one category by ID
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Category ID"
// @Success      200 {object} ResponseCategoryDTO "Category information"
// @Failure      404 {object} utils.HTTPError "Category not found"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /categories/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}
	idStr, ok := c.Params.Get("id")
	if !ok {
		utils.WriteHTTPError(c, h.log, service.MissingID)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteHTTPError(c, h.log, service.InvalidID)
		return
	}
	ctg, err := h.srv.GetByID(c.Request.Context(), int64(id), p.ID)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, DomainToResponseDTO(ctg))
}

// UpdateByID    godoc
// @Summary      Update category
// @Description  Update an category by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Category ID"
// @Param        request body UpdateCategoryDTO true "Update category body"
// @Success      200 {object} ResponseCategoryDTO "Updated category information"
// @Failure      400 {object} utils.HTTPError "Bad request"
// @Failure      404 {object} utils.HTTPError "Category not found"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /categories/{id} [patch]
func (h *Handler) UpdateByID(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}
	idStr, ok := c.Params.Get("id")
	if !ok {
		utils.WriteHTTPError(c, h.log, service.MissingID)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteHTTPError(c, h.log, service.InvalidID)
		return
	}
	dto := &UpdateCategoryDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		appErr := apperrors.BadJSON
		appErr.Msg = err.Error()
		utils.WriteHTTPError(c, h.log, appErr)
		return
	}

	ctg := UpdateCategoryDTOToDomain(dto)
	ctg.ID = int64(id)
	ctg.UserID = p.ID
	res, err := h.srv.UpdateByID(c.Request.Context(), ctg)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, DomainToResponseDTO(res))
}

// DeleteByID    godoc
// @Summary      Delete category
// @Description  Delete an category by ID
// @Tags         categories
// @Security     BearerAuth
// @Param        id path int true "Category ID"
// @Success      204
// @Failure      404 {object} utils.HTTPError "Category not found"
// @Failure      401 {object} utils.HTTPError "Unauthorized"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /categories/{id} [delete]
func (h *Handler) DeleteByID(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}
	idStr, ok := c.Params.Get("id")
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
