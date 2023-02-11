package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/si-bas/go-rest-boilerplate/config"
	"github.com/si-bas/go-rest-boilerplate/domain/model"
	"github.com/si-bas/go-rest-boilerplate/domain/repository"
)

type AuthService interface {
	ValidateUser(context.Context, model.ValidateUser) (*model.User, error)
	GenerateToken(context.Context, *model.User) (*model.JwtToken, error)
	ParseToken(context.Context, string) (*jwt.Token, error)
	GetClaims(context.Context, *jwt.Token) (jwt.MapClaims, error)
	GetUser(context.Context, uint32) (*model.User, error)
}

type authImpl struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authImpl{
		userRepo: userRepo,
	}
}

func (s *authImpl) ValidateUser(ctx context.Context, prerequisite model.ValidateUser) (*model.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, prerequisite.Email)
	if err != nil {
		return nil, err
	}

	if user.VerifyPassword(prerequisite.Password) != nil {
		return nil, errors.New("password is invalid")
	}
	return user, nil
}

func (s *authImpl) GenerateToken(ctx context.Context, user *model.User) (*model.JwtToken, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.ID
	claims["name"] = user.Name
	claims["exp"] = time.Now().Add(time.Second * time.Duration(rand.Int31n(config.Config.Jwt.ExpiresIn))).Unix()

	t, err := token.SignedString([]byte(config.Config.Jwt.Secret))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = user.ID
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	rtClaims["exp"] = time.Now().Add(time.Second * time.Duration(rand.Int31n(config.Config.Jwt.RefreshExpiresIn))).Unix()

	rt, err := refreshToken.SignedString([]byte(config.Config.Jwt.Secret))
	if err != nil {
		return nil, err
	}

	return &model.JwtToken{
		AccessToken:      t,
		ExpiresIn:        config.Config.Jwt.ExpiresIn,
		RefreshToken:     rt,
		RefreshExpiresIn: config.Config.Jwt.RefreshExpiresIn,
	}, nil
}

func (s *authImpl) ParseToken(ctx context.Context, accessToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}

		return []byte(config.Config.Jwt.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *authImpl) GetClaims(ctx context.Context, token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("error when get token claims")
	}

	return claims, nil
}

func (s *authImpl) GetUser(ctx context.Context, id uint32) (*model.User, error) {
	user, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
