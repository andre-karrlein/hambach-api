package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/andre-karrlein/hambach-api/cmd/content/handler"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Llongfile)
	logger.Println("Lambda one has started.")

	// The main goroutine in a Lambda might never run its deferred statements.
	// This is because of how the Lambda is shutdown.
	// https://docs.aws.amazon.com/lambda/latest/dg/runtimes-context.html#runtimes-lifecycle-shutdown
	defer logger.Println("Lambda one has stopped.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize a random seed for the Pokémon API.
	rand.Seed(time.Now().UnixNano())

	h := handler.New(logger)

	lambda.StartHandlerWithContext(ctx, h)
}
