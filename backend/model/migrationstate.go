package model

import "time"

// MigrationState はマイグレーション実行履歴を管理するモデルです。
type MigrationState struct {
	MigrationID  string    `json:"migrationId" dynamodbav:"MigrationID"`                       // 主キー
	AppliedAt    time.Time `json:"appliedAt" dynamodbav:"AppliedAt"`                           // 適用日時
	Description  string    `json:"description,omitempty" dynamodbav:"Description,omitempty"`   // マイグレーション内容説明
	Success      bool      `json:"success" dynamodbav:"Success"`                               // 成功/失敗フラグ
	ErrorMessage string    `json:"errorMessage,omitempty" dynamodbav:"ErrorMessage,omitempty"` // エラー詳細（失敗時）
}
