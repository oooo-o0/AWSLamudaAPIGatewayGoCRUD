package middleware

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func NewErrorResponse(status int, message string) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]string{"error": message})
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
		Headers:    DefaultHeaders(),
	}
}

func DefaultHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
	}
}
