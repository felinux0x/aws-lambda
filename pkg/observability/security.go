package observability

import (
	"log/slog"
)

// SensitiveString defines a type that automatically masks its value when logged.
// Big-Tech standard for PII (Personally Identifiable Information) protection.
type SensitiveString string

// LogValue implements the slog.LogValuer interface.
func (s SensitiveString) LogValue() slog.Value {
	return slog.StringValue("REDACTED")
}

// Handler with PII Protection example
// In a real scenario, you could implement a global middleware that checks keys like "email", "cpf", "password".
