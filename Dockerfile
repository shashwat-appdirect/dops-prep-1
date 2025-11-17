# Multi-stage Dockerfile for AppDirect Workshop Registration SPA

# Stage 1: Frontend Build
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy frontend files
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# Stage 2: Backend Build
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app/backend

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source
COPY backend/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./main.go

# Stage 3: Runtime
FROM alpine:latest

# Install ca-certificates and wget for health check
RUN apk --no-cache add ca-certificates wget

WORKDIR /app

# Copy frontend build from Stage 1
COPY --from=frontend-builder /app/frontend/dist ./static

# Copy backend binary from Stage 2
COPY --from=backend-builder /app/backend/app .

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/registrations/count || exit 1

# Run the application
CMD ["./app"]

