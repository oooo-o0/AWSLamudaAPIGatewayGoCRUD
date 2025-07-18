package service

import (
	"context"
	"fmt"
	"time"

	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
	"github.com/oooo-o0/kaigo-insurance-system/batch/interface/repository"
)

type UserServiceInterface interface {
	GetUsersInactiveSince(ctx context.Context, cutoff time.Time) ([]*model.User, error)
	DeleteUser(ctx context.Context, userID string) error
}

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUsersInactiveSince 指定日時よりログイン履歴が古いユーザー一覧を取得
func (s *UserService) GetUsersInactiveSince(ctx context.Context, cutoff time.Time) ([]*model.User, error) {
	// DynamoDBでLastLoginAt <= cutoff の条件で検索（GSI等を使う想定）
	users, err := s.userRepo.QueryUsersByLastLoginBefore(ctx, cutoff)
	if err != nil {
		return nil, fmt.Errorf("failed to query inactive users: %w", err)
	}
	return users, nil
}

// DeleteUser ユーザー削除（論理削除 or 物理削除は要件次第）
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	// 物理削除例
	err := s.userRepo.Delete(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user %s: %w", userID, err)
	}
	return nil
}
