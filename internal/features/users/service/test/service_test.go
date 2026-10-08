package service_test

import (
	"fmt"
	"testing"

	"github.com/SJMaks/financer_back/internal/core/auth"
	"github.com/SJMaks/financer_back/internal/features/users/domain"
	"github.com/SJMaks/financer_back/internal/features/users/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestService_Register(t *testing.T) {
	tm := NewMockTokenManager(t)

	tests := []struct {
		name      string
		setupMock func(repo *MockUserRepository, u *domain.User)
		u         *domain.User
		password  string
		want      *domain.User
		wantErr   error
	}{
		{
			name: "success",
			setupMock: func(repo *MockUserRepository, u *domain.User) {
				user := *u
				user.ID = 1
				user.Password = "hash"
				repo.EXPECT().Create(
					mock.Anything,
					u,
				).Return(&user, nil)
			},
			u: &domain.User{
				Login:    "alice",
				Fullname: "Alice",
			},
			password: "SECret123",
			want: &domain.User{
				ID:       1,
				Login:    "alice",
				Password: "hash",
				Fullname: "Alice",
			},
			wantErr: nil,
		},
		{
			name:      "login required error",
			setupMock: func(repo *MockUserRepository, u *domain.User) {},
			u: &domain.User{
				Fullname: "Alice",
			},
			password: "SECret123",
			want:     nil,
			wantErr:  service.LoginRequired,
		},
		{
			name:      "password required error",
			setupMock: func(repo *MockUserRepository, u *domain.User) {},
			u: &domain.User{
				Login:    "alice",
				Fullname: "Alice",
			},
			password: "",
			want:     nil,
			wantErr:  service.PasswordRequired,
		},
		{
			name:      "password required error",
			setupMock: func(repo *MockUserRepository, u *domain.User) {},
			u: &domain.User{
				Login: "alice",
			},
			password: "SECret123",
			want:     nil,
			wantErr:  service.FullnameRequired,
		},
		{
			name:      "invalid login error",
			setupMock: func(repo *MockUserRepository, u *domain.User) {},
			u: &domain.User{
				Login:    "alice and me",
				Fullname: "Alice",
			},
			password: "SECret123",
			want:     nil,
			wantErr:  service.InvalidLogin,
		},
		{
			name:      "invalid password error",
			setupMock: func(repo *MockUserRepository, u *domain.User) {},
			u: &domain.User{
				Login:    "alice",
				Fullname: "Alice",
			},
			password: "secret123",
			want:     nil,
			wantErr:  service.InvalidPassword,
		},
		{
			name:      "invalid full_name error",
			setupMock: func(repo *MockUserRepository, u *domain.User) {},
			u: &domain.User{
				Login:    "alice",
				Fullname: "Alice123",
			},
			password: "SECret123",
			want:     nil,
			wantErr:  service.InvalidFullname,
		},
		{
			name: "repo user already exists error",
			setupMock: func(repo *MockUserRepository, u *domain.User) {
				repo.EXPECT().Create(
					mock.Anything,
					u,
				).Return(nil, service.UserAlreadyExists)
			},
			u: &domain.User{
				Login:    "alice",
				Fullname: "Alice",
			},
			password: "SECret123",
			want:     nil,
			wantErr:  service.UserAlreadyExists,
		},
		{
			name: "repo internal error",
			setupMock: func(repo *MockUserRepository, u *domain.User) {
				repo.EXPECT().Create(
					mock.Anything,
					u,
				).Return(nil, fmt.Errorf("internal"))
			},
			u: &domain.User{
				Login:    "alice",
				Fullname: "Alice",
			},
			password: "SECret123",
			want:     nil,
			wantErr:  fmt.Errorf("internal"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUserRepository(t)
			tt.setupMock(repo, tt.u)

			s := service.NewService(repo, tm)

			got, gotErr := s.Register(t.Context(), tt.u, tt.password)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

func TestService_Login(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(tm *MockTokenManager, repo *MockUserRepository, u *domain.User)
		u         *domain.User
		want      string
		wantErr   error
	}{
		{
			name: "success",
			setupMock: func(tm *MockTokenManager, repo *MockUserRepository, u *domain.User) {
				hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
				require.NoError(t, err)

				user := *u
				user.ID = 1
				user.Password = string(hash)
				repo.EXPECT().GetByLogin(
					mock.Anything,
					u.Login,
				).Return(&user, nil)
				tm.EXPECT().Generate(auth.TokenPayload{
					ID:    user.ID,
					Login: user.Login,
				}).Return("token", nil)
			},
			u: &domain.User{
				Login:    "alice",
				Password: "SECret123",
			},
			want:    "token",
			wantErr: nil,
		},
		{
			name:      "login required error",
			setupMock: func(tm *MockTokenManager, repo *MockUserRepository, u *domain.User) {},
			u:         &domain.User{},
			want:      "",
			wantErr:   service.LoginRequired,
		},
		{
			name:      "password required error",
			setupMock: func(tm *MockTokenManager, repo *MockUserRepository, u *domain.User) {},
			u: &domain.User{
				Login: "alice",
			},
			want:    "",
			wantErr: service.PasswordRequired,
		},
		{
			name: "invalid credentials login error",
			setupMock: func(tm *MockTokenManager, repo *MockUserRepository, u *domain.User) {
				repo.EXPECT().GetByLogin(
					mock.Anything,
					u.Login,
				).Return(nil, service.UserNotFound)
			},
			u: &domain.User{
				Login:    "aliceson",
				Password: "SECret123",
			},
			want:    "",
			wantErr: service.InvalidCredentials,
		},
		{
			name: "invalid credentials password error",
			setupMock: func(tm *MockTokenManager, repo *MockUserRepository, u *domain.User) {
				user := *u
				user.ID = 1
				user.Password = "hash"
				repo.EXPECT().GetByLogin(
					mock.Anything,
					u.Login,
				).Return(&user, nil)
			},
			u: &domain.User{
				Login:    "alice",
				Password: "SECret",
			},
			want:    "",
			wantErr: service.InvalidCredentials,
		},
		{
			name: "repo internal error",
			setupMock: func(tm *MockTokenManager, repo *MockUserRepository, u *domain.User) {
				repo.EXPECT().GetByLogin(
					mock.Anything,
					u.Login,
				).Return(nil, fmt.Errorf("internal"))
			},
			u: &domain.User{
				Login:    "alice",
				Password: "SECret123",
			},
			want:    "",
			wantErr: fmt.Errorf("internal"),
		},
		{
			name: "tm internal error",
			setupMock: func(tm *MockTokenManager, repo *MockUserRepository, u *domain.User) {
				hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
				require.NoError(t, err)

				user := *u
				user.ID = 1
				user.Password = string(hash)
				repo.EXPECT().GetByLogin(
					mock.Anything,
					u.Login,
				).Return(&user, nil)
				tm.EXPECT().Generate(auth.TokenPayload{
					ID:    user.ID,
					Login: user.Login,
				}).Return("", fmt.Errorf("internal"))
			},
			u: &domain.User{
				Login:    "alice",
				Password: "SECret123",
			},
			want:    "",
			wantErr: fmt.Errorf("internal"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUserRepository(t)
			tm := NewMockTokenManager(t)

			tt.setupMock(tm, repo, tt.u)

			s := service.NewService(repo, tm)

			got, gotErr := s.Login(t.Context(), tt.u)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

func TestService_GetByID(t *testing.T) {
	tm := NewMockTokenManager(t)

	tests := []struct {
		name      string
		setupMock func(repo *MockUserRepository, id int64, user *domain.User)
		id        int64
		want      *domain.User
		wantErr   error
	}{
		{
			name: "success",
			setupMock: func(repo *MockUserRepository, id int64, user *domain.User) {
				repo.EXPECT().GetByID(
					mock.Anything,
					id,
				).Return(user, nil)
			},
			id: 1,
			want: &domain.User{
				ID:       1,
				Login:    "alice",
				Password: "hash",
				Fullname: "Alice",
			},
			wantErr: nil,
		},
		{
			name: "repo user not found error",
			setupMock: func(repo *MockUserRepository, id int64, user *domain.User) {
				repo.EXPECT().GetByID(
					mock.Anything,
					id,
				).Return(nil, service.UserNotFound)
			},
			id:      1,
			want:    nil,
			wantErr: service.UserNotFound,
		},
		{
			name: "repo internal error",
			setupMock: func(repo *MockUserRepository, id int64, user *domain.User) {
				repo.EXPECT().GetByID(
					mock.Anything,
					id,
				).Return(nil, fmt.Errorf("internal"))
			},
			id:      1,
			want:    nil,
			wantErr: fmt.Errorf("internal"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUserRepository(t)

			tt.setupMock(repo, tt.id, tt.want)

			s := service.NewService(repo, tm)

			got, gotErr := s.GetByID(t.Context(), tt.id)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

func TestService_UpdateByID(t *testing.T) {
	tm := NewMockTokenManager(t)

	newFullname := "Alice White"
	emptyFullname := ""
	invalidFullname := "Alice123"

	tests := []struct {
		name      string
		setupMock func(repo *MockUserRepository, u *domain.PatchUser, user *domain.User)
		u         *domain.PatchUser
		want      *domain.User
		wantErr   error
	}{
		{
			name: "success",
			setupMock: func(repo *MockUserRepository, u *domain.PatchUser, user *domain.User) {
				repo.EXPECT().UpdateByID(
					mock.Anything,
					u,
				).Return(user, nil)
			},
			u: &domain.PatchUser{
				ID:       1,
				Fullname: &newFullname,
			},
			want: &domain.User{
				ID:       1,
				Login:    "alice",
				Password: "hash",
				Fullname: "Alice White",
			},
			wantErr: nil,
		},
		{
			name:      "full_name required error",
			setupMock: func(repo *MockUserRepository, u *domain.PatchUser, user *domain.User) {},
			u: &domain.PatchUser{
				ID:       1,
				Fullname: &emptyFullname,
			},
			want:    nil,
			wantErr: service.FullnameRequired,
		},
		{
			name:      "invalid full_name",
			setupMock: func(repo *MockUserRepository, u *domain.PatchUser, user *domain.User) {},
			u: &domain.PatchUser{
				ID:       1,
				Fullname: &invalidFullname,
			},
			want:    nil,
			wantErr: service.InvalidFullname,
		},
		{
			name: "repo internal error",
			setupMock: func(repo *MockUserRepository, u *domain.PatchUser, user *domain.User) {
				repo.EXPECT().UpdateByID(
					mock.Anything,
					u,
				).Return(nil, fmt.Errorf("internal"))
			},
			u: &domain.PatchUser{
				ID:       1,
				Fullname: &newFullname,
			},
			want:    nil,
			wantErr: fmt.Errorf("internal"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUserRepository(t)

			tt.setupMock(repo, tt.u, tt.want)

			s := service.NewService(repo, tm)

			got, gotErr := s.UpdateByID(t.Context(), tt.u)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

func TestService_DeleteByID(t *testing.T) {
	tm := NewMockTokenManager(t)

	tests := []struct {
		name      string
		setupMock func(repo *MockUserRepository, id int64)
		id        int64
		wantErr   error
	}{
		{
			name: "success",
			setupMock: func(repo *MockUserRepository, id int64) {
				repo.EXPECT().DeleteByID(
					mock.Anything,
					id,
				).Return(nil)
			},
			id:      1,
			wantErr: nil,
		},
		{
			name: "repo user not found error",
			setupMock: func(repo *MockUserRepository, id int64) {
				repo.EXPECT().DeleteByID(
					mock.Anything,
					id,
				).Return(service.UserNotFound)
			},
			id:      1,
			wantErr: service.UserNotFound,
		},
		{
			name: "repo internal error",
			setupMock: func(repo *MockUserRepository, id int64) {
				repo.EXPECT().DeleteByID(
					mock.Anything,
					id,
				).Return(fmt.Errorf("internal"))
			},
			id:      1,
			wantErr: fmt.Errorf("internal"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUserRepository(t)

			tt.setupMock(repo, tt.id)

			s := service.NewService(repo, tm)

			gotErr := s.DeleteByID(t.Context(), tt.id)

			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
