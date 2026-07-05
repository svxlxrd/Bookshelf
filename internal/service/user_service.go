package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/bookshelf/monolith/internal/domain"
	"github.com/bookshelf/monolith/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUsernameExists     = errors.New("username already exists")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidUsername    = errors.New("invalid username")
	ErrInvalidEmail       = errors.New("invalid email")
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func (s *UserService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	// валидация
	if len(req.Username) < 3 {
		return nil, ErrInvalidUsername
	}

	if len(req.Password) < 8 {
		return nil, ErrInvalidPassword
	}

	if req.Email == "" {
		return nil, ErrInvalidEmail
	}

	// проверка уникальности
	emailExists, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, ErrUserExists
	}

	usernameExists, err := s.repo.UsernameExists(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, ErrUsernameExists
	}

	// хеширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate hash: %w", err)
	}

	// создание пользователя
	user := &domain.User{
		Username: req.Username,
		Email: req.Email,
		PasswordHash: string(hash),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// генерируем токен и возвращаем ответ
	return &domain.AuthResponse{
			User: domain.UserPublic{
				Username: user.Username,
				Email:    user.Email,
			},
		}, nil
}