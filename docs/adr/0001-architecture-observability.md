# 1. Use Graviton2, AWS SAM, and Clean Architecture

Date: 2026-04-23

## Status

Accepted

## Context

We need a boilerplate for AWS Lambda in Go that is production-ready, highly performant, and maintainable by a large engineering team (Big Tech standards). The architecture must separate AWS infrastructure concerns from the core business logic.

## Decision

1.  **Hardware:** We will use `arm64` (Graviton2) with `provided.al2023`. It provides better performance per watt and is roughly 20% cheaper than x86_64.
2.  **Architecture:** We adopt the **Service/Handler Pattern** (a simplified Hexagonal Architecture). Handlers act as inbound adapters (translating AWS events). Services hold pure domain logic. This ensures domain logic is 100% unit-testable without mocking AWS resources.
3.  **Observability:** We will use `log/slog` for structured logging with injected `AWSTraceID` for zero-search correlation in CloudWatch. We will use AWS CloudWatch Embedded Metric Format (EMF) via standard output (`fmt.Println`) to avoid heavy SDK dependencies and reduce cold start times.
4.  **Deployment Strategy:** We configure `DeploymentPreference` in SAM to use Canary releases (`Canary10Percent5Minutes`) coupled with CloudWatch Alarms to automatically rollback upon error spikes.

## Consequences

-   **Positive:** Highly decoupled code, excellent test coverage (including Fuzzing and Benchmarks), safe automated rollouts, and granular observability.
-   **Negative:** Adds slight boilerplate overhead compared to a single-file Lambda script. Developers must adhere strictly to dependency injection and interface definitions.
