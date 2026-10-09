package transport_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SJMaks/financer_back/internal/core/auth"
	"github.com/SJMaks/financer_back/internal/features/users/domain"
	"github.com/SJMaks/financer_back/internal/features/users/service"
	"github.com/SJMaks/financer_back/internal/features/users/transport"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestHandler_Register(t *testing.T) {
	log := zap.NewNop().Sugar()

	user := &domain.User{
		ID: 1,
	}

	tests := []struct {
		name       string
		body       string
		setupMock  func(srv *MockUserService)
		wantStatus int
	}{
		{
			name: "success",
			body: `{
				"login": "alice",
				"password": "secret123",
				"full_name": "Alice"
			}`,
			setupMock: func(srv *MockUserService) {
				srv.EXPECT().Register(
					mock.Anything,
					&domain.User{
						Login:    "alice",
						Password: "secret123",
						Fullname: "Alice",
					},
					"secret123",
				).Return(user, nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "bad json",
			body: `{
				"login:
			}`,
			setupMock: func(srv *MockUserService) {
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service login required error",
			body: `{}`,
			setupMock: func(srv *MockUserService) {
				srv.EXPECT().Register(
					mock.Anything,
					&domain.User{},
					"",
				).Return(nil, service.LoginRequired)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service password required error",
			body: `{
				"login": "alice"
			}`,
			setupMock: func(srv *MockUserService) {
				srv.EXPECT().Register(
					mock.Anything,
					&domain.User{
						Login: "alice",
					},
					"",
				).Return(nil, service.PasswordRequired)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service full_name required error",
			body: `{
				"login": "alice",
				"password": "secret123"
			}`,
			setupMock: func(srv *MockUserService) {
				srv.EXPECT().Register(
					mock.Anything,
					&domain.User{
						Login:    "alice",
						Password: "secret123",
					},
					"secret123",
				).Return(nil, service.FullnameRequired)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service invalid login error",
			body: `{
				"login": "alice some",
				"password": "secret123",
				"full_name": "Alice"
			}`,
			setupMock: func(srv *MockUserService) {
				srv.EXPECT().Register(
					mock.Anything,
					&domain.User{
						Login:    "alice some",
						Password: "secret123",
						Fullname: "Alice",
					},
					"secret123",
				).Return(nil, service.InvalidLogin)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service invalid password error",
			body: `{
				"login": "alice",
				"password": "secret123 some",
				"full_name": "Alice"
			}`,
			setupMock: func(srv *MockUserService) {
				srv.EXPECT().Register(
					mock.Anything,
					&domain.User{
						Login:    "alice",
						Password: "secret123 some",
						Fullname: "Alice",
					},
					"secret123 some",
				).Return(nil, service.InvalidPassword)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service invalid full_name error",
			body: `{
				"login": "alice",
				"password": "secret123",
				"full_name": "Alice123"
			}`,
			setupMock: func(srv *MockUserService) {
				srv.EXPECT().Register(
					mock.Anything,
					&domain.User{
						Login:    "alice",
						Password: "secret123",
						Fullname: "Alice123",
					},
					"secret123",
				).Return(nil, service.InvalidFullname)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service internal error",
			body: `{
				"login": "alice",
				"password": "secret123",
				"full_name": "Alice"
			}`,
			setupMock: func(srv *MockUserService) {
				srv.EXPECT().Register(
					mock.Anything,
					&domain.User{
						Login:    "alice",
						Password: "secret123",
						Fullname: "Alice",
					},
					"secret123",
				).Return(nil, errors.New("something went wrong"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewMockUserService(t)
			mid := NewMockAuthMiddleware(t)

			tt.setupMock(srv)

			h := transport.NewHandler(srv, log, mid)

			r := gin.New()
			r.POST("/api/v1/auth/register", h.Register)

			req := httptest.NewRequest(
				"POST",
				"/api/v1/auth/register",
				strings.NewReader(tt.body))

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestHandler_Login(t *testing.T) {
	log := zap.NewNop().Sugar()

	tests := []struct {
		name       string
		body       string
		setupMock  func(srv *MockUserService)
		wantStatus int
	}{
		{
			name: "success",
			body: `{
				"login": "alice",
				"password": "secret123"
			}`,
			setupMock: func(srv *MockUserService) {
				srv.EXPECT().Login(
					mock.Anything,
					&domain.User{
						Login:    "alice",
						Password: "secret123",
					},
				).Return("token", nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "bad json",
			body: `{
				"login:
			}`,
			setupMock: func(srv *MockUserService) {
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service login required error",
			body: `{}`,
			setupMock: func(srv *MockUserService) {
				srv.EXPECT().Login(
					mock.Anything,
					&domain.User{},
				).Return("", service.LoginRequired)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service password required error",
			body: `{
				"login": "alice"
			}`,
			setupMock: func(srv *MockUserService) {
				srv.EXPECT().Login(
					mock.Anything,
					&domain.User{
						Login: "alice",
					},
				).Return("", service.PasswordRequired)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service invalid credentials error",
			body: `{
				"login": "alice",
				"password": "password"
			}`,
			setupMock: func(srv *MockUserService) {
				srv.EXPECT().Login(
					mock.Anything,
					&domain.User{
						Login:    "alice",
						Password: "password",
					},
				).Return("", service.InvalidCredentials)
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "service internal error",
			body: `{
				"login": "alice",
				"password": "secret123"
			}`,
			setupMock: func(srv *MockUserService) {
				srv.EXPECT().Login(
					mock.Anything,
					&domain.User{
						Login:    "alice",
						Password: "secret123",
					},
				).Return("", errors.New("something went wrong"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewMockUserService(t)
			mid := NewMockAuthMiddleware(t)

			tt.setupMock(srv)

			h := transport.NewHandler(srv, log, mid)

			r := gin.New()
			r.POST("/api/v1/auth/login", h.Login)

			req := httptest.NewRequest(
				"POST",
				"/api/v1/auth/login",
				strings.NewReader(tt.body))

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestHandler_GetMe(t *testing.T) {
	log := zap.NewNop().Sugar()

	user := &domain.User{
		ID:       1,
		Login:    "alice",
		Fullname: "Alice",
	}

	tests := []struct {
		name       string
		token      string
		setupMock  func(mid *MockAuthMiddleware, srv *MockUserService)
		wantStatus int
	}{
		{
			name:  "success",
			token: "1",
			setupMock: func(mid *MockAuthMiddleware, srv *MockUserService) {
				mid.EXPECT().UserFromContext(
					mock.Anything,
				).Return(&auth.TokenPayload{
					ID:    1,
					Login: "alice",
				}, true)
				srv.EXPECT().GetByID(
					mock.Anything,
					int64(1),
				).Return(user, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "invalid token",
			token: "invalid",
			setupMock: func(mid *MockAuthMiddleware, srv *MockUserService) {
				mid.EXPECT().UserFromContext(
					mock.Anything,
				).Return(&auth.TokenPayload{}, false)
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:  "service internal error",
			token: "1",
			setupMock: func(mid *MockAuthMiddleware, srv *MockUserService) {
				mid.EXPECT().UserFromContext(
					mock.Anything,
				).Return(&auth.TokenPayload{
					ID:    1,
					Login: "alice",
				}, true)
				srv.EXPECT().GetByID(
					mock.Anything,
					int64(1),
				).Return(nil, fmt.Errorf("something went wrong"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewMockUserService(t)
			mid := NewMockAuthMiddleware(t)

			tt.setupMock(mid, srv)

			h := transport.NewHandler(srv, log, mid)

			r := gin.New()
			r.GET("/api/v1/users/me", h.GetMe)

			req := httptest.NewRequest(
				"GET",
				"/api/v1/users/me",
				strings.NewReader(""))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestHandler_UpdateMe(t *testing.T) {
	log := zap.NewNop().Sugar()

	user := &domain.User{
		ID:       1,
		Login:    "alice",
		Fullname: "Alice",
	}

	tests := []struct {
		name       string
		token      string
		body       string
		setupMock  func(mid *MockAuthMiddleware, srv *MockUserService)
		wantStatus int
	}{
		{
			name:  "success",
			token: "1",
			body: `{
				"full_name": "Alice White"
			}`,
			setupMock: func(mid *MockAuthMiddleware, srv *MockUserService) {
				mid.EXPECT().UserFromContext(
					mock.Anything,
				).Return(&auth.TokenPayload{
					ID:    1,
					Login: "alice",
				}, true)
				fullname := "Alice White"
				updUser := user
				updUser.Fullname = fullname
				srv.EXPECT().UpdateByID(
					mock.Anything,
					&domain.PatchUser{
						ID:       1,
						Fullname: &fullname,
					},
				).Return(updUser, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "invalid token",
			token: "invalid",
			body: `{
				"full_name": "Alice White"
			}`,
			setupMock: func(mid *MockAuthMiddleware, srv *MockUserService) {
				mid.EXPECT().UserFromContext(
					mock.Anything,
				).Return(&auth.TokenPayload{}, false)
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:  "bad json",
			token: "1",
			body: `{
				"full_name
			}`,
			setupMock: func(mid *MockAuthMiddleware, srv *MockUserService) {
				mid.EXPECT().UserFromContext(
					mock.Anything,
				).Return(&auth.TokenPayload{
					ID:    1,
					Login: "alice",
				}, true)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "service full_name required error",
			token: "1",
			body: `{
				"full_name": ""
			}`,
			setupMock: func(mid *MockAuthMiddleware, srv *MockUserService) {
				mid.EXPECT().UserFromContext(
					mock.Anything,
				).Return(&auth.TokenPayload{
					ID:    1,
					Login: "alice",
				}, true)
				fullname := ""
				srv.EXPECT().UpdateByID(
					mock.Anything,
					&domain.PatchUser{
						ID:       1,
						Fullname: &fullname,
					},
				).Return(nil, service.FullnameRequired)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "service invalid full_name error",
			token: "1",
			body: `{
				"full_name": "Alice 123 White"
			}`,
			setupMock: func(mid *MockAuthMiddleware, srv *MockUserService) {
				mid.EXPECT().UserFromContext(
					mock.Anything,
				).Return(&auth.TokenPayload{
					ID:    1,
					Login: "alice",
				}, true)
				fullname := "Alice 123 White"
				srv.EXPECT().UpdateByID(
					mock.Anything,
					&domain.PatchUser{
						ID:       1,
						Fullname: &fullname,
					},
				).Return(nil, service.InvalidFullname)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "service currency required error",
			token: "1",
			body: `{
				"currency": ""
			}`,
			setupMock: func(mid *MockAuthMiddleware, srv *MockUserService) {
				mid.EXPECT().UserFromContext(
					mock.Anything,
				).Return(&auth.TokenPayload{
					ID:    1,
					Login: "alice",
				}, true)
				currency := ""
				srv.EXPECT().UpdateByID(
					mock.Anything,
					&domain.PatchUser{
						ID:       1,
						Currency: &currency,
					},
				).Return(nil, service.CurrencyRequired)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "service invalid currency error",
			token: "1",
			body: `{
				"currency": "dollar 123"
			}`,
			setupMock: func(mid *MockAuthMiddleware, srv *MockUserService) {
				mid.EXPECT().UserFromContext(
					mock.Anything,
				).Return(&auth.TokenPayload{
					ID:    1,
					Login: "alice",
				}, true)
				currency := "dollar 123"
				srv.EXPECT().UpdateByID(
					mock.Anything,
					&domain.PatchUser{
						ID:       1,
						Currency: &currency,
					},
				).Return(nil, service.InvalidCurrency)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "service internal error",
			token: "1",
			body: `{
				"full_name": "Alice White"
			}`,
			setupMock: func(mid *MockAuthMiddleware, srv *MockUserService) {
				mid.EXPECT().UserFromContext(
					mock.Anything,
				).Return(&auth.TokenPayload{
					ID:    1,
					Login: "alice",
				}, true)
				fullname := "Alice White"
				srv.EXPECT().UpdateByID(
					mock.Anything,
					&domain.PatchUser{
						ID:       1,
						Fullname: &fullname,
					},
				).Return(nil, fmt.Errorf("something went wrong"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewMockUserService(t)
			mid := NewMockAuthMiddleware(t)

			tt.setupMock(mid, srv)

			h := transport.NewHandler(srv, log, mid)

			r := gin.New()
			r.PATCH("/api/v1/users/me", h.UpdateMe)

			req := httptest.NewRequest(
				"PATCH",
				"/api/v1/users/me",
				strings.NewReader(tt.body))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestHandler_DeleteMe(t *testing.T) {
	log := zap.NewNop().Sugar()

	tests := []struct {
		name       string
		token      string
		setupMock  func(mid *MockAuthMiddleware, srv *MockUserService)
		wantStatus int
	}{
		{
			name:  "success",
			token: "1",
			setupMock: func(mid *MockAuthMiddleware, srv *MockUserService) {
				mid.EXPECT().UserFromContext(
					mock.Anything,
				).Return(&auth.TokenPayload{
					ID:    1,
					Login: "alice",
				}, true)
				srv.EXPECT().DeleteByID(
					mock.Anything,
					int64(1),
				).Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:  "invalid token",
			token: "invalid",
			setupMock: func(mid *MockAuthMiddleware, srv *MockUserService) {
				mid.EXPECT().UserFromContext(
					mock.Anything,
				).Return(&auth.TokenPayload{}, false)
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:  "service user not found error",
			token: "1",
			setupMock: func(mid *MockAuthMiddleware, srv *MockUserService) {
				mid.EXPECT().UserFromContext(
					mock.Anything,
				).Return(&auth.TokenPayload{
					ID:    1,
					Login: "alice",
				}, true)
				srv.EXPECT().DeleteByID(
					mock.Anything,
					int64(1),
				).Return(service.UserNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:  "service internal error",
			token: "1",
			setupMock: func(mid *MockAuthMiddleware, srv *MockUserService) {
				mid.EXPECT().UserFromContext(
					mock.Anything,
				).Return(&auth.TokenPayload{
					ID:    1,
					Login: "alice",
				}, true)
				srv.EXPECT().DeleteByID(
					mock.Anything,
					int64(1),
				).Return(fmt.Errorf("something went wrong"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewMockUserService(t)
			mid := NewMockAuthMiddleware(t)

			tt.setupMock(mid, srv)

			h := transport.NewHandler(srv, log, mid)

			r := gin.New()
			r.DELETE("/api/v1/users/me", h.DeleteMe)

			req := httptest.NewRequest(
				"DELETE",
				"/api/v1/users/me",
				strings.NewReader(""))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
