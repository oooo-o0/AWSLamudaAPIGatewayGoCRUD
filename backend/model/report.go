package model

import "time"

type ReportType string

const (
	ReportTypeDaily   ReportType = "daily"
	ReportTypeWeekly  ReportType = "weekly"
	ReportTypeMonthly ReportType = "monthly"
)

type Report struct {
	ReportID      string     `json:"reportId" dynamodbav:"ReportID"`
	UserID        string     `json:"userId" dynamodbav:"UserID"` // 作成者または対象ユーザー
	ReportType    ReportType `json:"reportType" dynamodbav:"ReportType"`
	Title         string     `json:"title" dynamodbav:"Title"`
	Content       string     `json:"content" dynamodbav:"Content"`             // JSONやテキストなど自由形式
	GeneratedDate time.Time  `json:"generatedDate" dynamodbav:"GeneratedDate"` // レポート生成日時
	CreatedAt     time.Time  `json:"createdAt" dynamodbav:"CreatedAt"`
	UpdatedAt     time.Time  `json:"updatedAt" dynamodbav:"UpdatedAt"`
}
