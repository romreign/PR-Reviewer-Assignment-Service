# Linter Configuration Guide

This document describes the linter configuration for the PR-Reviewer-Assignment-Service project.

## Overview

The project uses **golangci-lint** for comprehensive code quality analysis. 
The configuration is defined in `.golangci.yml`.
## Installation

To install golangci-lint:

```bash
# Using Go
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Or download the binary directly
# https://github.com/golangci/golangci-lint/releases
```

## Running the Linter

### Run all linters
```bash
golangci-lint run
```

### Run specific linter
```bash
golangci-lint run --enable errcheck
```

### Fix issues automatically (where possible)
```bash
golangci-lint run --fix
```

### Run with verbose output
```bash
golangci-lint run -v
```

## Enabled Linters

### Error Handling and Validation
- **errcheck**: Checks for unchecked errors in function calls. Critical for identifying missing error handling.
- **errname**: Ensures sentinel errors are prefixed with `Err` and error types are suffixed with `Error`.
- **errorlint**: Finds code that will cause problems with the error wrapping scheme (Go 1.13+).

### Code Quality 
- **godot**: Ensures comments end with proper punctuation.
- **gofmt**: Checks code formatting (automated formatting).
- **goimports**: Auto-fixes imports and applies `gofmt` rules.
- **govet**: Reports suspicious constructs and potential bugs.
- **gosimple**: Suggests simplifications in code logic.
- **ineffassign**: Detects unused variable assignments.
- **staticcheck**: Advanced static analysis for bugs and performance issues.
- **typecheck**: Type-checking similar to Go compiler front-end.
- **unconvert**: Removes unnecessary type conversions.
- **unused**: Finds unused constants, variables, functions, and types.

### Security
- **gosec**: Security-focused linter for potential vulnerabilities.
- **securego**: Examines code for security concerns.

### Performance
- **prealloc**: Recommends pre-allocating slices when size is known.

### Style and Consistency
- **stylecheck**: Go style guide compliance checker.
- **revive**: Highly configurable linter with modern replacements for `golint`.
- **nolintlint**: Ensures `nolint` directives are properly formatted.

### Maintainability
- **maintidx**: Measures maintainability index of functions (target: > 10).
- **misspell**: Finds common English spelling mistakes in comments.

### Network Safety
- **noctx**: Warns about HTTP requests without context.Context.

## Configuration Details

### Timeout
The linter has a **5-minute timeout** to prevent hanging on large codebases.

### Excluded Patterns

The following are excluded from linting for practical reasons:

1. **Test Files (`_test.go`)**: Excluded from strict error checking and security checks (more lenient for tests).
2. **Generated Code (`internal/api/`)**: Excluded from style checks (code is auto-generated).
3. **Vendor Directory**: Completely skipped.
4. **Hidden Directories**: Skipped by default (`.git`, `.github`, etc.).

## Project-Specific Rules

### Repository Layer
- Must handle database errors explicitly
- Connection pool management must be properly closed
- No raw SQL queries without parameterization (security)

### Service Layer
- All business logic errors must be wrapped with context using `fmt.Errorf`
- Public methods must have documentation comments
- Complex logic should be in functions under 100 lines (maintidx check)

### Handler Layer
- All HTTP errors must be properly logged
- Response validation must precede business logic
- No unhandled errors from JSON marshal/unmarshal

### Testing
- Test names must follow `Test*` convention
- Benchmarks must use `Benchmark*` convention
- Use table-driven tests for multiple scenarios

## Running in CI/CD

The linter is integrated into GitHub Actions CI/CD pipeline:

```bash
# CI command
golangci-lint run --timeout 5m --out-format github-actions
```

## Example Linting Session

```bash
# Run linter with all checks
$ golangci-lint run
./internal/service/pull_request_service.go:42:2: unhandled error (errcheck)

# Fix the issue by adding error handling
$ vim ./internal/service/pull_request_service.go

# Re-run to confirm fix
$ golangci-lint run
✓ No issues found
```
