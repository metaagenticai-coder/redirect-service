# Redirect Service

A high-performance URL redirect service optimized for ultra-low latency. This service handles HTTP requests containing short codes and performs HTTP 301 redirects to the original long URLs.

## Overview

The redirect-service is a read-optimized microservice designed to handle 10,000+ requests per second with p99 latency under 50ms. It uses a cache-first architecture with Redis as the primary data source and PostgreSQL read replicas as fallback.

### Key Features

- **Ultra-Low Latency**: p99 latency < 50ms for redirects
- **High Throughput**: Handles 10,000+ requests/second
- **Cache-First Architecture**: >95% cache hit rate target
- **Horizontally Scalable**: Runs on Cloud Run with auto-scaling
- **High Availability**: 99.9% uptime with graceful degradation
- **Minimal Resource Footprint**: 256Mi memory, 1 vCPU

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Chi v5 (lightweight HTTP router)
- **Primary Data Source**: Redis 7.0 (Memorystore)
- **Fallback Data Source**: PostgreSQL 15 (Cloud SQL Read Replica)
- **Container Runtime**: Cloud Run
- **Observability**: Cloud Logging, Cloud Monitoring, Cloud Trace

## Architecture

```
Request → Load Balancer → Redirect Service
                              ↓
                         Redis Cache (primary)
                              ↓ (cache miss)
                         PostgreSQL Read Replica (fallback)
```

## Directory Structure

```
.
├── cmd/
│   └── redirect-service/     # Main application entry point
├── internal/
│   ├── handler/              # HTTP request handlers
│   ├── middleware/           # HTTP middleware (logging, metrics, etc.)
│   ├── service/              # Business logic
│   ├── repository/           # Database access layer
│   ├── cache/                # Redis cache layer
│   └── config/               # Configuration management
├── pkg/
│   ├── logger/               # Structured logging
│   └── metrics/              # Prometheus metrics
├── deployments/              # Deployment configurations
├── tests/                    # Integration tests
├── Dockerfile                # Container image definition
├── Makefile                  # Build and development tasks
└── go.mod                    # Go module dependencies
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker (for containerized development)
- Redis 7.0+ (for local testing)
- PostgreSQL 15+ (for local testing)

### Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

Key environment variables:

- `PORT`: HTTP server port (default: 8081)
- `REDIS_HOST`: Redis connection string
- `REDIS_PASSWORD`: Redis authentication password
- `DATABASE_URL`: PostgreSQL connection string
- `LOG_LEVEL`: Logging level (debug, info, warn, error)

### Local Development

```bash
# Install dependencies
make deps

# Run tests
make test

# Run the service locally
make run

# Build binary
make build
```

### Docker

```bash
# Build Docker image
make docker-build

# Run container
make docker-run
```

## API Endpoints

### GET /{shortCode}

Redirects to the original long URL.

**Response (301 Moved Permanently)**:
```
HTTP/1.1 301 Moved Permanently
Location: https://example.com/original/url
Cache-Control: public, max-age=300
X-Cache-Status: HIT
```

**Response (404 Not Found)**:
```json
{
  "error": "not_found",
  "message": "Short code not found",
  "short_code": "invalid1"
}
```

### GET /health

Liveness probe for container health checks.

**Response (200 OK)**:
```json
{
  "status": "healthy",
  "service": "redirect-service",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### GET /ready

Readiness probe that validates connectivity to dependencies.

**Response (200 OK)**:
```json
{
  "status": "ready",
  "service": "redirect-service",
  "checks": {
    "redis": "ok",
    "database": "ok"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### GET /metrics

Prometheus-compatible metrics endpoint.

## Performance Characteristics

### Latency Targets

- p50: < 10ms (cache hits)
- p95: < 30ms
- p99: < 50ms

### Throughput

- Target: 10,000 requests/second
- Max scaling: 500 Cloud Run instances
- Concurrency: 250 requests per instance

### Cache Performance

- Target cache hit rate: >95%
- Cache TTL: 3600 seconds (1 hour)
- Eviction policy: allkeys-lru

## Deployment

### Cloud Run

Deploy to Google Cloud Run:

```bash
# Build and push image
gcloud builds submit --tag gcr.io/PROJECT_ID/redirect-service

# Deploy to Cloud Run
gcloud run deploy redirect-service \
  --image gcr.io/PROJECT_ID/redirect-service \
  --platform managed \
  --region us-central1 \
  --min-instances 5 \
  --max-instances 500 \
  --memory 256Mi \
  --cpu 1 \
  --concurrency 250 \
  --timeout 10s
```

Or use the deployment configuration:

```bash
gcloud run services replace deployments/cloudrun.yaml
```

## Monitoring and Observability

### Metrics

Key metrics exported to Cloud Monitoring:

- `redirect_requests_total`: Total redirect requests
- `redirect_latency_seconds`: Request latency histogram
- `cache_hit_total`: Cache hits
- `cache_miss_total`: Cache misses
- `cache_hit_rate`: Cache hit rate percentage
- `redis_connection_errors_total`: Redis connection errors
- `db_connection_errors_total`: Database connection errors

### Logging

Structured JSON logs exported to Cloud Logging:

```json
{
  "timestamp": "2024-01-15T10:30:00.123Z",
  "severity": "INFO",
  "service": "redirect-service",
  "method": "GET",
  "path": "/aB3xY9",
  "status": 301,
  "latency_ms": 3.2,
  "cache_hit": true
}
```

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run linter
make lint

# Format code
make fmt
```

## Contributing

1. Follow Go best practices and idioms
2. Write unit tests for new functionality
3. Ensure all tests pass before submitting
4. Use `make fmt` to format code
5. Use `make lint` to check for issues

## License

Copyright © 2024 MetaAgenticAI
