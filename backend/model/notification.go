package model

import (
	"time"
)

// NotificationStatus 通知ステータス定義
type NotificationStatus string

const (
	StatusUnread  NotificationStatus = "UNREAD"
	StatusRead    NotificationStatus = "READ"
	StatusDeleted NotificationStatus = "DELETED"
)

// Notification 通知ドメインモデル
type Notification struct {
	NotificationID string             `json:"notification_id"` // PK
	UserID         string             `json:"user_id"`         // 通知対象ユーザーID
	Title          string             `json:"title"`
	Message        string             `json:"message"`
	Status         NotificationStatus `json:"status"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}
