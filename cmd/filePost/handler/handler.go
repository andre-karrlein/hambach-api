package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

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

	f := model.UploadedFile{}

	json.Unmarshal([]byte(request.Body), &f)

	err = saveFile(handler, f)
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
	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}

	return response, nil
}

func saveFile(handler lambdaHandler, file model.UploadedFile) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	uploader := s3manager.NewUploader(sess)

	parseData(file)
	fileData, err := os.Open("/tmp/" + file.Name)
	if err != nil {
		log.Fatalf("Unable to open file %q, %v", file.Name, err)
	}

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("hambach"),
		Key:    aws.String(file.Name),
		Body:   fileData,
	})
	if err != nil {
		// Print the error and exit.
		log.Fatalf("Unable to upload %q to hambach, %v", file.Name, err)
	}

	return nil
}

func parseData(file model.UploadedFile) {
	dec, err := base64.StdEncoding.DecodeString(file.Data)
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.Create("/tmp/" + file.Name)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		log.Fatalln(err)
	}
	if err := f.Sync(); err != nil {
		log.Fatalln(err)
	}
}
