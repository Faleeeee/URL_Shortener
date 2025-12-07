# URL Shortener Service

A production-ready URL shortening service built with Go, similar to bit.ly or TinyURL. This service converts long URLs into short, shareable links with click tracking and comprehensive API support.

[![Go Version](https://img.shields.io/badge/Go-1.23-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

---

## 📋 Table of Contents

- [Problem Description](#problem-description)
- [Features](#features)
- [Quick Start](#quick-start)
- [API Documentation](#api-documentation)
- [Architecture & Design Decisions](#architecture--design-decisions)
- [Technical Trade-offs](#technical-trade-offs)
- [Challenges & Solutions](#challenges--solutions)
- [Testing](#testing)
- [Limitations & Future Improvements](#limitations--future-improvements)
- [Production Readiness](#production-readiness)

---

## 🎯 Problem Description

### The Challenge

Users have long URLs like:
```
https://example.com/very/long/path/to/resource?param1=value1&param2=value2
```

And want to shorten them to:
```
http://short.url/abc123
```

### Requirements

1. **Create Short URLs**: Convert any long URL into a compact, shareable link
2. **Redirect**: Automatically redirect users from short URL to original URL
3. **Analytics**: Track click counts for each shortened URL
4. **Management**: List and retrieve information about created URLs
5. **Custom Aliases**: Optionally allow users to specify custom short codes

---

## ✨ Features

### Core Functionality
- ✅ **URL Shortening**: Generate short, random 6-character aliases
- ✅ **Custom Aliases**: Support for user-defined short codes
- ✅ **Fast Redirects**: 302 redirects with async click tracking
- ✅ **Click Analytics**: Real-time counter increments
- ✅ **Pagination**: Efficient listing of all URLs

### Security & Validation
- ✅ **URL Validation**: Format checking with regex
- ✅ **Private URL Blocking**: Prevents localhost and private IP addresses
- ✅ **Input Sanitization**: Alphanumeric-only aliases
- ✅ **Collision Handling**: Automatic retry with new codes

### Performance
- ✅ **Database Indexes**: Unique index on alias for O(1) lookups
- ✅ **Connection Pooling**: Configured for 25 max connections
- ✅ **Atomic Operations**: Race-condition-free click counting

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.23+**: [Install Go](https://golang.org/doc/install)
- **Docker & Docker Compose**: [Install Docker](https://docs.docker.com/get-docker/)
- **PostgreSQL** (or use Docker Compose)

### 1. Clone the Repository

```bash
git clone <your-repo-url>
cd URL-Shortener-Service
```

### 2. Configure Environment Variables

```bash
# Copy the example environment file
cp .env.example .env

# Edit .env with your database credentials if needed
# The default values work with the Docker Compose setup
```

### 3. Start the Database

```bash
# Using Docker Compose
sudo docker compose up -d

# Wait for database to be ready (5 seconds)
sleep 5
```

### 4. Run Database Migrations

```bash
sudo docker exec url_shortener_db psql -U postgres -d url_shortener -f /migrations/000001_create_urls_table.up.sql
```

### 5. Install Dependencies

```bash
go mod download
```

### 6. Run the Service

```bash
go run cmd/api/main.go
```

The service will start on `http://localhost:8080`

### 7. Access Swagger Documentation

Open your browser to:
```
http://localhost:8080/swagger/index.html
```

---

## ⚙️ Configuration

The service uses environment variables for configuration. All settings are defined in the `.env` file.

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SERVER_PORT` | Port the server listens on | `8080` | No |
| `DATABASE_URL` | PostgreSQL connection string | - | Yes |
| `JWT_SECRET` | Secret key for JWT token signing | - | Yes |
| `JWT_EXPIRATION` | JWT token expiration duration | `24h` | No |

### Example `.env` File

```bash
# Server Configuration
SERVER_PORT=8080

# Database Configuration
DATABASE_URL=postgres://postgres:123456@localhost:5432/url_shortener?sslmode=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION=24h
```

### Configuration Loading

The service loads configuration in the following order:
1. Reads from `.env` file in the project root
2. Falls back to system environment variables if `.env` not found
3. Uses default values for optional settings
4. Fails with clear error message if required variables are missing

> [!IMPORTANT]
> **Production Security**: Always use strong, randomly generated values for `JWT_SECRET` in production environments.

---

## 📡 API Documentation

### Base URL
```
http://localhost:8080
```

### Endpoints

#### 1. Create Short URL

**POST** `/url/shorten`

Create a shortened URL with optional custom alias.

**Request Body:**
```json
{
  "url": "https://www.example.com/very/long/url",
  "alias": "my-link"  // Optional custom alias
}
```

**Response (200 OK):**
```json
{
  "alias": "my-link",
  "short_url": "http://localhost:8080/my-link",
  "original_url": "https://www.example.com/very/long/url"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid URL format, URL too long, or invalid alias
- `409 Conflict`: Custom alias already exists
- `500 Internal Server Error`: Failed to create short URL

**cURL Example:**
```bash
curl -X POST http://localhost:8080/url/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'
```

---

#### 2. Redirect to Original URL

**GET** `/:alias`

Redirects to the original URL and increments click counter.

**Parameters:**
- `alias` (path): The short code

**Response:**
- `302 Found`: Redirects to original URL
- `404 Not Found`: Short URL not found

**Browser Example:**
```
http://localhost:8080/my-link
```

---

#### 3. Get URL Information

**GET** `/url/links/:alias`

Retrieve metadata about a shortened URL.

**Response (200 OK):**
```json
{
  "alias": "my-link",
  "original_url": "https://www.example.com",
  "click_count": 42,
  "created_at": "2025-12-05T15:30:00Z",
  "updated_at": "2025-12-05T16:45:00Z"
}
```

**cURL Example:**
```bash
curl http://localhost:8080/url/links/my-link
```

---

#### 4. List All URLs

**GET** `/url/links`

Retrieve a paginated list of all shortened URLs.

**Query Parameters:**
- `limit` (optional): Items per page (default: 50, max: 100)
- `offset` (optional): Offset for pagination (default: 0)

**Response (200 OK):**
```json
{
  "urls": [
    {
      "id": 1,
      "alias": "abc123",
      "original_url": "https://www.google.com",
      "click_count": 5,
      "created_at": "2025-12-05T15:30:00Z",
      "updated_at": "2025-12-05T15:35:00Z"
    }
  ],
  "count": 1,
  "limit": 50,
  "offset": 0
}
```

**cURL Example:**
```bash
curl "http://localhost:8080/url/links?limit=20&offset=0"
```

---

## 🏗️ Architecture & Design Decisions

### 1. Database Choice: **PostgreSQL**

#### Why PostgreSQL?

✅ **Chosen for:**
- **ACID Compliance**: Ensures data integrity for concurrent requests
- **Unique Constraints**: Database-level prevention of duplicate aliases
- **Atomic Operations**: `UPDATE ... SET count = count + 1` prevents race conditions
- **Indexing**: Fast O(1) lookups on alias column
- **Transactions**: Support for multi-step operations
- **Reliability**: Battle-tested in production environments

❌ **Trade-offs vs. NoSQL (MongoDB, Redis, DynamoDB):**
- **Vertical Scaling Limitation**: PostgreSQL scales vertically (bigger machines) while NoSQL scales horizontally (more machines)
- **Complexity**: Requires more setup than embedded databases like SQLite
- **Cost**: More expensive than serverless options for low traffic

#### Mitigation Strategy:
- Add **read replicas** for horizontal read scaling
- Implement **Redis caching** for hot URLs (80/20 rule)
- Use **connection pooling** to maximize throughput
- Consider **database sharding** for extreme scale (100M+ URLs)

---

### 2. Short Code Generation: **Base62 + Cryptographic Random**

#### Algorithm

```go
Characters: [0-9A-Za-z] = 62 possibilities
Length: 6 characters
Total combinations: 62^6 = 56,800,235,584 (56.8 billion)
```

#### Why This Approach?

✅ **Advantages:**
- **High Collision Resistance**: 56.8 billion combinations ensure virtually no collisions
- **Compact**: Only 6 characters (user-friendly)
- **Unpredictable**: Cryptographic randomness prevents URL guessing
- **Stateless**: No need for distributed counter synchronization

❌ **Alternatives Considered:**

| Approach | Why Not Chosen |
|----------|----------------|
| **Auto-increment ID + base62** | Predictable (security risk), exposes URL count |
| **MD5/SHA hash + truncate** | Collision possible, longer codes (8-10 chars) |
| **Snowflake ID** | Requires distributed coordination, overkill |
| **UUID** | Too long (36 chars) for "short" URL |

#### Collision Handling

```go
1. Generate random 6-character base62 code
2. Attempt INSERT into database
3. If unique constraint violation → retry (max 3 times)
4. If still fails → return error (astronomically rare)
```

**Collision Probability**: With 1 million URLs, probability ≈ 0.001% (negligible)

---

### 3. API Design: **REST**

#### Why REST over GraphQL/gRPC?

✅ **REST chosen because:**
- **Simplicity**: Well-understood by all developers
- **Perfect fit**: CRUD operations map naturally to HTTP methods
- **Caching**: Browser and CDN caching works out-of-the-box
- **Redirects**: Native HTTP 302 redirect support
- **Tooling**: Swagger/OpenAPI for documentation

❌ **GraphQL**: Overkill for simple CRUD, no redirect support
❌ **gRPC**: Requires protobuf, no browser support without proxy

---

### 4. Concurrency Strategy

#### Problem: Race Conditions

**Scenario 1**: Two users create URLs simultaneously
**Solution**: Database unique constraint on `alias` column

```sql
CREATE UNIQUE INDEX idx_alias ON urls(alias);
```

**Scenario 2**: Multiple click events for same URL
**Solution**: Atomic SQL update

```sql
UPDATE urls SET click_count = click_count + 1 WHERE alias = ?
```

**Scenario 3**: Read-modify-write collision
**Solution**: Use `QueryRow` + `Exec` with transactions

---

### 5. Project Structure: **Clean Architecture**

```
├── cmd/api/              # Application entry point
├── internal/
│   ├── domain/          # Business entities & validation
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic
│   ├── handler/         # HTTP handlers (API)
│   ├── server/          # Server & routing
│   ├── database/        # Database connection
│   └── config/          # Configuration loading from env
├── migrations/          # SQL migrations
├── docs/               # Swagger documentation
├── .env.example        # Environment variables template
└── .env                # Environment variables (gitignored)
```

**Benefits:**
- **Separation of Concerns**: Each layer has single responsibility
- **Testability**: Easy to mock dependencies
- **Maintainability**: Changes isolated to specific layers
- **Scalability**: Can split into microservices later

---

## ⚖️ Technical Trade-offs

### Trade-off 1: 302 (Temporary) vs 301 (Permanent) Redirect

**Choice**: 302 Temporary Redirect

✅ **Why 302:**
- Browsers always request server (increments counter accurately)
- Original URL can be changed if needed
- Better for analytics and tracking

❌ **Why NOT 301:**
- Browsers cache permanently (bypasses server)
- Cannot update original URL
- Click counter would be inaccurate

---

### Trade-off 2: Synchronous vs Asynchronous Click Counting

**Choice**: Asynchronous (fire-and-forget)

```go
go h.service.IncrementClickCount(alias)
c.Redirect(http.StatusFound, url.OriginalURL)
```

✅ **Advantages:**
- **Fast redirects**: User doesn't wait for counter update
- **Better UX**: Sub-millisecond response times

❌ **Disadvantages:**
- **Potential data loss**: If server crashes before update
- **Eventual consistency**: Counter might lag slightly

**Mitigation**: PostgreSQL write-ahead log ensures durability even if goroutine fails

---

### Trade-off 3: Pagination Limit (Max 100)

**Choice**: Hard cap at 100 items per page

✅ **Why:**
- Prevents abuse (fetching millions of records)
- Protects database from expensive queries
- Reduces network payload

❌ **Trade-off:**
- Users need multiple requests for large datasets

**Alternative**: Could implement cursor-based pagination for better performance

---

### Trade-off 4: No Rate Limiting (Yet)

**Decision**: Not implemented in v1

✅ **Why deferred:**
- Adds complexity (Redis, token bucket algorithm)
- YAGNI (You Ain't Gonna Need It) for MVP
- Can add later via middleware

⚠️ **Risk:**
- Vulnerable to abuse (spam, DDOS)

**Mitigation Plan**:
- Use reverse proxy (Nginx) with `limit_req`
- Implement API key system
- Add Redis-based rate limiter middleware

---

## 🔥 Challenges & Solutions

### Challenge 1: Concurrent URL Creation

**Problem**: Two requests with same long URL arrive simultaneously

**Solution**: Retry logic with exponential backoff
```go
for i := 0; i < MaxRetries; i++ {
    alias := GenerateShortCode()
    if err := repo.Create(alias); err == nil {
        return alias, nil
    }
    // On duplicate, retry
}
```

**Alternative considered**: Check if URL exists first → Race condition still possible

---

### Challenge 2: Private URL Prevention

**Problem**: User could shorten `http://localhost:9090/admin` and share it

**Solution**: Blacklist common private patterns
```go
if strings.Contains(host, "localhost") ||
   strings.HasPrefix(host, "127.") ||
   strings.HasPrefix(host, "192.168.") { ... }
```

**Limitation**: Doesn't catch all private ranges (e.g., `172.16-31.x.x`)

**Future**: Use CIDR matching library for comprehensive check

---

### Challenge 3: Click Counter Race Conditions

**Problem**: Multiple clicks → lost updates

**Bad Approach** (race condition):
```go
url := repo.FindByAlias(alias)
url.ClickCount++
repo.Update(url)  // Lost update!
```

**Good Approach** (atomic):
```sql
UPDATE urls SET click_count = click_count + 1 WHERE alias = ?
```

**Learning**: Always use atomic operations for counters

---

### Challenge 4: Swagger Code Generation

**Problem**: Swagger docs out of sync with code

**Solution**: Use `swag` annotations in code
```go
// @Summary Create a shortened URL
// @Param request body domain.ShortenRequest true "URL to shorten"
func (h *URLHandler) ShortenURL(c *gin.Context) { ... }
```

Then auto-generate:
```bash
swag init -g cmd/api/main.go
```

**Benefit**: Single source of truth (code)

---

## 🧪 Testing

### Run All Tests

```bash
go test -v ./...
```

### Run Tests with Coverage

```bash
go test -v -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Test Results

```
✅ Short Code Generation:
   - Uniqueness test (1000 iterations): PASS
   - Length validation: PASS
   - Base62 character check: PASS

✅ URL Validation:
   - Valid HTTP/HTTPS: PASS
   - Localhost blocking: PASS
   - Private IP blocking (192.168, 10, 172.16): PASS
   - URL length limit: PASS

✅ Alias Validation:
   - Alphanumeric + hyphen + underscore: PASS
   - Special character rejection: PASS
   - Length limit: PASS
```

### Manual Testing

```bash
# Create short URL
curl -X POST http://localhost:8080/url/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'

# Response: {"alias":"a1B2c3","short_url":"http://localhost:8080/a1B2c3", ...}

# Test redirect
curl -L http://localhost:8080/a1B2c3

# Get URL info
curl http://localhost:8080/url/links/a1B2c3

# List all URLs
curl http://localhost:8080/url/links
```

---

## 🚧 Limitations & Future Improvements

### Current Limitations

| Limitation | Impact | Priority |
|------------|--------|----------|
| **No Rate Limiting** | Vulnerable to abuse | HIGH |
| **No URL Expiration** | Database grows indefinitely | MEDIUM |
| **No Analytics Dashboard** | Limited insights | LOW |
| **No Custom Domains** | Only localhost:8080 | LOW |
| **No URL Validation for Malicious Sites** | Phishing risk | MEDIUM |

### Future Improvements

#### Phase 1: Security & Reliability
- [ ] **Rate Limiting**: 100 requests/hour per IP
- [ ] **API Keys**: Authentication for paid tiers
- [ ] **URL Blacklist**: Block known malicious domains
- [ ] **HTTPS Support**: TLS certificates via Let's Encrypt

#### Phase 2: Features
- [ ] **QR Code Generation**: Auto-generate QR codes for short URLs
- [ ] **Expiration**: Auto-delete after N days/clicks
- [ ] **Password Protection**: Secure short URLs with password
- [ ] **Custom Domains**: Support `go.yourcompany.com`

#### Phase 3: Analytics
- [ ] **Click Analytics**: Track user agent, referrer, geo-location
- [ ] **Admin Dashboard**: Web UI for URL management
- [ ] **Real-time Stats**: WebSocket for live click updates

#### Phase 4: Scale
- [ ] **Redis Caching**: Cache hot URLs (80/20 rule)
- [ ] **Read Replicas**: Scale PostgreSQL reads
- [ ] **CDN Integration**: Cloudflare for global redirects
- [ ] **Database Sharding**: Partition by hash(alias)

---

## 🏭 Production Readiness

### What's Missing for Production?

| Requirement | Status | Solution |
|-------------|--------|----------|
| **SSL/TLS** | ❌ Not implemented | Use Nginx reverse proxy + Let's Encrypt |
| **Monitoring** | ❌ Not implemented | Add Prometheus + Grafana |
| **Logging** | ⚠️ Basic only | Integrate Logrus/Zap with structured logging |
| **Error Tracking** | ❌ Not implemented | Sentry or Rollbar integration |
| **CI/CD** | ❌ Not implemented | GitHub Actions for test + deploy |
| **Load Balancer** | ❌ Not implemented | Nginx or AWS ALB |
| **Database Backups** | ⚠️ Manual | Automated daily backups to S3 |
| **Health Checks** | ✅ Implemented | `/health` endpoint |

### Deployment Architecture (Proposed)

```
Internet
   ↓
Cloudflare CDN (DDoS protection, caching)
   ↓
Nginx Load Balancer (SSL termination, rate limiting)
   ↓
Go Service (3 replicas, Docker containers)
   ↓
PostgreSQL Primary + 2 Read Replicas
   ↓
Redis Cache (hot URL cache)
```

### Estimated Capacity

**Current Setup (Single Instance):**
- **RPS**: ~5,000 requests/second
- **URLs**: 56 billion (limited by 6-char base62 space)
- **Database**: 100M URLs ≈ 20 GB storage

**With Scaling (Horizontal + Caching):**
- **RPS**: 50,000+ requests/second
- **Cost**: ~$500/month (AWS t3.medium x3 + RDS + ElastiCache)

---

## 🛠️ Development Commands

### Makefile Commands

```bash
make help              # Show all available commands
make build             # Build binary to bin/urlshortener
make run               # Run the service locally
make test              # Run all tests
make test-coverage     # Generate coverage report
make swagger           # Regenerate Swagger docs
make docker-up         # Start Docker containers
make docker-down       # Stop Docker containers
make migrate-up        # Run database migrations
make migrate-down      # Rollback migrations
make clean             # Remove build artifacts
```

---

## 📊 Database Schema

```sql
CREATE TABLE urls (
    id BIGSERIAL PRIMARY KEY,
    alias VARCHAR(16) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    click_count BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_alias ON urls(alias);
CREATE INDEX idx_created_at ON urls(created_at);
```

**Index Strategy:**
- `idx_alias`: Unique index for O(1) alias lookups
- `idx_created_at`: For analytics queries (newest URLs first)

---

## 🤝 Contributing

This is an assignment project, but contributions are welcome for:
- Bug fixes
- Performance improvements
- Documentation improvements

---

## 📄 License

MIT License - feel free to use for learning or commercial projects.

---

## 🙏 Acknowledgments

Built as a technical assignment showcasing:
- Clean architecture in Go
- RESTful API design
- Database optimization
- Concurrency handling
- Production-ready thinking

**Technologies Used:**
- Go 1.23
- Gin Web Framework
- PostgreSQL 16
- Swagger/OpenAPI
- Docker & Docker Compose