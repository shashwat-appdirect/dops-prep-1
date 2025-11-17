# AppDirect India AI Workshop Registration SPA

A production-ready React SPA with Golang backend for event registration and management.

## Project Structure

```
.
├── frontend/          # React + Vite + TypeScript application
├── backend/           # Golang REST API
├── docker-compose.yml # Local development setup
└── README.md
```

## Features

- **Hero Section** with event branding and CTA buttons
- **Sessions & Speakers** grid display
- **Registration Form** with live counter and success confirmation
- **Location** section with embedded Google Maps
- **Admin Dashboard** (password-protected) with:
  - Attendee management
  - Speaker CRUD operations
  - Session CRUD operations
  - Analytics with designation breakdown pie chart

## Prerequisites

- Node.js 18+ and npm/yarn
- Go 1.21+
- Docker and Docker Compose (for containerized development)
- Firebase project with Firestore enabled
- Firebase service account JSON file

## Environment Setup

### Backend

Create `backend/.env`:

```env
FIREBASE_SERVICE_ACCOUNT=/path/to/service-account.json
# OR base64 encoded:
# FIREBASE_SERVICE_ACCOUNT=base64:eyJ0eXAiOiJKV1QiLCJhbGc...
SUBSCOLLECTION_ID=workshop-2024
ADMIN_PASSWORD=your-secure-password
PORT=8080
CORS_ORIGIN=http://localhost:5173
```

### Frontend

Create `frontend/.env`:

```env
VITE_API_URL=http://localhost:8080
```

## Development

### Using Docker Compose

```bash
docker-compose up --build
```

Frontend: http://localhost:5173
Backend: http://localhost:8080

### Manual Setup

#### Backend

```bash
cd backend
go mod download
go run main.go
```

#### Frontend

```bash
cd frontend
npm install
npm run dev
```

## Production Build

### Using Docker (Recommended)

```bash
# Build Docker image
make docker-build
# or
docker build -t appdirect-workshop:latest .

# Run locally
make docker-run
# or with custom env vars
docker run -p 8080:8080 \
  -e FIREBASE_SERVICE_ACCOUNT="base64:..." \
  -e SUBSCOLLECTION_ID="workshop-2024" \
  -e ADMIN_PASSWORD="your-password" \
  -e PORT=8080 \
  -e CORS_ORIGIN="http://localhost:8080" \
  appdirect-workshop:latest
```

### Manual Build

#### Backend

```bash
cd backend
go build -o app
./app
```

#### Frontend

```bash
cd frontend
npm run build
# Serve dist/ directory with your preferred web server
```

## Testing

See [TESTING.md](TESTING.md) for detailed testing information.

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration
```

## Google Cloud Run Deployment

The application is ready for Cloud Run deployment:

1. Build and push Docker image to Google Container Registry or Artifact Registry
2. Deploy to Cloud Run with environment variables:
   - `FIREBASE_SERVICE_ACCOUNT` (base64 encoded JSON)
   - `SUBSCOLLECTION_ID`
   - `ADMIN_PASSWORD`
   - `PORT` (set automatically by Cloud Run)
   - `CORS_ORIGIN` (your Cloud Run URL)

The Dockerfile includes:
- Multi-stage build for optimized image size
- Frontend and backend combined in single image
- Static file serving for SPA
- Health check endpoint
- Non-root user for security

## Security Notes

- Never commit `.env` files or service account JSON files
- Use environment variables for all sensitive data
- Admin password should be strong and kept secure
- CORS origin should be configured for production

## API Documentation

### Public Endpoints

- `POST /api/register` - Register for event
- `GET /api/registrations/count` - Get registration count

### Admin Endpoints (require authentication)

- `POST /api/admin/login` - Admin login
- `GET /api/admin/attendees` - List attendees
- `GET /api/admin/attendees/:id` - Get attendee details
- `GET /api/admin/speakers` - List speakers
- `POST /api/admin/speakers` - Create speaker
- `PUT /api/admin/speakers/:id` - Update speaker
- `DELETE /api/admin/speakers/:id` - Delete speaker
- `GET /api/admin/sessions` - List sessions
- `POST /api/admin/sessions` - Create session
- `PUT /api/admin/sessions/:id` - Update session
- `DELETE /api/admin/sessions/:id` - Delete session
- `GET /api/admin/analytics/designations` - Get designation breakdown

