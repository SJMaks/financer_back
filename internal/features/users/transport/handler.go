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

// Register      godoc
// @Summary      Register a user
// @Description  Resgister a new user in the app
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body UserRegisterDTO true "Register request body"
// @Success      201 {object} UserResponseDTO "Successfully registered user"
// @Failure      400 {object} utils.HTTPError "Bad request"
// @Failure      409 {object} utils.HTTPError "User already exists"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /auth/register [post]
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

// Login         godoc
// @Summary      Log in
// @Description  Log in to the app
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body LoginRequestDTO true "Login request body"
// @Success      200 {object} LoginResponseDTO "JWT of a successfully authenticated user"
// @Failure      400 {object} utils.HTTPError "Bad request"
// @Failure      401 {object} utils.HTTPError "Invalid credentials"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /auth/login [post]
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

// GetMe         godoc
// @Summary      Get current
// @Description  Get current user information
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} UserResponseDTO "Current user information"
// @Failure      404 {object} utils.HTTPError "User not found"
// @Failure      401 {object} utils.HTTPError "Invalid JWT"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /users/me [get]
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

// UpdateMe      godoc
// @Summary      Update current
// @Description  Update current user information
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body UpdateRequestDTO true "Update request body"
// @Success      200 {object} UserResponseDTO "Updated user information"
// @Failure      400 {object} utils.HTTPError "Bad request"
// @Failure      404 {object} utils.HTTPError "User not found"
// @Failure      401 {object} utils.HTTPError "Invalid JWT"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /users/me [patch]
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

// DeleteMe      godoc
// @Summary      Delete current
// @Description  Delete current user
// @Tags         users
// @Security     BearerAuth
// @Success      204 "User successfully deleted"
// @Failure      404 {object} utils.HTTPError "User not found"
// @Failure      401 {object} utils.HTTPError "Invalid JWT"
// @Failure      500 {object} utils.HTTPError "Internal error"
// @Router       /users/me [delete]
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
