package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"kaigo-insurance-system/backend/common"
	"kaigo-insurance-system/backend/middleware"
	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/backend/repository"
	"kaigo-insurance-system/backend/service"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type NotificationHandler struct {
	service service.NotificationService
}

func NewNotificationHandler() (*NotificationHandler, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	dbClient := dynamodb.NewFromConfig(cfg)

	tableName := common.GetEnv("TABLE_NAME", "Notification")
	repo := repository.NewDynamoNotificationRepository(dbClient, tableName)
	svc := service.NewNotificationService(repo)

	return &NotificationHandler{
		service: svc,
	}, nil
}

func (h *NotificationHandler) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 認証ミドルウェア
	userID, err := middleware.ExtractUserIDFromRequest(req)
	if err != nil {
		return common.RespondError(http.StatusUnauthorized, "Unauthorized"), nil
	}

	path := req.Path
	method := req.HTTPMethod

	// /notifications または /notifications/{notification_id}
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	switch method {
	case http.MethodPost:
		if len(pathParts) == 1 {
			return h.createNotification(ctx, req, userID)
		}
	case http.MethodGet:
		if len(pathParts) == 1 {
			// 一覧取得（ユーザー別）
			return h.listNotifications(ctx, userID)
		}
		if len(pathParts) == 2 {
			return h.getNotification(ctx, pathParts[1], userID)
		}
	case http.MethodPut:
		if len(pathParts) == 2 {
			// 例: ステータス更新（既読化）
			return h.updateNotificationStatus(ctx, pathParts[1], req, userID)
		}
	case http.MethodDelete:
		if len(pathParts) == 2 {
			return h.deleteNotification(ctx, pathParts[1], userID)
		}
	}

	return common.RespondError(http.StatusBadRequest, "Invalid request"), nil
}

func (h *NotificationHandler) createNotification(ctx context.Context, req events.APIGatewayProxyRequest, userID string) (events.APIGatewayProxyResponse, error) {
	var input model.Notification
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return common.RespondError(http.StatusBadRequest, "Invalid request body"), nil
	}

	// UserIDは認証情報と一致させるか、管理者の場合は別途対応（ここは簡易的にuserID固定）
	input.UserID = userID

	if input.NotificationID == "" {
		return common.RespondError(http.StatusBadRequest, "notification_id is required"), nil
	}
	if input.Title == "" || input.Message == "" {
		return common.RespondError(http.StatusBadRequest, "title and message are required"), nil
	}

	err := h.service.CreateNotification(ctx, &input)
	if err != nil {
		return common.RespondError(http.StatusInternalServerError, err.Error()), nil
	}
	return common.RespondJSON(http.StatusCreated, input), nil
}

func (h *NotificationHandler) getNotification(ctx context.Context, notificationID, userID string) (events.APIGatewayProxyResponse, error) {
	notification, err := h.service.GetNotification(ctx, notificationID)
	if err != nil {
		return common.RespondError(http.StatusNotFound, "Notification not found"), nil
	}

	if notification.UserID != userID {
		return common.RespondError(http.StatusForbidden, "Forbidden"), nil
	}

	return common.RespondJSON(http.StatusOK, notification), nil
}

func (h *NotificationHandler) listNotifications(ctx context.Context, userID string) (events.APIGatewayProxyResponse, error) {
	notifications, err := h.service.ListNotificationsByUser(ctx, userID)
	if err != nil {
		return common.RespondError(http.StatusInternalServerError, err.Error()), nil
	}
	return common.RespondJSON(http.StatusOK, notifications), nil
}

func (h *NotificationHandler) updateNotificationStatus(ctx context.Context, notificationID string, req events.APIGatewayProxyRequest, userID string) (events.APIGatewayProxyResponse, error) {
	// リクエストボディにstatusがあることを想定
	var body struct {
		Status model.NotificationStatus `json:"status"`
	}
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return common.RespondError(http.StatusBadRequest, "Invalid request body"), nil
	}

	// 既存通知取得して所有者チェック
	notification, err := h.service.GetNotification(ctx, notificationID)
	if err != nil {
		return common.RespondError(http.StatusNotFound, "Notification not found"), nil
	}
	if notification.UserID != userID {
		return common.RespondError(http.StatusForbidden, "Forbidden"), nil
	}

	// ステータス更新
	err = h.service.MarkAsRead(ctx, notificationID)
	if err != nil {
		return common.RespondError(http.StatusInternalServerError, err.Error()), nil
	}
	notification.Status = body.Status
	return common.RespondJSON(http.StatusOK, notification), nil
}

func (h *NotificationHandler) deleteNotification(ctx context.Context, notificationID string, userID string) (events.APIGatewayProxyResponse, error) {
	notification, err := h.service.GetNotification(ctx, notificationID)
	if err != nil {
		return common.RespondError(http.StatusNotFound, "Notification not found"), nil
	}
	if notification.UserID != userID {
		return common.RespondError(http.StatusForbidden, "Forbidden"), nil
	}

	err = h.service.DeleteNotification(ctx, notificationID)
	if err != nil {
		return common.RespondError(http.StatusInternalServerError, err.Error()), nil
	}

	return common.RespondJSON(http.StatusNoContent, nil), nil
}
