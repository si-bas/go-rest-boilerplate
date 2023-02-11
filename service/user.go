package service

import (
	"context"
	"errors"

	"github.com/si-bas/go-rest-boilerplate/domain/model"
	"github.com/si-bas/go-rest-boilerplate/domain/repository"
	"github.com/si-bas/go-rest-boilerplate/shared/helper/pagination"
)

type UserService interface {
	Create(context.Context, model.CreateUser) (*model.User, error)
	EmailIsUsed(context.Context, string) (bool, error)
	ListPaginate(context.Context, model.UserFilter, pagination.Param) ([]model.User, *pagination.Param, error)
	Detail(context.Context, uint32) (*model.User, error)
}

type userImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userImpl{
		userRepo: userRepo,
	}
}

func (s *userImpl) Create(ctx context.Context, payload model.CreateUser) (*model.User, error) {
	if emailIsUsed, err := s.EmailIsUsed(ctx, payload.Email); emailIsUsed || err != nil {
		if err != nil {
			return nil, err
		}

		return nil, errors.New("email already used")
	}

	newUser := model.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	}
	if err := s.userRepo.Insert(ctx, &newUser); err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (s *userImpl) EmailIsUsed(ctx context.Context, email string) (bool, error) {
	count, err := s.userRepo.CountByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	return *count > 0, nil
}

func (s *userImpl) ListPaginate(ctx context.Context, filter model.UserFilter, query pagination.Param) ([]model.User, *pagination.Param, error) {
	users, meta, err := s.userRepo.GetPaginate(ctx, filter, query)
	if err != nil {
		return nil, nil, err
	}

	return users, meta, nil
}

func (s *userImpl) Detail(ctx context.Context, id uint32) (*model.User, error) {
	user, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
