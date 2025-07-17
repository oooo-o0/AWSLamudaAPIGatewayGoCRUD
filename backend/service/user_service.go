package service

import (
	"context"
	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/backend/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id string) error
	GetUsersByRole(ctx context.Context, role string) ([]*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) error {
	return s.repo.Put(ctx, user)
}

func (s *userService) GetUser(ctx context.Context, id string) (*model.User, error) {
	return s.repo.Get(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) GetUsersByRole(ctx context.Context, role string) ([]*model.User, error) {
	return s.repo.ListByRole(ctx, role)
}

func (s *userService) UpdateUserRole(ctx context.Context, userID, role string) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Role = role
	return s.repo.Update(ctx, user)
}
