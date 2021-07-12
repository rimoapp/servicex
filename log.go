package servicex

import (
	"context"
	"log"

	"cloud.google.com/go/logging"
)

var logger *logging.Logger

func newLogger(ctx context.Context) *logging.Logger {
	if logger != nil {
		return logger
	}
	cloudProject, err := GetEnv("GOOGLE_CLOUD_PROJECT")
	if err != nil {
		log.Fatalf("failed to get GOOGLE_CLOUD_PROJECT from env: %v", err)
	}
	client, err := logging.NewClient(ctx, cloudProject)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	logger = client.Logger("ServeHTTP")
	return logger
}

// LogError formats its arguments according to the format, analogous to fmt.Printf,
// and records the text as logx message at debug level.
// The logx message will be associated with the platform request linked with the context.
func LogError(ctx context.Context, text string) {
	newLogger(ctx).StandardLogger(logging.Error).Println(text)
}

// LogWarning is like LogError, but the severity is warning level.
func LogWarning(ctx context.Context, text string) {
	newLogger(ctx).StandardLogger(logging.Warning).Println(text)
}

// LogInfo is like LogError, but the severity is info level.
func LogInfo(ctx context.Context, text string) {
	newLogger(ctx).StandardLogger(logging.Info).Println(text)
}
