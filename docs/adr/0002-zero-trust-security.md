# 2. Zero-Trust Security & Secrets Management

Date: 2026-04-23

## Status

Accepted

## Context

Production Lambda functions often handle sensitive data (PII) and credentials. We need to prevent accidental leakage of these secrets in environment variables or logs, and ensure that our IAM policies are as restrictive as possible.

## Decision

1.  **Secrets:** We will use the **AWS Parameters and Secrets Lambda Extension**. This prevents storing sensitive keys (like API_KEY or DB_PASSWORD) in cleartext environment variables. The application will fetch secrets via a local HTTP call (localhost:2773), which is cached and managed by the extension.
2.  **Least Privilege:** All IAM policies in `template.yaml` must be granular. We avoid `ManagedPolicyArns` and prefer `Policies` that target specific resources (ex: `DynamoDBCrudPolicy` limited to a single table name).
3.  **Dependency Safety:** We integrate `govulncheck` in our GitHub Actions. This tool officially analyzes our dependency graph for known security vulnerabilities.
4.  **PII Masking:** We implement the `slog.LogValuer` interface for sensitive types. Any variable typed as `SensitiveString` will automatically be masked as `REDACTED` in CloudWatch logs.

## Consequences

-   **Positive:** Significantly reduced attack surface. Credential leaks via logs or environment dumps are virtually eliminated. Dependencies are continuously audited for safety.
-   **Negative:** Adds a slight complexity to local development (need to simulate the Secrets Extension or use mock values).
