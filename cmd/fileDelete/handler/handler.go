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
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/andre-karrlein/hambach-api/model"
	"github.com/andre-karrlein/hambach-api/util"
)

// ErrNoFiles indicates that API failed in some way
var ErrNoFiles = errors.New("failed to get files")

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

	id, ok := request.PathParameters["id"]
	if !ok {
		handler.logger.Println("No id given.")
		response.StatusCode = http.StatusBadGateway
	}

	err = deleteFile(handler, id)
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

	response.StatusCode = http.StatusAccepted
	response.Body = string(data)
	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}

	return response, nil
}

func deleteFile(handler lambdaHandler, id string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create S3 service client
	svc := s3.New(sess)
	file := getFileById(handler, svc, id)
	key := file.Key

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String("hambach"), Key: aws.String(key)})
	if err != nil {
		log.Fatalf("Unable to delete object %q from bucket hambach, %v", key, err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String("hambach"),
		Key:    aws.String(key),
	})

	return nil
}

func getFileById(handler lambdaHandler, svc *s3.S3, id string) model.File {

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String("hambach")})
	if err != nil {
		log.Fatalln(err)
	}

	var file model.File
	for _, item := range resp.Contents {
		if *item.ETag == id {
			file = model.File{
				ID:           *item.ETag,
				Key:          *item.Key,
				LastModified: item.LastModified.String(),
			}
		}
	}

	return file
}
