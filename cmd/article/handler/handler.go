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
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/andre-karrlein/hambach-api/model"
	"github.com/andre-karrlein/hambach-api/util"
)

// ErrNoArticle indicates that API failed in some way
var ErrNoArticle = errors.New("failed to get Article")

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
		data, err := json.MarshalIndent(getSpecificArticle(handler, id), "", "    ")
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

	data, err := json.MarshalIndent(getAllArticle(handler), "", "    ")
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

func getAllArticle(handler lambdaHandler) []model.Article {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	out, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("Articles-6tvyheax4jd3dk4j2tertab7ru-staging"),
	})

	if err != nil {
		panic(err)
	}

	content := []model.Article{}
	for _, s := range out.Items {
		item := model.Article{}

		err = dynamodbattribute.UnmarshalMap(s, &item)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
		content = append(content, item)
	}

	return content
}

func getSpecificArticle(handler lambdaHandler, id string) model.Article {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	expr, err := expression.NewBuilder().WithFilter(
		expression.Equal(expression.Name("id"), expression.Value(id)),
	).Build()
	if err != nil {
		panic(err)
	}

	out, err := svc.Scan(&dynamodb.ScanInput{
		TableName:                 aws.String("Articles-6tvyheax4jd3dk4j2tertab7ru-staging"),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	if err != nil {
		panic(err)
	}

	content := model.Article{}

	err = dynamodbattribute.UnmarshalMap(out.Items[0], &content)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return content
}
