package handler_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/feliperosa/aws-lambda-go/internal/config"
	"github.com/feliperosa/aws-lambda-go/internal/handler"
	"github.com/feliperosa/aws-lambda-go/internal/service"
)

func TestHandleAPIGateway(t *testing.T) {
	// Setup mock config and logger
	cfg := &config.AppConfig{
		Environment: "test",
		TableName:   "mock-table",
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	
	// Inject the actual service for integration testing
	// (or you could inject a mock if the service had complex dependencies)
	svc := service.NewHelloService()
	h := handler.New(cfg, logger, svc)

	tests := []struct {
		name           string
		request        events.APIGatewayProxyRequest
		expectedStatus int
		expectedMsg    string
	}{
		{
			name: "Success with query parameter",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "GET",
				Path:       "/hello",
				QueryStringParameters: map[string]string{
					"name": "Felipe",
				},
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "Hello, Felipe!",
		},
		{
			name: "Validation error (name too short)",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "GET",
				Path:       "/hello",
				QueryStringParameters: map[string]string{
					"name": "A",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "", // Should contain error field
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := h.HandleAPIGateway(context.Background(), tt.request)

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			var body handler.ResponseBody
			_ = json.Unmarshal([]byte(resp.Body), &body)

			if tt.expectedStatus == http.StatusOK && body.Message != tt.expectedMsg {
				t.Errorf("expected message '%s', got '%s'", tt.expectedMsg, body.Message)
			}
			
			if tt.expectedStatus != http.StatusOK && body.Status != "error" {
				t.Errorf("expected status 'error', got '%s'", body.Status)
			}
		})
	}
}
