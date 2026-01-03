# Quick Start Guide

Get up and running in 5 minutes! ⚡

## Prerequisites Check

```bash
# Check Go version (need 1.25+)
go version

# Check PostgreSQL (need 12+)
psql --version

# Check Docker (optional)
docker --version
```

## 5-Minute Setup

### 1. Clone & Configure (1 min)

```bash
git clone <repository-url>
cd backend-svc
cp .env.template .env
# Edit .env with your credentials
```

### 2. Start Database (1 min)

```bash
# Using Docker (recommended)
docker run --name postgres-ig \
  -e POSTGRES_USER=ig_user \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=ig_test \
  -p 5432:5432 -d postgres:15

# Or use your local PostgreSQL
```

### 3. Initialize Database (1 min)

```bash
# Replace <LONG_LIVED_ACCESS_TOKEN> in database/init.sql first
psql -h localhost -U ig_user -d ig_test -f database/init.sql
```

### 4. Install & Run (2 min)

```bash
go mod download
go run cmd/main.go
```

### 5. Test It

```bash
# In another terminal
curl http://localhost:8080/ready
curl http://localhost:8080/media
```

🎉 **Done!** Service is running at `http://localhost:8080`

---

## Daily Commands

### Development
```bash
make run              # Start the service
make test             # Run all tests
make pre-commit       # Before committing
```

### Testing
```bash
make test-integration # Integration tests
make test-unit        # Unit tests
go test -v ./...      # Verbose output
```

### Docker
```bash
make docker-build     # Build image
make docker-run       # Run container
make docker-stop      # Stop container
```

---

## Common Tasks

### Add a New Test

1. Write test in `tests/integration/your_test.go`
2. Run: `go test ./tests/integration/ -run YourTest -v`
3. Update README.md test table
4. Update CHANGELOG.md
5. Run: `make pre-commit`
6. Commit: `git commit -m "test: add your test"`

### Fix a Bug

1. Write test that reproduces the bug
2. Fix the code
3. Verify: `make test`
4. Update CHANGELOG.md under "Fixed"
5. Run: `make pre-commit`
6. Commit: `git commit -m "fix: description"`

### Add a Feature

1. Write tests first (TDD)
2. Implement feature
3. Update README if needed
4. Update CHANGELOG.md under "Added"
5. Run: `make pre-commit`
6. Commit: `git commit -m "feat: description"`

---

## Troubleshooting

### Service won't start
```bash
# Check environment variables
cat .env

# Check database connection
psql -h localhost -U ig_user -d ig_test -c "SELECT 1"

# Check logs
tail -f logs.txt
```

### Tests failing
```bash
# Clean everything
make clean

# Reinstall dependencies
make deps

# Run tests with verbose output
make test -v
```

### Database issues
```bash
# Recreate database
dropdb -h localhost -U ig_user ig_test
createdb -h localhost -U ig_user ig_test
psql -h localhost -U ig_user -d ig_test -f database/init.sql
```

---

## Environment Variables Explained

| Variable | Purpose | Example |
|----------|---------|---------|
| `APP_ID` | Facebook App ID | `123456789` |
| `APP_SECRET` | Facebook App Secret | `abcdef123456` |
| `IG_USER_ID` | Instagram Business Account ID | `987654321` |
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:pass@localhost/db` |

---

## File Structure (Most Important)

```
backend-svc/
├── cmd/main.go              # Start here - entry point
├── api/handlers.go          # API endpoints
├── internal/
│   ├── instagram/service.go # Instagram API logic
│   ├── cache/memory.go      # In-memory cache
│   └── token/runtime.go     # Token management
├── tests/integration/       # Your tests go here
├── .env                     # Your configuration
└── README.md                # Full documentation
```

---

## Architecture in 30 Seconds

```
User Request
    ↓
API Handler (/media)
    ↓
In-Memory Cache (instant response)
    ↑
Scheduler (syncs every 15 min)
    ↑
Instagram API
```

**Key Points:**
- All media served from memory (fast!)
- Token stored in: PostgreSQL → Disk → Memory
- Auto-refresh token before expiration
- Tests ensure reliability

---

## Need More Info?

- **Full setup**: See [README.md](README.md)
- **Contributing**: See [CONTRIBUTING.md](CONTRIBUTING.md)
- **Changes**: See [CHANGELOG.md](CHANGELOG.md)
- **Architecture**: Check `Arch.png`

---

## Quick Links

- Test your setup: `curl http://localhost:8080/ready`
- View all media: `curl http://localhost:8080/media`
- Run tests: `make test`
- Pre-commit: `make pre-commit`

**Remember**: Every commit must pass `make pre-commit` ✅
