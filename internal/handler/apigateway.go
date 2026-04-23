// Package handler implements the Inbound Adapters for AWS Lambda.
// It is responsible for parsing events, invoking services, and formatting responses.
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/feliperosa/aws-lambda-go/internal/config"
	"github.com/feliperosa/aws-lambda-go/internal/service"
	"github.com/feliperosa/aws-lambda-go/pkg/observability"
)

// ResponseBody defines the standard JSON structure for all API responses.
type ResponseBody struct {
	Message string `json:"message,omitempty"` // Message on success
	Status  string `json:"status"`            // Status: success or error
	Error   string `json:"error,omitempty"`   // Error details
}

// Handler contains all dependencies required for processing events.
type Handler struct {
	Config  *config.AppConfig
	Logger  *slog.Logger
	Service service.HelloService
}

// New initializes the Handler with its required dependencies.
func New(cfg *config.AppConfig, logger *slog.Logger, svc service.HelloService) *Handler {
	return &Handler{
		Config:  cfg,
		Logger:  logger,
		Service: svc,
	}
}

// HandleAPIGateway processes incoming REST requests from Amazon API Gateway.
// It demonstrates the use of structured logging with Trace correlation.
func (h *Handler) HandleAPIGateway(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Logger enrichment: ensures every log line contains the X-Ray Trace ID.
	log := observability.LoggerWithTrace(ctx, h.Logger)

	log.InfoContext(ctx, "HTTP request received",
		"path", request.Path,
		"method", request.HTTPMethod,
		"requestId", request.RequestContext.RequestID,
	)

	// 1. Map Input: API Gateway events to Domain Service request.
	name := request.QueryStringParameters["name"]

	svcReq := service.HelloRequest{
		Name: name,
	}

	// 2. Business Logic Execution: Decoupled from AWS events.
	// Context is passed down to support timeouts and cancellations.
	resp, err := h.Service.SayHello(ctx, svcReq)
	if err != nil {
		log.ErrorContext(ctx, "Service execution failed", "error", err)
		return h.buildErrorResponse(err)
	}

	// 3. Response Formatting: Standardized JSON output.
	respBody := ResponseBody{
		Message: resp.Greeting,
		Status:  "success",
	}

	bodyBytes, _ := json.Marshal(respBody)
	return h.buildResponse(http.StatusOK, string(bodyBytes)), nil
}

// buildErrorResponse maps application errors to the correct HTTP status codes.
func (h *Handler) buildErrorResponse(err error) (events.APIGatewayProxyResponse, error) {
	// TODO: Use errors.As() to differentiate between Validation and Internal errors.
	return h.buildResponse(http.StatusBadRequest, fmt.Sprintf(`{"status":"error", "error":"%s"}`, err.Error())), nil
}

// buildResponse is a helper to construct a valid APIGatewayProxyResponse.
func (h *Handler) buildResponse(statusCode int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*", // CORS Support
		},
		Body: body,
	}
}
