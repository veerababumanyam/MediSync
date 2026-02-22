# MediSync ETL Service Dockerfile
# Multi-stage build for optimized production image
#
# This Dockerfile builds the ETL service that:
#   - Syncs data from Tally ERP and HIMS to PostgreSQL data warehouse
#   - Publishes events to NATS JetStream
#   - Caches data in Redis
#   - Exports metrics for Prometheus
#
# Build:
#   docker build -t medisync-etl -f cmd/etl/Dockerfile .
#
# Run:
#   docker run --env-file .env medisync-etl

# ============================================================================
# Stage 1: Builder
# ============================================================================
# Using official Go image for building
FROM golang:1.26-alpine AS builder

# Install build dependencies
# - git: for fetching Go modules
# - ca-certificates: for HTTPS connections
# - tzdata: for timezone support
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata

# Set working directory
WORKDIR /app

# Copy go module files first (better layer caching)
COPY go.mod go.sum ./

# Download dependencies
# This layer is cached unless go.mod or go.sum changes
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the ETL binary
# CGO_ENABLED=0: Produce statically linked binary
# -ldflags: Strip debug info for smaller binary
# -trimpath: Remove file system paths for reproducible builds
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev') -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -trimpath \
    -o /app/bin/etl \
    ./cmd/etl

# ============================================================================
# Stage 2: Runtime
# ============================================================================
# Using distroless for minimal attack surface
# Contains: ca-certificates, /etc/passwd (nonroot user), tzdata
FROM gcr.io/distroless/static-debian12:nonroot

# Labels for container metadata
LABEL org.opencontainers.image.title="MediSync ETL Service"
LABEL org.opencontainers.image.description="ETL service for syncing Tally ERP and HIMS data to PostgreSQL warehouse"
LABEL org.opencontainers.image.vendor="MediSync"
LABEL org.opencontainers.image.source="https://github.com/medisync/medisync"
LABEL com.medisync.service="etl"

# Copy timezone data for proper time handling
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy CA certificates for HTTPS connections
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary from builder stage
COPY --from=builder /app/bin/etl /etl

# Copy migrations directory (needed for embedded migrations if used)
# Note: Migrations are typically run separately via migrate CLI,
# but including them allows the ETL service to verify schema if needed
COPY --from=builder /app/migrations /migrations

# Set working directory
WORKDIR /

# Environment variables with defaults
# These can be overridden at runtime via docker run --env or docker-compose
ENV TZ=UTC
ENV LOG_LEVEL=info
ENV LOG_FORMAT=json

# Health check
# The ETL service should expose a /health endpoint on the metrics port
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD ["/etl", "health"]

# Expose metrics port (Prometheus scraping)
EXPOSE 9100

# Run as non-root user (provided by distroless/nonroot)
USER nonroot:nonroot

# Run the ETL service
ENTRYPOINT ["/etl"]

# Default command (can be overridden)
# The service runs in daemon mode by default
CMD ["serve"]
