package repository

import (
	"context"

	"github.com/si-bas/go-rest-boilerplate/domain/model"
	"github.com/si-bas/go-rest-boilerplate/shared/helper/pagination"
	"gorm.io/gorm"
)

type UserRepository interface {
	FilteredDb(model.UserFilter) *gorm.DB

	Insert(context.Context, *model.User) error
	GetPaginate(context.Context, model.UserFilter, pagination.Param) ([]model.User, *pagination.Param, error)
	GetFiltered(context.Context, model.UserFilter) ([]model.User, error)
	CountByEmail(context.Context, string) (*int64, error)
	FindById(context.Context, uint32) (*model.User, error)
	FindByEmail(context.Context, string) (*model.User, error)
}

type userImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userImpl{
		db: db,
	}
}

func (r *userImpl) Insert(ctx context.Context, user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userImpl) GetFiltered(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	var users []model.User
	err := r.FilteredDb(filter).Find(&users).Error

	return users, err
}

func (r *userImpl) GetPaginate(ctx context.Context, filter model.UserFilter, param pagination.Param) ([]model.User, *pagination.Param, error) {
	var users []model.User

	filteredDb := r.FilteredDb(filter)

	if err := filteredDb.Scopes(pagination.Paginate(model.User{}, &param, filteredDb)).Find(&users).Error; err != nil {
		return nil, nil, err
	}

	return users, &param, nil
}

func (r *userImpl) CountByEmail(ctx context.Context, email string) (*int64, error) {
	var count int64
	if err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *userImpl) FilteredDb(filter model.UserFilter) *gorm.DB {
	chain := r.db.Model(&model.User{})

	if filter.Keyword != "" {
		searchVal := "%" + filter.Keyword + "%"
		chain.Where("name LIKE ? OR email LIKE ?", searchVal, searchVal)
	}

	return chain
}

func (r *userImpl) FindById(ctx context.Context, id uint32) (*model.User, error) {
	var user model.User
	if err := r.db.Model(&model.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
func (r *userImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.Model(&model.User{}).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
