package handler

import (
	"context"
	"encoding/json"
	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/backend/repository"
	"kaigo-insurance-system/backend/service"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

var userService service.UserService

func init() {
	sess := session.Must(session.NewSession())
	repo := repository.NewUserRepository(sess, os.Getenv("TABLE_NAME"))
	userService = service.NewUserService(repo)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "POST":
		return handleCreate(ctx, req)
	case "GET":
		return handleGet(ctx, req)
	case "PUT":
		return handleUpdate(ctx, req)
	case "DELETE":
		return handleDelete(ctx, req)
	default:
		return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed}, nil
	}
}

func handleCreate(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var user model.User
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	if err := userService.CreateUser(ctx, &user); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusCreated}, nil
}

func handleGet(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.QueryStringParameters["user_id"]
	user, err := userService.GetUser(ctx, id)
	if err != nil || user == nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, nil
	}

	body, _ := json.Marshal(user)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func handleUpdate(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var user model.User
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	if err := userService.UpdateUser(ctx, &user); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
}

func handleDelete(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.QueryStringParameters["user_id"]
	if err := userService.DeleteUser(ctx, id); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: http.StatusNoContent}, nil
}

func main() {
	lambda.Start(Handler)
}
