# Backend Service for Obesity Killer 

A high-performance backend service that synchronizes and serves Instagram media content with intelligent token management and caching strategies.

[![Pre-Commit Tests](https://github.com/Arogyam-Health/Backend-svc/workflows/Pre-Commit%20Tests/badge.svg)](https://github.com/Arogyam-Health/Backend-svc/actions)

## Documentation

- **[Documentation Index](docs/DOCUMENTATION_INDEX.md)** - Overview of all documentation
- **[Quick Start Guide](QUICKSTART.md)** - Get running in 5 minutes
- **[Contributing Guide](CONTRIBUTING.md)** - How to contribute and commit guidelines
- **[Testing Guide](docs/TESTING_GUIDE.md)** - Testing strategy and patterns
- **[Changelog](CHANGELOG.md)** - Version history and changes
- **[Architecture Diagram](Arch.png)** - Visual system overview
- **[Pre-Commit Checklist](docs/PRE_COMMIT_CHECKLIST.md)** - Quick reference (printable)
---

## 🏗️ Architecture

![Architecture Diagram](Arch.png)

### System Overview

The service follows a layered architecture with the following key components:

1. **API Layer** (`api/`): HTTP handlers for media endpoints
2. **Service Layer** (`internal/instagram/`): Business logic for Instagram API integration
3. **Cache Layer** (`internal/cache/`): In-memory storage for fast media access
4. **Token Management** (`internal/token/`): Multi-layer token storage (PostgreSQL → Disk → Memory)
5. **Scheduler** (`internal/scheduler/`): Cron jobs for token refresh and media sync
6. **Bootstrap** (`internal/bootstrap/`): Service initialization and token loading

### Key Features

- **Multi-layer Token Storage**: Tokens are stored in PostgreSQL, disk (token.json), and in-memory for redundancy
- **Automatic Token Refresh**: Tokens are refreshed 7 days before expiration
- **In-Memory Media Cache**: All media served from memory for low latency
- **Scheduled Media Sync**: Periodic synchronization with Instagram API
- **Concurrent Access Support**: Thread-safe cache for multiple simultaneous users

---

## 📋 Prerequisites

- **Go**: 1.25.0 ( currently required for render deployment )
- **PostgreSQL**: 12 or higher (needed for access token management)

---

## Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd backend-svc
```

### 2. Environment Configuration

Create a `.env` file in the project root by copying the template:

```bash
cp .env.template .env
```

Edit `.env` with your actual credentials:

```env
# Instagram/Facebook API Credentials
APP_ID=<YOUR_APP_ID>
APP_SECRET=<YOUR_APP_SECRET>
IG_USER_ID=<YOUR_INSTAGRAM_USER_ID>

# Server Configuration
PORT=8080

# API Base URLs
FB_API_BASE_URL=https://graph.facebook.com/v24.0/
IG_GRAPH_API_BASE_URL=https://graph.instagram.com/

# PostgreSQL Connection
DATABASE_URL=postgres://ig_user:password@localhost:5432/ig_test?sslmode=disable
```

**How to obtain credentials:**
- `APP_ID` and `APP_SECRET`: Create an app at [Facebook Developers](https://developers.facebook.com/)
- `IG_USER_ID`: Your Instagram Business Account ID
- Generate a long-lived access token using Facebook's Access Token Tool

### 3. Database Setup

Start PostgreSQL and initialize the database:

```bash
# Start PostgreSQL (if using Docker)
docker run --name postgres-ig \
  -e POSTGRES_USER=ig_user \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=ig_test \
  -p 5432:5432 -d postgres:15

# Run initialization script
psql -h localhost -U ig_user -d ig_test -f database/init.sql
```

**Important:** Before running `init.sql`, replace `<LONG_LIVED_ACCESS_TOKEN>` with your actual Instagram long-lived access token.

### 4. Install Dependencies

```bash
go mod download
```

### 5. Run the Service

```bash
go run cmd/main.go
```

The service will start on `http://localhost:8080` (or your configured PORT).

### 6. Verify Setup

```bash
# Check if service is ready
curl http://localhost:8080/ready

# Fetch all media
curl http://localhost:8080/media

# Fetch specific media by IDs
curl "http://localhost:8080/media?ids=MEDIA_ID_1,MEDIA_ID_2"
```

---

## 🐳 Docker Deployment

### Build and Run

```bash
# Build image
docker build -t instagram-backend:latest .

# Run container
docker run -d \
  --name instagram-backend \
  --env-file .env \
  -p 8080:8080 \
  instagram-backend:latest
```

---

## 🧪 Testing

### Test Structure

Tests are organized in `tests/` directory:
- `tests/integration/`: Integration tests
- `tests/helpers/`: Test utilities and fixtures
- `tests/dummy/`: Mock Instagram API server

### Available Tests

#### Token Management Tests

| Test | File | Description |
|------|------|-------------|
| **Bootstrap Token Loading** | `bootstrap.test.go` | Verifies token is loaded from PostgreSQL to disk (token.json) and in-memory on service startup |
| **Token Refresh Before Expiry** | `refresh.test.go` | Tests automatic token refresh 7 days before expiration. New token stored in PostgreSQL, disk, and in-memory |

**Manual Verification Tests:**

| Test | How to Verify | Expected Behavior |
|------|---------------|-------------------|
| **Container Restart Recovery** | 1. Start service<br>2. Delete `token.json`<br>3. Restart container | Token reloaded from PostgreSQL, saved to disk, loaded to memory. Service continues normally. |

#### Media Management Tests

| Test | File | Description |
|------|------|-------------|
| **Media Bootstrap** | `media.test.go::TestMediaBootstrap` | Verifies all media is fetched from Instagram and cached in-memory on service startup |
| **New Media Sync** | `media.test.go::TestMediaSyncAddsNewItems` | Tests scheduled sync job that fetches new media from Instagram and updates in-memory cache |
| **Concurrent Access** | `concurrency.test.go` | Validates multiple users can access `/media` endpoint simultaneously with low latency (served from in-memory cache) |

#### Error Handling Tests

| Test | File | Description |
|------|------|-------------|
| **Instagram API Down** | `instagram_down.test.go` | Ensures service remains stable when Instagram API is unavailable |

### Running Tests

```bash
# Run all tests
make test

# Run integration tests only
go test ./tests/integration/...

# Run specific test
go test ./tests/integration/ -run TestBootstrapLoadsFromPostgres

# Run with verbose output
go test -v ./tests/...

# Run with coverage
go test -cover ./...
```

### Pre-Commit Testing

**⚠️ CRITICAL: All commits MUST pass pre-commit tests.**

```bash
make pre-commit
```

This command runs:
1. ✅ Code formatting check (`gofmt`)
2. ✅ Static analysis (`go vet`)
3. ✅ All unit and integration tests
4. ✅ Build verification
5. ✅ Reminds you to update CHANGELOG.md

**CI/CD Pipeline**: GitHub Actions automatically runs these checks on every push and pull request. See [`.github/workflows/pre-commit.yml`](.github/workflows/pre-commit.yml)

**Failed pre-commit?** See the [Contributing Guide](CONTRIBUTING.md#-common-issues) for troubleshooting.
Quick Reference

For detailed instructions, see the [Contributing Guide](CONTRIBUTING.md).

**Required before every commit:**
1. ✅ Run `make pre-commit`
2. ✅ Update `CHANGELOG.md`
3. ✅ Follow commit message format

### Adding New Tests (Summary)

1. Write test in `tests/integration/` or `internal/<module>/`
2. Run: `go test ./tests/integration/ -run YourTest -v`
3. Update README test table (above)
4. Update CHANGELOG.md
5. Run: `make pre-commit`
6. Commit: `git commit -m "test: add test for X"`

**Full details**: See [Contributing Guide - Adding New Tests](CONTRIBUTING.md#-adding-new-tests)

### Commit Message Format

```
<type>: <short description>

[optional body]
```

**Types**: `feat`, `fix`, `test`, `docs`, `refactor`, `chore`, `perf`

**Examples**:
- `feat: add pagination to media endpoint`
- `fix: resolve race condition in token refresh`
- `test: add test for concurrent media access`
- `docs: update setup instructions`

**Full guide**: See [Contributing Guide - Commit Message Format](CONTRIBUTING.md#-commit-message-format)
- `fix`: Bug fix
- `test`: Adding/modifying tests
- `docs`: Documentation changes
- `refactor`: Code refactoring
- `chore`: Maintenance tasks

---

## 📂 Project Structure

```
backend-svc/
├── api/                      # HTTP handlers
│   └── handlers.go          # Media and health check endpoints
├── cmd/
│   └── main.go              # Application entry point
├── database/
│   └── init.sql             # Database schema and seed data
├── internal/
│   ├── bootstrap/           # Service initialization
│   │   └── token.go        # Token bootstrap logic
│   ├── cache/               # In-memory storage
│   │   ├── memory.go       # Cache implementation
│   │   └── memory.test.go  # Cache unit tests
│   ├── config/              # Configuration management
│   │   ├── config.go       # Environment config
│   │   └── postgres.go     # Database connection
│   ├── instagram/           # Instagram API integration
│   │   ├── client.go       # HTTP client wrapper
│   │   ├── models.go       # Data models
│   │   └── service.go      # Business logic
│   ├── scheduler/           # Cron jobs
│   │   ├── cron.go         # Scheduler implementation
│   │   └── cron.test.go    # Scheduler tests
│   └── token/               # Token management
│       ├── disk.go          # Disk persistence
│       ├── postgres.go      # Database operations
│       └── runtime.go       # In-memory token storage
├── tests/
│   ├── dummy/               # Mock servers
│   ├── helpers/             # Test utilities
│   └── integration/         # Integration tests
├── .env.template            # Environment variables template
├── Dockerfile               # Container definition
├── Makefile                 # Build and test commands
├── CHANGELOG.md             # Version history
└── README.md                # This file
```

---

## 🔄 How It Works

### Service Startup Flow

1. **Load Configuration**: Read environment variables from `.env`
2. **Connect to PostgreSQL**: Establish database connection
3. **Bootstrap Token**: 
   - Fetch access token from PostgreSQL
   - Save to disk (`token.json`)
   - Load into in-memory runtime
4. **Bootstrap Media**:
   - Fetch all media from Instagram API
   - Store in in-memory cache
5. **Start Schedulers**:
   - Media sync every 15 minutes
   - Token refresh check daily (refreshes if <7 days to expiry)
6. **Start HTTP Server**: Listen for API requests

### Token Management Flow

```
PostgreSQL (Source of Truth)
    ↓
Disk (token.json - Persistence)
    ↓
In-Memory (Runtime - Fast Access)
```

**Refresh Cycle:**
- Scheduler checks token expiration daily
- If <7 days remaining:
  1. Request new token from Instagram API
  2. Save to PostgreSQL
  3. Save to disk
  4. Update in-memory

### Media Serving Flow

```
User Request → API Handler → In-Memory Cache → JSON Response
```

Benefits:
- No database queries per request
- No Instagram API calls per request
- Sub-millisecond response times
- Handles concurrent users efficiently

---

## 🛠️ Troubleshooting

### Token Issues

**Problem**: Service fails to start with "token expired"
**Solution**: 
1. Generate new long-lived token
2. Update in PostgreSQL: `UPDATE instagram_tokens SET access_token = 'NEW_TOKEN', expires_at = now() + interval '60 days';`
3. Restart service

### Database Connection Issues

**Problem**: "failed to connect postgres"
**Solution**:
1. Verify PostgreSQL is running: `psql -h localhost -U ig_user -d ig_test`
2. Check `DATABASE_URL` in `.env`
3. Ensure database and user exist

### Media Not Syncing

**Problem**: New Instagram posts not appearing
**Solution**:
1. Check logs for Instagram API errors
2. Verify `IG_USER_ID` is correct
3. Confirm token has `instagram_basic` and `pages_read_engagement` permissions
4. Manually trigger sync: restart service

---

## 📊 API Endpoints

| Endpoint | Method | Description | Example |
|----------|--------|-------------|---------|
| `/ready` | GET | Health check | `curl http://localhost:8080/ready` |
| `/media` | GET | Get all media | `curl http://localhost:8080/media` |
| `/media?ids=<ids>` | GET | Get specific media | `curl http://localhost:8080/media?ids=123,456` |

---

## 🎯 Future Plans

- **PostgreSQL Media Storage**: Store media data in PostgreSQL to protect against Instagram API failures and prevent memory bloat
- **Pagination**: Add pagination support for large media collections
- **Webhooks**: Real-time Instagram updates via webhooks
- **Redis Cache**: Optional Redis cache layer for distributed deployments
- **Metrics**: Prometheus metrics for monitoring

---

## 📄 License

[Your License Here]

---

## 👥 Contributors

[Your Team Here]

---

## 📞 Support

For issues or questions:
- Create an issue in the repository
- Contact the development team
- Check the [CHANGELOG.md](CHANGELOG.md) for recent updates