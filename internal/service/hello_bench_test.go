package service_test

import (
	"context"
	"strings"
	"testing"

	"github.com/feliperosa/aws-lambda-go/internal/service"
)

// BenchmarkSayHello measures the performance and memory allocations of the core business logic.
// Run with: go test -bench=. -benchmem ./internal/service
func BenchmarkSayHello(b *testing.B) {
	svc := service.NewHelloService()
	req := service.HelloRequest{Name: "PerformanceTest"}
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs() // Important: Shows memory allocations per operation
	
	for i := 0; i < b.N; i++ {
		_, err := svc.SayHello(ctx, req)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// FuzzSayHello uses Go 1.18+ fuzzing to automatically generate random inputs
// to test the robustness of our validation and business logic.
// Run with: go test -fuzz=FuzzSayHello ./internal/service
func FuzzSayHello(f *testing.F) {
	svc := service.NewHelloService()
	ctx := context.Background()

	// Seed corpus: provide some examples of valid and invalid inputs
	f.Add("Felipe")
	f.Add("")
	f.Add("A")
	f.Add("Voldemort")
	f.Add(strings.Repeat("A", 100)) // Too long string

	f.Fuzz(func(t *testing.T, randomName string) {
		req := service.HelloRequest{Name: randomName}
		resp, err := svc.SayHello(ctx, req)

		// 1. Validation Logic checks
		if len(randomName) < 2 || len(randomName) > 50 {
			if err == nil {
				t.Errorf("expected validation error for string of length %d, got none", len(randomName))
			}
			return
		}

		// 2. Business Logic checks
		if randomName == "Voldemort" {
			if err == nil {
				t.Error("expected domain error for 'Voldemort', got none")
			}
			return
		}

		// 3. Success checks
		if err != nil {
			t.Errorf("unexpected error for valid input '%s': %v", randomName, err)
			return
		}
		expectedGreet := "Hello, " + randomName + "!"
		if resp.Greeting != expectedGreet {
			t.Errorf("expected '%s', got '%s'", expectedGreet, resp.Greeting)
		}
	})
}
