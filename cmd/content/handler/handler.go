package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/andre-karrlein/hambach-api/model"
	"github.com/andre-karrlein/hambach-api/util"
)

// ErrNoContent indicates that API failed in some way
var ErrNoContent = errors.New("failed to get Content")

type lambdaHandler struct {
	customString string
	logger       *log.Logger
}

// New creates a new handler for Lambda one.
func New(logger *log.Logger) lambda.Handler {
	return util.NewHandlerV1(lambdaHandler{
		logger: logger,
	})
}

// Handle implements util.LambdaHTTPV1 interface. It contains the logic for the handler.
func (handler lambdaHandler) Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (response *events.APIGatewayProxyResponse, err error) {
	response = &events.APIGatewayProxyResponse{}

	key := request.QueryStringParameters["appkey"]

	app_key := os.Getenv("READ_KEY")

	if key != app_key {
		response.StatusCode = http.StatusBadGateway
		response.Body = string("Invalid APP Key!")

		return response, nil
	}
	id, ok := request.PathParameters["id"]
	if ok && id != "" {
		data, err := json.MarshalIndent(getSpecificContent(handler, id), "", "    ")
		if err != nil {
			handler.logger.Print("Failed to JSON marshal response.\nError: %w", err)
			response.StatusCode = http.StatusInternalServerError
			return response, nil
		}

		response.StatusCode = http.StatusOK
		response.Body = string(data)
		response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}

		return response, nil
	}

	data, err := json.MarshalIndent(getAllContent(handler), "", "    ")
	if err != nil {
		handler.logger.Print("Failed to JSON marshal response.\nError: %w", err)
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	response.StatusCode = http.StatusOK
	response.Body = string(data)
	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}

	return response, nil
}

func getAllContent(handler lambdaHandler) []model.Content {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	out, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("hambach.content"),
	})

	if err != nil {
		panic(err)
	}

	content := []model.Content{}
	for _, s := range out.Items {
		item := model.Content{}

		err = dynamodbattribute.UnmarshalMap(s, &item)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
		content = append(content, item)
	}

	return content
}

func getSpecificContent(handler lambdaHandler, id string) model.Content {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	out, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("hambach.content"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		panic(err)
	}

	content := model.Content{}

	err = dynamodbattribute.UnmarshalMap(out.Item, &content)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return content
}
