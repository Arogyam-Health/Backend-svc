# Pre-Commit Checklist

**Print this and keep it at your desk! 📋**

---

## ✅ Before Every Commit

```bash
# 1. Format your code
make fmt

# 2. Run all tests
make test

# 3. Run pre-commit checks
make pre-commit
```

---

## ✅ Did You Update?

- [ ] **CHANGELOG.md** - Added your changes under `[Unreleased]`
- [ ] **README.md** - Updated if you added tests or features
- [ ] **Test files** - Tests written and passing

---

## ✅ Commit Message Format

```
<type>: <short description>

[optional details]
```

**Common Types:**
- `feat` - New feature
- `fix` - Bug fix
- `test` - New or updated tests
- `docs` - Documentation only
- `refactor` - Code restructuring

**Examples:**
```bash
git commit -m "feat: add pagination to media endpoint"
git commit -m "fix: resolve token refresh race condition"
git commit -m "test: add test for concurrent access"
git commit -m "docs: update setup instructions"
```

---

## ✅ If Pre-Commit Fails

### Formatting Issues
```bash
make fmt
make pre-commit
```

### Test Failures
```bash
# Run specific test
go test ./tests/integration/ -run TestName -v

# Clean and retry
make clean
make deps
make test
```

### Build Failures
```bash
# Check for syntax errors
go build ./cmd/main.go

# Check imports
go mod tidy
```

---

## ⚡ Quick Commands

| Command | Purpose |
|---------|---------|
| `make run` | Start the service |
| `make test` | Run all tests |
| `make test-integration` | Run integration tests |
| `make pre-commit` | **Required before commit** |
| `make fmt` | Format code |
| `make clean` | Clean build artifacts |
| `make help` | Show all commands |

---

## 🚫 Never Commit Without

1. ❌ Running `make pre-commit`
2. ❌ Updating `CHANGELOG.md`
3. ❌ Proper commit message format
4. ❌ Passing tests

---

## 📚 Need Help?

- **Full guide**: See [CONTRIBUTING.md](../CONTRIBUTING.md)
- **Quick start**: See [QUICKSTART.md](../QUICKSTART.md)
- **Detailed docs**: See [README.md](../README.md)

---

**Remember**: `make pre-commit` before every commit! ✅

No exceptions. It keeps our codebase clean and prevents broken builds.
