# Testing Strategy & Architecture

This document explains our testing strategy and how it relates to the system architecture.

## System Architecture Overview

Refer to [Arch.png](../Arch.png) for the visual diagram.

### Component Layers

```
┌─────────────────────────────────────────┐
│         HTTP API Layer                  │
│  (handlers.go - /media, /ready)         │
└───────────────┬─────────────────────────┘
                │
┌───────────────▼─────────────────────────┐
│      Service Layer                      │
│  (instagram/service.go)                 │
│  Business logic & orchestration         │
└───────┬────────────────────┬────────────┘
        │                    │
┌───────▼────────┐   ┌──────▼──────────┐
│  Cache Layer   │   │  Token Manager  │
│  (memory.go)   │   │  (runtime.go)   │
│  In-Memory     │   │  Multi-layer    │
└───────┬────────┘   └──────┬──────────┘
        │                    │
┌───────▼────────────────────▼──────────┐
│     External Dependencies             │
│  - Instagram API                      │
│  - PostgreSQL Database                │
│  - Disk Storage (token.json)          │
└───────────────────────────────────────┘
```

---

## Testing Layers

### 1. Unit Tests
**Location**: `internal/<module>/*.test.go`

**What We Test**:
- Individual functions in isolation
- Cache operations (add, get, update)
- Token validation logic
- Scheduler timing logic

**Example**: `internal/cache/memory.test.go`
```go
func TestCacheSetAndGet(t *testing.T) {
    cache := NewStore()
    media := []Media{{ID: "1", Caption: "Test"}}
    
    cache.SetMedia(media)
    result := cache.GetAllMedia()
    
    if len(result) != 1 {
        t.Errorf("expected 1 item, got %d", len(result))
    }
}
```

**Purpose**: Verify individual components work correctly in isolation.

---

### 2. Integration Tests
**Location**: `tests/integration/*.test.go`

**What We Test**:
- Multiple components working together
- End-to-end workflows
- Database interactions
- External API mocking

**Categories**:

#### A. Token Management Integration
Tests the token lifecycle across all storage layers:

```
PostgreSQL → Disk → Memory → API Request
    ↓
Refresh → PostgreSQL → Disk → Memory
```

**Tests**:
1. `TestBootstrapLoadsFromPostgres` - Startup flow
2. `TestTokenRefreshBeforeExpiry` - Refresh flow

#### B. Media Management Integration
Tests the media sync and caching flow:

```
Instagram API → Service → Cache → HTTP Response
```

**Tests**:
1. `TestMediaBootstrap` - Initial sync
2. `TestMediaSyncAddsNewItems` - Incremental sync
3. `TestConcurrentMediaAccess` - Thread safety

#### C. Error Handling Integration
Tests system resilience:

```
Instagram API (Down) → Service → Error Handling → Graceful Degradation
```

**Tests**:
1. `TestInstagramDownDoesNotCrashService` - API failure recovery

---

### 3. Manual Verification Tests

Some scenarios require manual verification:

**Container Restart Test**:
```bash
# 1. Start service
docker run --name test-svc instagram-backend

# 2. Verify token loaded
curl http://localhost:8080/ready

# 3. Simulate disk loss
docker exec test-svc rm /app/token.json

# 4. Restart container
docker restart test-svc

# 5. Verify token reloaded from PostgreSQL
docker logs test-svc | grep "BOOTSTRAP"
curl http://localhost:8080/media
```

**Expected**: Service recovers and continues serving requests.

---

## Test Design Patterns

### Pattern 1: Arrange-Act-Assert (AAA)

```go
func TestExample(t *testing.T) {
    // Arrange: Set up test data
    db := setupTestDB(t)
    service := NewService(db)
    
    // Act: Execute the operation
    result, err := service.FetchMedia()
    
    // Assert: Verify the outcome
    if err != nil {
        t.Fatal(err)
    }
    if len(result) == 0 {
        t.Error("expected media, got none")
    }
}
```

### Pattern 2: Table-Driven Tests

```go
func TestCacheOperations(t *testing.T) {
    tests := []struct {
        name     string
        input    []Media
        expected int
    }{
        {"empty", []Media{}, 0},
        {"single", []Media{{ID: "1"}}, 1},
        {"multiple", []Media{{ID: "1"}, {ID: "2"}}, 2},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cache := NewStore()
            cache.SetMedia(tt.input)
            
            if got := len(cache.GetAllMedia()); got != tt.expected {
                t.Errorf("expected %d, got %d", tt.expected, got)
            }
        })
    }
}
```

### Pattern 3: Mock External Dependencies

```go
// tests/dummy/instagram.go
func StartDummyServer() *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Mock Instagram API responses
        json.NewEncoder(w).Encode(mockData)
    }))
}

// In test:
func TestWithMockAPI(t *testing.T) {
    srv := dummy.StartDummyServer()
    defer srv.Close()
    
    // Test against mock server
    client := instagram.NewClient(srv.URL)
    // ... test logic
}
```

---

## Test Data Management

### Test Database

**Setup**: Each test gets a fresh database:

```go
func setupTestDB(t *testing.T) *sql.DB {
    db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    
    // Clean slate
    db.Exec("DELETE FROM instagram_tokens")
    
    t.Cleanup(func() {
        db.Close()
    })
    
    return db
}
```

### Test Fixtures

**Location**: `tests/helpers/`

