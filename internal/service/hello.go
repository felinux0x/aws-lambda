// Package service contains the core business logic of the application.
// This layer is pure Go and should remain agnostic of the transport layer (AWS Lambda).
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/feliperosa/aws-lambda-go/pkg/observability"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// HelloRequest defines the required input for the greeting service.
type HelloRequest struct {
	// Name must be between 2 and 50 characters.
	Name string `json:"name" validate:"required,min=2,max=50"`
}

// HelloResponse defines the result of the greeting service.
type HelloResponse struct {
	Greeting string `json:"greeting"`
}

// HelloService defines the contract for our business logic.
// Using an interface facilitates mocking during unit tests.
type HelloService interface {
	SayHello(ctx context.Context, req HelloRequest) (HelloResponse, error)
}

type helloService struct {
	// Repository dependencies would be added here.
}

// NewHelloService creates a concrete implementation of HelloService.
func NewHelloService() HelloService {
	return &helloService{}
}

// SayHello performs validation and generates a greeting.
// It also demonstrates the emission of business metrics via EMF.
func (s *helloService) SayHello(ctx context.Context, req HelloRequest) (HelloResponse, error) {
	// 1. Validation: Fail early if data is malformed.
	if err := validate.Struct(req); err != nil {
		return HelloResponse{}, fmt.Errorf("validation error: %w", err)
	}

	// 2. Domain Rules: Enforce business constraints.
	if req.Name == "Voldemort" {
		return HelloResponse{}, errors.New("cannot name the unspeakable")
	}

	// 3. Observability: Emit a custom business metric.
	// This happens synchronously but is ingested asynchronously by CloudWatch.
	observability.LogEMF("MyApplication/BusinessMetrics",
		map[string]string{"Environment": "production"},
		map[string]float64{"SuccessfulGreetings": 1},
	)

	return HelloResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}
