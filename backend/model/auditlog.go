package model

import "time"

// AuditLog は操作ログを表します。
type AuditLog struct {
	LogID       string    `json:"logId" dynamodbav:"LogID"`
	UserID      string    `json:"userId" dynamodbav:"UserID"`
	EventType   string    `json:"eventType" dynamodbav:"EventType"`                   // 例: CREATE, UPDATE, DELETE, LOGIN
	Description string    `json:"description" dynamodbav:"Description"`               // 操作内容や詳細メッセージ
	Timestamp   time.Time `json:"timestamp" dynamodbav:"Timestamp"`                   // イベント発生日時
	Metadata    string    `json:"metadata,omitempty" dynamodbav:"Metadata,omitempty"` // 任意の追加情報(JSONなど
}
