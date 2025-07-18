package user

import (
	"context"
	"fmt"
	"time"

	"github.com/oooo-o0/kaigo-insurance-system/batch/interface/job_handler"
	"github.com/oooo-o0/kaigo-insurance-system/batch/service"
)

const (
	// 非アクティブ期間（例：180日）
	inactiveDays = 180
)

type InactiveUserCleanupJob struct {
	userService service.UserServiceInterface
}

func NewInactiveUserCleanupJob(userSvc service.UserServiceInterface) *InactiveUserCleanupJob {
	return &InactiveUserCleanupJob{
		userService: userSvc,
	}
}

// Execute はLambdaの起動ポイント
func (job *InactiveUserCleanupJob) Execute(ctx context.Context, payload []byte) error {
	fmt.Println("InactiveUserCleanupJob started")

	cutoffDate := time.Now().AddDate(0, 0, -inactiveDays)
	// 非アクティブユーザーを取得
	users, err := job.userService.GetUsersInactiveSince(ctx, cutoffDate)
	if err != nil {
		return fmt.Errorf("failed to get inactive users: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("No inactive users found")
		return nil
	}

	fmt.Printf("Found %d inactive users. Deleting...\n", len(users))

	for _, user := range users {
		err := job.userService.DeleteUser(ctx, user.UserID)
		if err != nil {
			// ログだけ残し続行
			fmt.Printf("Failed to delete user %s: %v\n", user.UserID, err)
		} else {
			fmt.Printf("Deleted user %s\n", user.UserID)
		}
	}

	fmt.Println("InactiveUserCleanupJob finished")
	return nil
}

var _ job_handler.JobHandler = (*InactiveUserCleanupJob)(nil)
