# Test Documentation Template

Use this template when adding new tests to the project.

## Test Information

**Test Name**: `TestYourFeatureName`

**File Location**: `tests/integration/your_feature.test.go` (or `internal/module/file.test.go` for unit tests)

**Category**: 
- [ ] Token Management
- [ ] Media Management
- [ ] Error Handling
- [ ] Performance
- [ ] Other: _____________

**Type**:
- [ ] Integration Test
- [ ] Unit Test
- [ ] Manual Verification Test

---

## Test Description

### What It Tests

A brief description of what this test verifies. Be specific about the behavior being tested.

**Example**: 
> Verifies that when a new media item is added to Instagram, the scheduled sync job fetches it and updates the in-memory cache without disrupting active requests.

### Why It's Important

Explain why this test is necessary and what problems it prevents.

**Example**:
> Ensures users always see the latest content without manual restarts. Prevents stale data from being served after new posts are published.

### Test Scenario

Describe the step-by-step flow of what happens in the test.

**Example**:
1. Start service with 5 media items cached
2. Mock Instagram API adds new item (ID: "new_media_123")
3. Trigger sync scheduler
4. Verify cache contains 6 items
5. Verify new item is accessible via `/media` endpoint

---

## Test Code

```go
package integration

import (
    "testing"
    // your imports
)

func TestYourFeatureName(t *testing.T) {
    // Setup
    // ...
    
    // Action
    // ...
    
    // Assert
    if got != expected {
        t.Fatalf("expected %v, got %v", expected, got)
    }
}
```

---

## Documentation Updates

### 1. README.md Update

Add to the appropriate test table in README.md:

```markdown
| **Your Test Name** | `your_feature.test.go` | Brief description of what it verifies |
```

### 2. CHANGELOG.md Update

Add to the `[Unreleased]` section:

```markdown
### Added
- Added TestYourFeatureName to verify [behavior description]
```

### 3. Test Execution

Commands to run your test:

```bash
# Run specific test
go test ./tests/integration/ -run TestYourFeatureName -v

# Run all tests in category
go test ./tests/integration/... -v

# Run with coverage
go test -cover ./tests/integration/ -run TestYourFeatureName
```

---

## Dependencies

List any setup required for the test to run:

- [ ] PostgreSQL database
- [ ] Mock Instagram API server
- [ ] Specific environment variables
- [ ] Test fixtures or seed data
- [ ] Other: _____________

**Environment Variables Needed**:
```bash
DATABASE_URL=postgres://...
IG_USER_ID=test_user
# Add others
```

---

## Edge Cases Covered

List the edge cases this test handles:

- [ ] Empty data set
- [ ] Null/nil values
- [ ] Concurrent access
- [ ] API failure scenarios
- [ ] Timeout scenarios
- [ ] Invalid input
- [ ] Large data sets
- [ ] Other: _____________

---

## Expected Outcomes

### Success Criteria

What must be true for the test to pass?

**Example**:
- ✅ Cache updated with new media
- ✅ No duplicate entries created
- ✅ All media IDs are unique
- ✅ Response time < 100ms

### Failure Indicators

What causes this test to fail?

**Example**:
- ❌ New media not added to cache
- ❌ Existing media lost during sync
- ❌ API timeout errors
- ❌ Race conditions detected

---

## Related Tests

List any related tests that should be considered:

- `TestMediaBootstrap` - Initial media loading
- `TestConcurrentMediaAccess` - Thread safety
- Related manual verification steps

---

## Maintenance Notes

### Known Limitations

Document any limitations or assumptions:

**Example**:
> This test uses a mock Instagram API server. It does not test actual Instagram API rate limiting or network issues.

### Future Improvements

Ideas for expanding or improving this test:

**Example**:
- Add test for pagination when sync returns >50 items
- Test behavior when Instagram API returns HTTP 429 (rate limit)
- Add performance benchmarking

---

## Checklist Before Committing

- [ ] Test written and passing locally
- [ ] Test added to README.md test table
- [ ] CHANGELOG.md updated
- [ ] Test includes clear assertions
- [ ] Test is idempotent (can run multiple times)
- [ ] Test cleans up resources (files, connections)
- [ ] Test has descriptive name
- [ ] Test failure messages are clear
- [ ] `make pre-commit` passes
- [ ] Code reviewed by at least one team member

---

## Example Commit Message

```
test: add test for new media sync

- Added TestMediaSyncAddsNewItems
- Verifies scheduler fetches new Instagram media
- Ensures cache updates without data loss
- Updated README.md and CHANGELOG.md

Closes #123
```

---

**Note**: Fill out all relevant sections before submitting your test for review. Remove sections that don't apply with a note explaining why.
