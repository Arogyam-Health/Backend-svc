# Contributing Guide

Thank you for contributing to the Instagram Media Backend Service! This guide will help you get started.

## 🚦 Before You Commit

**Every commit MUST pass pre-commit checks.** This is not optional.

```bash
make pre-commit
```

This will:
1. ✅ Format your code (`gofmt`)
2. ✅ Run static analysis (`go vet`)
3. ✅ Run all tests (unit + integration)
4. ✅ Verify the build succeeds

## 📝 Commit Checklist

Before committing, ensure:

- [ ] Code is formatted: `make fmt`
- [ ] All tests pass: `make test`
- [ ] Pre-commit checks pass: `make pre-commit`
- [ ] **CHANGELOG.md is updated** (see below)
- [ ] Commit message follows the format (see below)

## 📋 Updating CHANGELOG.md

**Required for EVERY commit that adds, changes, or fixes functionality.**

### How to Update

1. Open `CHANGELOG.md`
2. Find the `[Unreleased]` section
3. Add your changes under the appropriate category:
   - **Added**: New features, tests, files
   - **Changed**: Modifications to existing functionality
   - **Fixed**: Bug fixes

### Example

```markdown
## [Unreleased]

### Added
- Added pagination support for media endpoint
- Added TestMediaPagination integration test

### Changed
- Updated cache implementation to support larger datasets

### Fixed
- Fixed race condition in token refresh scheduler
```

### Bad Example (Don't do this)

```markdown
### Added
- Updated stuff
- Fixed things
- Changed code
```

Be specific and clear!

## 💬 Commit Message Format

Use this format for all commits:

```
<type>: <short description>

[optional body with more details]

[optional footer with issue references]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `test`: Adding or modifying tests
- `docs`: Documentation changes
- `refactor`: Code refactoring (no functionality change)
- `chore`: Maintenance tasks, dependencies
- `perf`: Performance improvements

### Examples

**Good commits:**
```
feat: add pagination to media endpoint

- Added page and limit query parameters
- Updated cache to support offset-based retrieval
- Added TestMediaPagination integration test
```

```
fix: resolve race condition in token refresh

The token refresh scheduler had a race condition when
multiple goroutines accessed the token simultaneously.
Added mutex locking to prevent concurrent access.

Fixes #42
```

```
test: add test for media deletion cascade

Ensures that when media is deleted from Instagram,
it's also removed from our cache during the next sync.
```

**Bad commits:**
```
Updated code
Fixed stuff
Changes
WIP
asdf
```

## 🧪 Adding New Tests

### Step 1: Write Your Test

**Integration tests** go in `tests/integration/`:
```go
package integration

import "testing"

func TestYourNewFeature(t *testing.T) {
    // Setup
    // Action
    // Assert
}
```

**Unit tests** go alongside the code in `internal/<module>/`:
```go
package mymodule

import "testing"

func TestMyFunction(t *testing.T) {
    // Test logic
}
```

### Step 2: Run the Test

```bash
# Run your specific test
go test ./tests/integration/ -run TestYourNewFeature -v

# Run all tests
make test
```

### Step 3: Update Documentation

Add your test to the README.md test table:

```markdown
| Test | File | Description |
|------|------|-------------|
| **Your Test Name** | `your_file.test.go` | What it tests |
```

### Step 4: Update CHANGELOG.md

```markdown
## [Unreleased]

### Added
- Added TestYourNewFeature to verify media deletion behavior
```

### Step 5: Pre-Commit and Commit

```bash
make pre-commit
git add .
git commit -m "test: add test for media deletion"
```

## 🔄 Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feat/your-feature-name
```

### 2. Make Your Changes

Write code, add tests, update docs.

### 3. Run Tests Frequently

```bash
# Quick check during development
make quick-check

# Full test suite
make test
```

### 4. Update CHANGELOG.md

Document your changes in the `[Unreleased]` section.

### 5. Pre-Commit Checks

```bash
make pre-commit
```

### 6. Commit

```bash
git add .
git commit -m "feat: add your feature"
```

### 7. Push and Create PR

```bash
git push origin feat/your-feature-name
```

Create a Pull Request with:
- Clear description of changes
- Link to related issues
- Screenshots if applicable

## 🏗️ Code Style Guidelines

### Go Best Practices

- **Formatting**: Use `gofmt` (automatic with `make fmt`)
- **Naming**: Use camelCase for local variables, PascalCase for exported
- **Comments**: Public functions must have doc comments
- **Error Handling**: Always check and handle errors
- **Testing**: Write table-driven tests when possible

### Example: Good Test Structure

```go
func TestMediaCache(t *testing.T) {
    tests := []struct {
        name     string
        input    []Media
        expected int
    }{
        {"empty cache", []Media{}, 0},
        {"single item", []Media{{ID: "1"}}, 1},
        {"multiple items", []Media{{ID: "1"}, {ID: "2"}}, 2},
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

## 🐛 Debugging Tests

### Run Single Test

```bash
go test ./tests/integration/ -run TestSpecificTest -v
```

### Enable Verbose Output

```bash
go test -v ./...
```

### Use Test Environment

Tests use `.env.test` for configuration. Create it if needed:

```bash
cp .env.template .env.test
```

## 📊 Coverage Reports

Generate a coverage report:

```bash
make test-coverage
```

Open `coverage.html` in your browser to see detailed coverage.

## ❓ Common Issues

### "Tests failing locally but passed before"

1. Clean and rebuild:
   ```bash
   make clean
   make deps
   make test
   ```

2. Check if `.env.test` is properly configured

3. Ensure PostgreSQL is running and accessible

### "Pre-commit checks failing on formatting"

```bash
make fmt
make pre-commit
```

### "Database connection errors in tests"

Check `DATABASE_URL` in `.env.test`:
```bash
DATABASE_URL=postgres://ig_user:password@localhost:5432/ig_test?sslmode=disable
```

Verify PostgreSQL is running:
```bash
psql -h localhost -U ig_user -d ig_test -c "SELECT 1"
```

## 🎯 Quick Reference

```bash
# Daily workflow
make fmt                 # Format code
make test               # Run all tests
make pre-commit         # Pre-commit checks

# Testing
make test-integration   # Integration tests only
make test-unit          # Unit tests only
make test-coverage      # Generate coverage report

# Building
make build              # Build binary
make run                # Run locally
make docker-build       # Build Docker image

# Utilities
make clean              # Clean artifacts
make deps               # Install dependencies
```

## 🤝 Getting Help

- **Questions?** Ask in the team chat or create an issue
- **Bug found?** Open an issue with reproduction steps
- **Feature idea?** Discuss in team meeting or create an issue

## ✨ Recognition

Quality contributions are recognized! Keep these guidelines in mind:
- Write clear, tested code
- Document your changes
- Help review others' PRs
- Share knowledge with the team

Thank you for helping make this project better! 🚀
