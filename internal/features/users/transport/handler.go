package transport

import (
	"errors"
	"net/http"

	"github.com/SJMaks/financer_back/internal/core/middleware"
	"github.com/SJMaks/financer_back/internal/features/users/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	srv *service.Service
	log *zap.SugaredLogger
}

func NewHandler(srv *service.Service, log *zap.SugaredLogger) *Handler {
	return &Handler{
		srv: srv,
		log: log,
	}
}

func (h *Handler) Register(c *gin.Context) {
	dto := &UserRegisterDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := RegisterDTOToDomain(*dto)
	res, err := h.srv.Register(c.Request.Context(), user, dto.Password)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	u := DomainToResponse(*res)

	c.JSON(http.StatusCreated, u)
}

func (h *Handler) Login(c *gin.Context) {
	dto := &LoginRequestDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := LoginDTOToDomain(*dto)
	token, err := h.srv.Login(c.Request.Context(), user)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, LoginResponseDTO{
		Token: token,
	})
}

func (h *Handler) GetMe(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	u, err := h.srv.GetByID(c.Request.Context(), p.ID)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	res := DomainToResponse(*u)
	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateMe(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	dto := &UpdateRequestDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := UpdateDTOToDomain(*dto)
	user.ID = p.ID
	u, err := h.srv.UpdateByID(c.Request.Context(), user)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	res := DomainToResponse(*u)
	c.JSON(http.StatusOK, res)
}

func (h *Handler) DeleteMe(c *gin.Context) {
	p, ok := middleware.UserFromContext(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.srv.DeleteByID(c.Request.Context(), p.ID); err != nil {
		h.writeServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.EmptyLogin) ||
		errors.Is(err, service.EmptyFullname) ||
		errors.Is(err, service.WrongPasswordFormat):
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case errors.Is(err, service.UserNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	case errors.Is(err, service.InvalidCredentials):
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	case errors.Is(err, service.UserAlreadyExists):
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	default:
		h.log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
}
