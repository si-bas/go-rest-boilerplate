package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	password := []byte(u.Password)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return
}

func (u *User) VerifyPassword(password string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
	}
	return
}

type UserFind struct {
	ID uint32 `uri:"id" binding:"required"`
}

type UserFilter struct {
	Keyword string `query:"q" form:"q" url:"q" json:"keyword"`
	Name    string `query:"name" form:"name" url:"name" json:"name"`
	Email   string `query:"email" form:"email" url:"email" json:"email"`
}

type UserListRequest struct {
	UserFilter
	Limit uint              `query:"limit,omitempty" form:"limit"`
	Page  uint              `query:"page,omitempty" form:"page"`
	Sort  map[string]string `query:"sort,omitempty" form:"sort"`
}

type CreateUser struct {
	Name     string
	Email    string
	Password string
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
}

type CreateUserResponse struct {
	Id        uint32    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
