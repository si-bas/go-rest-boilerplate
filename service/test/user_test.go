package test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"github.com/si-bas/go-rest-boilerplate/domain/model"
	repoMocks "github.com/si-bas/go-rest-boilerplate/domain/repository/mocks"
	"github.com/si-bas/go-rest-boilerplate/service"
	"github.com/si-bas/go-rest-boilerplate/shared/helper/pagination"
	"gorm.io/gorm"
)

type userMock struct {
	userRepo repoMocks.UserRepository
}

func TestUserEmailIsUsed(t *testing.T) {
	testCases := []struct {
		name     string
		mockFunc func(mock *userMock)
		wantErr  error
	}{
		{
			name: "email is used",
			mockFunc: func(listMock *userMock) {
				countResult := int64(1)
				listMock.userRepo.On("CountByEmail", mock.Anything, mock.Anything).Return(&countResult, nil)
			},
		},
		{
			name: "email is not used",
			mockFunc: func(listMock *userMock) {
				countResult := int64(0)
				listMock.userRepo.On("CountByEmail", mock.Anything, mock.Anything).Return(&countResult, nil)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			listMock := userMock{
				userRepo: repoMocks.UserRepository{},
			}
			if tc.mockFunc != nil {
				tc.mockFunc(&listMock)
			}

			svc := service.NewUserService(&listMock.userRepo)
			result, err := svc.EmailIsUsed(context.TODO(), "newuser@mail.com")

			assert.Equal(t, tc.wantErr, err)
			listMock.userRepo.AssertExpectations(t)

			if err == nil {
				if result {
					assert.Equal(t, result, true)
				} else {
					assert.Equal(t, result, false)
				}
			}
		})
	}
}

func TestUserCreate(t *testing.T) {
	newUser := model.CreateUser{
		Name:     "new user",
		Email:    "newuser@mail.com",
		Password: "secret",
	}

	testCases := []struct {
		name     string
		mockFunc func(mock *userMock)
		wantErr  error
	}{
		{
			name: "success create user",
			mockFunc: func(listMock *userMock) {
				countResult := int64(0)
				listMock.userRepo.On("CountByEmail", mock.Anything, mock.Anything).Return(&countResult, nil)
				listMock.userRepo.On("Insert", mock.Anything, &model.User{
					Name:     newUser.Name,
					Email:    newUser.Email,
					Password: newUser.Password,
				}).Return(nil)
			},
		},
		{
			name: "failed create user - email already exists",
			mockFunc: func(listMock *userMock) {
				countResult := int64(1)
				listMock.userRepo.On("CountByEmail", mock.Anything, mock.Anything).Return(&countResult, nil)
			},
			wantErr: errors.New("email already used"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			listMock := userMock{
				userRepo: repoMocks.UserRepository{},
			}
			if tc.mockFunc != nil {
				tc.mockFunc(&listMock)
			}

			svc := service.NewUserService(&listMock.userRepo)
			result, err := svc.Create(context.TODO(), newUser)

			assert.Equal(t, tc.wantErr, err)
			listMock.userRepo.AssertExpectations(t)

			if err == nil {
				assert.Equal(t, result, model.User{
					Name:     newUser.Name,
					Email:    newUser.Email,
					Password: newUser.Password,
				})
			}
		})
	}
}

func TestUserDetail(t *testing.T) {
	existingUser := model.User{
		ID:    1,
		Name:  "user",
		Email: "user@mail.com",
	}

	testCases := []struct {
		name     string
		mockFunc func(mock *userMock)
		wantErr  error
	}{
		{
			name: "success get user detail",
			mockFunc: func(listMock *userMock) {
				listMock.userRepo.On("FindById", mock.Anything, uint32(1)).Return(&existingUser, nil)
			},
		},
		{
			name: "failed get user detail - not found",
			mockFunc: func(listMock *userMock) {
				listMock.userRepo.On("FindById", mock.Anything, uint32(1)).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: gorm.ErrRecordNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			listMock := userMock{
				userRepo: repoMocks.UserRepository{},
			}
			if tc.mockFunc != nil {
				tc.mockFunc(&listMock)
			}

			svc := service.NewUserService(&listMock.userRepo)
			result, err := svc.Detail(context.TODO(), uint32(1))

			assert.Equal(t, tc.wantErr, err)
			listMock.userRepo.AssertExpectations(t)

			if err == nil {
				assert.Equal(t, result, &existingUser)
			}
		})
	}
}

func TestUserListPaginate(t *testing.T) {
	users := []model.User{
		{
			ID:    1,
			Name:  "user",
			Email: "user@mail.com",
		},
		{
			ID:    2,
			Name:  "another user",
			Email: "auser@mail.com",
		},
	}

	testCases := []struct {
		name     string
		mockFunc func(mock *userMock)
		wantErr  error
	}{
		{
			name: "success get list of users",
			mockFunc: func(listMock *userMock) {
				listMock.userRepo.On("GetPaginate", mock.Anything, model.UserFilter{}, pagination.Param{}).Return(users, &pagination.Param{}, nil)
			},
		},
		{
			name: "error get list of users - no users",
			mockFunc: func(listMock *userMock) {
				listMock.userRepo.On("GetPaginate", mock.Anything, model.UserFilter{}, pagination.Param{}).Return(nil, nil, gorm.ErrRecordNotFound)
			},
			wantErr: gorm.ErrRecordNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			listMock := userMock{
				userRepo: repoMocks.UserRepository{},
			}
			if tc.mockFunc != nil {
				tc.mockFunc(&listMock)
			}

			svc := service.NewUserService(&listMock.userRepo)
			result, _, err := svc.ListPaginate(context.TODO(), model.UserFilter{}, pagination.Param{})

			assert.Equal(t, tc.wantErr, err)
			listMock.userRepo.AssertExpectations(t)

			if err == nil {
				assert.Equal(t, result, users)
			}
		})
	}
}
