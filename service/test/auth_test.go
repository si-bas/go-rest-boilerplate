package test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/mock"
	"github.com/si-bas/go-rest-boilerplate/config"
	"github.com/si-bas/go-rest-boilerplate/domain/model"
	repoMocks "github.com/si-bas/go-rest-boilerplate/domain/repository/mocks"
	"github.com/si-bas/go-rest-boilerplate/service"
	"gorm.io/gorm"
)

type authMock struct {
	userRepo repoMocks.UserRepository
}

func TestAuthValidateUser(t *testing.T) {
	user := model.User{
		ID:       1,
		Email:    "newuser@mail.com",
		Password: "$2a$10$6ItIWM3fUWVmY1GzGU4pzOGUtUVXUbbkHVA1F9fEvlkJchvHU7XF2",
	}

	testCases := []struct {
		name     string
		mockFunc func(mock *authMock)
		wantErr  error
	}{
		{
			name: "email and password is valid",
			mockFunc: func(listMock *authMock) {
				listMock.userRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(&user, nil)
			},
		},
		{
			name: "email is invalid",
			mockFunc: func(listMock *authMock) {
				listMock.userRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "password is invalid",
			mockFunc: func(listMock *authMock) {
				user.Password = "xxx"
				listMock.userRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(&user, nil)
			},
			wantErr: errors.New("password is invalid"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			listMock := authMock{
				userRepo: repoMocks.UserRepository{},
			}
			if tc.mockFunc != nil {
				tc.mockFunc(&listMock)
			}

			svc := service.NewAuthService(&listMock.userRepo)
			result, err := svc.ValidateUser(context.TODO(), model.ValidateUser{
				Email:    user.Email,
				Password: "admin",
			})
			assert.Equal(t, tc.wantErr, err)
			listMock.userRepo.AssertExpectations(t)

			if err == nil {
				assert.Equal(t, result, &user)
			}
		})
	}
}

func TestAuthGenerateToken(t *testing.T) {
	config.Config = &config.Cfg{}
	config.Config.Jwt.Secret = "test-secret"
	config.Config.Jwt.ExpiresIn = 1800
	config.Config.Jwt.RefreshExpiresIn = 3600

	var user model.User

	testCases := []struct {
		name     string
		mockFunc func(mock *authMock)
		wantErr  error
	}{
		{
			name: "success generate token",
			mockFunc: func(listMock *authMock) {
				user = model.User{
					ID:    1,
					Name:  "user",
					Email: "newuser@mail.com",
				}
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			listMock := authMock{
				userRepo: repoMocks.UserRepository{},
			}
			if tc.mockFunc != nil {
				tc.mockFunc(&listMock)
			}

			svc := service.NewAuthService(&listMock.userRepo)
			result, err := svc.GenerateToken(context.TODO(), &user)
			assert.Equal(t, tc.wantErr, err)
			listMock.userRepo.AssertExpectations(t)

			if err == nil {
				assert.IsEqual(result, &model.JwtToken{})
			}
		})
	}
}

func TestAuthParseToken(t *testing.T) {
	config.Config = &config.Cfg{}
	config.Config.Jwt.Secret = "test-secret"

	var accessToken string

	testCases := []struct {
		name     string
		mockFunc func(mock *authMock)
		wantErr  error
	}{
		{
			name: "success parse token",
			mockFunc: func(listMock *authMock) {
				accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzYxMDIyMTYsIm5hbWUiOiJ1c2VyIiwic3ViIjoxfQ.GyvQnVavFf4M4aHyKKlTAs2kjYsWTTituUDV4lVfvDI"
			},
		},
		{
			name: "failed parse token",
			mockFunc: func(listMock *authMock) {
				accessToken = "xxxx"
			},
			wantErr: jwt.ValidationError{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			listMock := authMock{
				userRepo: repoMocks.UserRepository{},
			}
			if tc.mockFunc != nil {
				tc.mockFunc(&listMock)
			}

			svc := service.NewAuthService(&listMock.userRepo)
			result, err := svc.ParseToken(context.TODO(), accessToken)

			assert.IsEqual(tc.wantErr, err)
			listMock.userRepo.AssertExpectations(t)

			if err == nil {
				assert.IsEqual(result, &jwt.Token{})
			}
		})
	}
}

func TestAuthGetClaims(t *testing.T) {
	config.Config = &config.Cfg{}
	config.Config.Jwt.Secret = "test-secret"

	var jwtToken jwt.Token

	testCases := []struct {
		name     string
		mockFunc func(mock *authMock)
		wantErr  error
	}{
		{
			name: "success get claims",
			mockFunc: func(listMock *authMock) {
				jwtToken = jwt.Token{
					Claims: jwt.MapClaims{},
				}
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			listMock := authMock{
				userRepo: repoMocks.UserRepository{},
			}
			if tc.mockFunc != nil {
				tc.mockFunc(&listMock)
			}

			svc := service.NewAuthService(&listMock.userRepo)
			result, err := svc.GetClaims(context.TODO(), &jwtToken)

			assert.IsEqual(tc.wantErr, err)
			listMock.userRepo.AssertExpectations(t)

			if err == nil {
				assert.IsEqual(result, jwt.MapClaims{})
			}
		})
	}
}

func TestAuthGetUser(t *testing.T) {
	user := model.User{
		ID:    1,
		Name:  "user",
		Email: "newuser@mail.com",
	}

	testCases := []struct {
		name     string
		mockFunc func(mock *authMock)
		wantErr  error
	}{
		{
			name: "success get user",
			mockFunc: func(listMock *authMock) {
				listMock.userRepo.On("FindById", mock.Anything, mock.Anything).Return(&user, nil)
			},
		},
		{
			name: "failed get user",
			mockFunc: func(listMock *authMock) {
				listMock.userRepo.On("FindById", mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: gorm.ErrRecordNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			listMock := authMock{
				userRepo: repoMocks.UserRepository{},
			}
			if tc.mockFunc != nil {
				tc.mockFunc(&listMock)
			}

			svc := service.NewAuthService(&listMock.userRepo)
			result, err := svc.GetUser(context.TODO(), user.ID)

			assert.Equal(t, tc.wantErr, err)
			listMock.userRepo.AssertExpectations(t)

			if err == nil {
				assert.Equal(t, result, &user)
			}
		})
	}

}
