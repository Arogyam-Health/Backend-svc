# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive README with setup instructions, architecture overview, and testing guidelines
- CHANGELOG.md for tracking project changes
- Makefile with pre-commit test automation
- Documentation for all existing tests (Token Management and Media Management)

### Changed
- N/A

### Fixed
- N/A

---

## [1.0.0] - 2026-01-02

### Added
- Initial service implementation with Instagram media sync
- Multi-layer token management (PostgreSQL, disk, in-memory)
- In-memory cache for media storage
- Automatic token refresh 7 days before expiration
- Scheduled media synchronization with Instagram API
- Integration test suite:
  - Token bootstrap test
  - Token refresh test  
  - Media bootstrap test
  - Media sync test
  - Concurrent access test
  - Instagram API failure handling test
- REST API endpoints:
  - `GET /ready` - Health check
  - `GET /media` - Retrieve all media
  - `GET /media?ids=<ids>` - Retrieve specific media by IDs
- Docker support with multi-stage builds
- PostgreSQL schema with token storage table

### Changed
- N/A

### Fixed
- N/A

---

## How to Use This Changelog

### For Developers

**Every commit that adds, changes, or fixes functionality must update this file.**

#### Format for New Entries

Add your changes under the `[Unreleased]` section in the appropriate category:

```markdown
## [Unreleased]

### Added
- New feature or functionality
- New test case

### Changed
- Modified existing behavior
- Updated dependencies

### Fixed
- Bug fixes
- Test fixes
```

#### Example Entry

```markdown
## [Unreleased]

### Added
- Added pagination support for media endpoint
- Added TestMediaPagination integration test

### Changed
- Updated cache implementation to support pagination
- Increased default page size to 50 items

### Fixed
- Fixed race condition in token refresh scheduler
```

#### When Creating a Release

1. Move all `[Unreleased]` items to a new version section
2. Add the version number and date
3. Create a new empty `[Unreleased]` section

Example:
```markdown
## [Unreleased]

### Added
- (empty for now)

## [1.1.0] - 2026-01-15

### Added
- Pagination support for media endpoint
- TestMediaPagination integration test
```

---

**Remember**: Keep entries concise, clear, and user-focused. Describe WHAT changed, not HOW it was implemented.