```go
// Insert test token
func InsertToken(db *sql.DB, token string, expiresAt time.Time) {
    db.Exec(`
        INSERT INTO instagram_tokens (id, access_token, expires_at)
        VALUES (TRUE, $1, $2)
        ON CONFLICT (id) DO UPDATE SET access_token = $1, expires_at = $2
    `, token, expiresAt)
}
```

---

## Coverage Goals

| Component | Target Coverage | Current |
|-----------|----------------|---------|
| Cache | 90% | ✅ 95% |
| Token Management | 85% | ✅ 88% |
| API Handlers | 80% | ✅ 82% |
| Scheduler | 75% | ✅ 78% |
| Instagram Service | 80% | ✅ 85% |

**Generate coverage report**:
```bash
make test-coverage
open coverage.html
```

---

## Test Execution Flow

### Development Workflow

```bash
# 1. Write test first (TDD)
vim tests/integration/my_feature.test.go

# 2. Run test (should fail)
go test ./tests/integration/ -run MyFeature -v

# 3. Implement feature
vim internal/service/my_feature.go

# 4. Run test again (should pass)
go test ./tests/integration/ -run MyFeature -v

# 5. Run all tests
make test

# 6. Pre-commit
make pre-commit
```

### CI/CD Pipeline

```
Push/PR → GitHub Actions
    ↓
1. Checkout code
    ↓
2. Setup Go + PostgreSQL
    ↓
3. Install dependencies
    ↓
4. Check formatting (make fmt)
    ↓
5. Run linter (make vet)
    ↓
6. Run all tests (make test)
    ↓
7. Build application (make build)
    ↓
8. Check CHANGELOG updated
    ↓
Pass ✅ → Merge allowed
Fail ❌ → Fix required
```

---

## Common Testing Scenarios

### Scenario 1: Adding a New Endpoint

1. **Write handler test**:
   ```go
   func TestNewEndpoint(t *testing.T) {
       req := httptest.NewRequest("GET", "/new", nil)
       w := httptest.NewRecorder()
       
       NewHandler(store).ServeHTTP(w, req)
       
       if w.Code != http.StatusOK {
           t.Errorf("expected 200, got %d", w.Code)
       }
   }
   ```

2. **Implement handler**
3. **Add integration test**
4. **Update README & CHANGELOG**
5. **Run `make pre-commit`**

### Scenario 2: Modifying Cache Logic

1. **Update unit tests** in `internal/cache/memory.test.go`
2. **Update integration tests** that use cache
3. **Run tests**: `go test ./internal/cache/... -v`
4. **Verify no regressions**: `make test`
5. **Update CHANGELOG**

### Scenario 3: Changing Database Schema

1. **Update `database/init.sql`**
2. **Update test helpers** in `tests/helpers/`
3. **Run integration tests**: `make test-integration`
4. **Document in README**
5. **Update CHANGELOG** under "Changed"

---

## Best Practices

### ✅ DO

- Write tests before code (TDD)
- Use descriptive test names: `TestFeature_Scenario_ExpectedResult`
- Clean up resources (files, connections) in tests
- Use `t.Helper()` for test helper functions
- Mock external dependencies (Instagram API, time)
- Test edge cases (empty, nil, max values)
- Run tests in isolation (no dependencies between tests)

### ❌ DON'T

- Commit without running `make pre-commit`
- Skip updating CHANGELOG for test additions
- Write tests that depend on execution order
- Use hardcoded timeouts (use mock time)
- Leave debug logs or `fmt.Println` in tests
- Test implementation details (test behavior)
- Ignore test failures ("it works on my machine")

---

## Debugging Failed Tests

### Step 1: Run Test in Verbose Mode

```bash
go test ./tests/integration/ -run TestName -v
```

### Step 2: Check Test Logs

Look for:
- SQL errors → Database connection issue
- HTTP errors → Mock server not started
- Nil pointer → Missing initialization
- Timeout → Infinite loop or deadlock

### Step 3: Add Debug Output

```go
func TestDebug(t *testing.T) {
    result := cache.GetAllMedia()
    t.Logf("Cache contents: %+v", result)
    
    // ... rest of test
}
```

### Step 4: Run in Isolation

```bash
# Clean everything
make clean
rm token.json test_token.json

# Fresh start
make deps
go test ./tests/integration/ -run TestName -v
```

### Step 5: Check Environment

```bash
# Verify .env.test exists
cat .env.test

# Check database connection
psql $DATABASE_URL -c "SELECT 1"

# Check Go version
go version
```

---

## Adding New Test Categories

When adding a completely new category of tests (e.g., "Authentication"):

1. **Create directory** (if needed): `tests/integration/auth/`
2. **Write tests**: `auth_test.go`
3. **Add section to README**:
   ```markdown
   #### Authentication Tests
   | Test | File | Description |
   |------|------|-------------|
   | **Test Name** | `auth_test.go` | What it tests |
   ```
4. **Update CHANGELOG**
5. **Document in this file**

---

## Resources

- **Go Testing**: https://golang.org/pkg/testing/
- **Table-Driven Tests**: https://github.com/golang/go/wiki/TableDrivenTests
- **Test Fixtures**: https://github.com/go-testfixtures/testfixtures
- **Mocking**: https://github.com/golang/mock

---

## Questions?

- See [CONTRIBUTING.md](../CONTRIBUTING.md) for general guidelines
- See [README.md](../README.md) for project overview
- Ask in team chat or create an issue

---

**Remember**: Good tests are documentation. Write tests that explain how the system should behave. 📝
