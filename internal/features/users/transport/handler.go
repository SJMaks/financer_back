package transport

import (
	"context"
	"net/http"

	"github.com/SJMaks/financer_back/internal/core/apperrors"
	"github.com/SJMaks/financer_back/internal/core/auth"
	"github.com/SJMaks/financer_back/internal/core/utils"
	"github.com/SJMaks/financer_back/internal/features/users/domain"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserService interface {
	Register(ctx context.Context, u *domain.User, password string) (*domain.User, error)
	Login(ctx context.Context, u *domain.User) (string, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	UpdateByID(ctx context.Context, u *domain.PatchUser) (*domain.User, error)
	DeleteByID(ctx context.Context, id int64) error
}

type AuthMiddleware interface {
	UserFromContext(c *gin.Context) (*auth.TokenPayload, bool)
}

type Handler struct {
	srv UserService
	log *zap.SugaredLogger
	mid AuthMiddleware
}

func NewHandler(srv UserService, log *zap.SugaredLogger, mid AuthMiddleware) *Handler {
	return &Handler{
		srv: srv,
		log: log,
		mid: mid,
	}
}

func (h *Handler) Register(c *gin.Context) {
	dto := &UserRegisterDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		appErr := apperrors.BadJSON
		appErr.Msg = err.Error()
		utils.WriteHTTPError(c, h.log, appErr)
		return
	}

	user := RegisterDTOToDomain(*dto)
	res, err := h.srv.Register(c.Request.Context(), user, dto.Password)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	u := DomainToResponse(*res)

	c.JSON(http.StatusCreated, u)
}

func (h *Handler) Login(c *gin.Context) {
	dto := &LoginRequestDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		appErr := apperrors.BadJSON
		appErr.Msg = err.Error()
		utils.WriteHTTPError(c, h.log, appErr)
		return
	}

	user := LoginDTOToDomain(*dto)
	token, err := h.srv.Login(c.Request.Context(), user)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, LoginResponseDTO{
		Token: token,
	})
}

func (h *Handler) GetMe(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}
	u, err := h.srv.GetByID(c.Request.Context(), p.ID)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	res := DomainToResponse(*u)
	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateMe(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}

	dto := &UpdateRequestDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		appErr := apperrors.BadJSON
		appErr.Msg = err.Error()
		utils.WriteHTTPError(c, h.log, appErr)
		return
	}

	user := UpdateDTOToDomain(*dto)
	user.ID = p.ID
	u, err := h.srv.UpdateByID(c.Request.Context(), user)
	if err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	res := DomainToResponse(*u)
	c.JSON(http.StatusOK, res)
}

func (h *Handler) DeleteMe(c *gin.Context) {
	p, ok := h.mid.UserFromContext(c)
	if !ok {
		utils.WriteHTTPError(c, h.log, apperrors.Unauthorized)
		return
	}
	if err := h.srv.DeleteByID(c.Request.Context(), p.ID); err != nil {
		utils.WriteHTTPError(c, h.log, err)
		return
	}
	c.Status(http.StatusNoContent)
}
