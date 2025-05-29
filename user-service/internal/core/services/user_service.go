package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"user-service/internal/core/domain"
	"user-service/internal/core/utils"
)

type UserService interface {
	Register(ctx context.Context, email, password, name string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, id, email, name string) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	VerifyPassword(hashedPassword, plainPassword string) bool
}

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Register(ctx context.Context, email, password, name string) (*domain.User, error) {
	existingUser, err := s.repo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		ID:       uuid.New().String(),
		Email:    email,
		Name:     name,
		Password: hashedPassword,
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id, email, name string) (*domain.User, error) {
	existingUser, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	if email != existingUser.Email {
		userWithEmail, err := s.repo.FindByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		if userWithEmail != nil && userWithEmail.ID != id {
			return nil, errors.New("email already taken")
		}
	}

	existingUser.Email = email
	existingUser.Name = name

	if err := s.repo.Update(ctx, existingUser); err != nil {
		return nil, err
	}

	return existingUser, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	existingUser, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	return s.repo.Delete(ctx, id)
}

func (s *userService) VerifyPassword(hashedPassword, plainPassword string) bool {
	return utils.VerifyPassword(hashedPassword, plainPassword)
}
