package service_test

import (
	"context"
	"testing"

	"github.com/feliperosa/aws-lambda-go/internal/service"
)

func TestSayHello(t *testing.T) {
	svc := service.NewHelloService()

	tests := []struct {
		name          string
		req           service.HelloRequest
		wantErr       bool
		expectedGreet string
	}{
		{
			name:          "Valid name",
			req:           service.HelloRequest{Name: "Felipe"},
			wantErr:       false,
			expectedGreet: "Hello, Felipe!",
		},
		{
			name:    "Empty name (validation error)",
			req:     service.HelloRequest{Name: ""},
			wantErr: true,
		},
		{
			name:    "Name too short (validation error)",
			req:     service.HelloRequest{Name: "A"},
			wantErr: true,
		},
		{
			name:    "The Unspeakable Name (domain error)",
			req:     service.HelloRequest{Name: "Voldemort"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.SayHello(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("SayHello() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && resp.Greeting != tt.expectedGreet {
				t.Errorf("expected greet %s, got %s", tt.expectedGreet, resp.Greeting)
			}
		})
	}
}
