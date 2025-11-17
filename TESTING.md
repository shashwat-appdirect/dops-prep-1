# Testing Guide

## Test Structure

The codebase includes comprehensive unit and integration tests:

### Unit Tests

- **Middleware Tests** (`backend/internal/middleware/auth_test.go`): Tests authentication middleware with various token scenarios
- **Handler Tests** (`backend/internal/handlers/*_test.go`): Tests for all handler functions
  - Registration handlers
  - Admin handlers  
  - Speaker CRUD operations
  - Session CRUD operations

### Integration Tests

- **API Integration Tests** (`backend/integration/api_test.go`): Full HTTP endpoint testing

## Running Tests

### All Tests
```bash
make test
# or
cd backend && go test ./... -v
```

### Unit Tests Only
```bash
make test-unit
# or
cd backend && go test ./internal/... -v -short
```

### Integration Tests Only
```bash
make test-integration
# or
cd backend && go test ./integration/... -v
```

### Specific Package
```bash
cd backend && go test ./internal/middleware -v
```

## Test Notes

Some tests are skipped as they require full Firestore mock implementation. The test structure demonstrates the testing patterns:

- Tests use `testify` for assertions
- HTTP testing uses `httptest` package
- Mock database interface allows for testability
- Integration tests test full HTTP request/response cycle

## Mock Database

The `MockDB` in `backend/internal/database/mock.go` provides a basic in-memory mock. For full Firestore functionality testing, consider:

1. Using Firestore Emulator for integration tests
2. Creating a more comprehensive mock that implements all Firestore interfaces
3. Using table-driven tests with test fixtures

## Test Coverage

Current test coverage includes:
- ✅ Authentication middleware (all scenarios)
- ✅ Admin login endpoint
- ✅ Handler input validation
- ✅ HTTP status code verification
- ⚠️ Full Firestore operations (requires enhanced mock or emulator)

