package model

import "time"

// AuthToken は認証トークン情報を表します。
type AuthToken struct {
	TokenID   string    `json:"tokenId" dynamodbav:"TokenID"`
	UserID    string    `json:"userId" dynamodbav:"UserID"`
	Token     string    `json:"token" dynamodbav:"Token"`   // JWTなどのトークン文字列
	Expiry    time.Time `json:"expiry" dynamodbav:"Expiry"` // 有効期限
	CreatedAt time.Time `json:"createdAt" dynamodbav:"CreatedAt"`
	Revoked   bool      `json:"revoked" dynamodbav:"Revoked"` // 失効フラグ
}
