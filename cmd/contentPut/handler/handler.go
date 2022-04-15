package handler

import (
	"context"
	"encoding/json"
	"errors"
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

	app_key := os.Getenv("WRITE_KEY")

	if key != app_key {
		response.StatusCode = http.StatusBadGateway
		response.Body = string("Invalid APP Key!")

		return response, nil
	}

	c := model.Content{}

	json.Unmarshal([]byte(request.Body), &c)

	err = saveContent(handler, c)
	if err != nil {
		response.StatusCode = http.StatusBadGateway
		response.Body = string("Error creating or updating post")
	}

	data, err := json.MarshalIndent("", "", "    ")
	if err != nil {
		handler.logger.Print("Failed to JSON marshal response.\nError: %w", err)
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	response.StatusCode = http.StatusCreated
	response.Body = string(data)

	return response, nil
}

func saveContent(handler lambdaHandler, content model.Content) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(content)
	if err != nil {
		log.Fatalf("Got error marshalling content item: %s", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("hambach.content"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
