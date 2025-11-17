# Multi-stage Dockerfile for building frontend and backend into one image

# Stage 1: Frontend Build
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files
COPY frontend/package.json frontend/package-lock.json* ./

# Install dependencies
RUN npm ci

# Copy frontend source
COPY frontend/ .

# Build frontend
RUN npm run build

# Stage 2: Backend Build
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app/backend

# Install git for go modules
RUN apk add --no-cache git

# Copy go mod files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy backend source
COPY backend/ .

# Build backend binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./main.go

# Stage 3: Runtime
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy frontend build from Stage 1
COPY --from=frontend-builder /app/frontend/dist ./static

# Copy backend binary from Stage 2
COPY --from=backend-builder /app/backend/app .

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/registrations/count || exit 1

# Run backend
CMD ["./app"]

