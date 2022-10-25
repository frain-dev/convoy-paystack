package main

import (
	"context"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	webhooks "github.com/frain-dev/convoy-paystack"
)

func main() {
	ctx := context.Background()
	if err := funcframework.RegisterHTTPFunctionContext(ctx, "/", webhooks.WebhookEndpoint); err != nil {
		log.Printf("HTTPFunction: %v\n", err)
	}

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
