// Package main serves as the entry point for the AWS Lambda function.
// It handles the initialization of the execution environment (Cold Start)
// and dependency injection.
//
// Author: Felipe Rosa
// License: MIT
package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/feliperosa/aws-lambda-go/internal/config"
	"github.com/feliperosa/aws-lambda-go/internal/handler"
	"github.com/feliperosa/aws-lambda-go/internal/service"
)

var (
	// appHandler is initialized once during the Lambda Cold Start.
	// This instance is reused across multiple "warm" invocations.
	appHandler *handler.Handler
)

// init runs during the "Init Phase" of the Lambda lifecycle.
// AWS does not charge for the CPU time during this phase (up to 10s),
// making it the ideal place for heavy initialization logic.
func init() {
	// 1. Configuration Loading: Fail fast if environment variables are missing.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("CRITICAL: Failed to load configuration: %v", err)
	}

	// 2. Structured Logging: JSON format for CloudWatch Logs.
	var programLevel = new(slog.LevelVar)
	if cfg.LogLevel == "DEBUG" {
		programLevel.Set(slog.LevelDebug)
	} else {
		programLevel.Set(slog.LevelInfo)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: programLevel,
	}))
	slog.SetDefault(logger)

	// 3. Dependency Injection: Initialize services.
	// In Big Tech environments, this is where we'd inject DB repositories.
	helloService := service.NewHelloService()

	// 4. Handler Injection: Decoupling AWS Transport from Business Logic.
	appHandler = handler.New(cfg, logger, helloService)

	logger.Info("Lambda execution environment initialized successfully", 
		"env", cfg.Environment, 
		"arch", "arm64",
	)
}

func main() {
	// Start the Lambda runtime loop.
	// This call blocks until the function is terminated.
	lambda.Start(appHandler.HandleAPIGateway)
}
