# AppDirect India AI Workshop Registration SPA

A production-ready React SPA with Golang backend for event registration and management.

## Project Structure

```
.
├── frontend/          # React + Vite + TypeScript application
├── backend/           # Golang REST API
├── docker-compose.yml # Local development setup
├── Dockerfile         # Multi-stage Dockerfile for production
├── Makefile           # Build and test automation
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
- For Cloud Run: Google Cloud Project with Firestore API enabled

## Environment Setup

### Backend

Create `backend/.env`:

```env
# For local development (required)
FIREBASE_SERVICE_ACCOUNT=/path/to/service-account.json
# OR base64 encoded:
# FIREBASE_SERVICE_ACCOUNT=base64:eyJ0eXAiOiJKV1QiLCJhbGc...

# For Cloud Run (optional - uses ADC)
# FIREBASE_SERVICE_ACCOUNT can be omitted

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

### Using Makefile

```bash
# Install dependencies
make deps

# Run tests
make test
make test-unit
make test-integration

# Build
make build          # Backend only
make build-frontend # Frontend only
make build-all      # Both

# Run locally
make run            # Backend
make run-frontend   # Frontend

# Docker
make docker-build
make docker-run
```

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

### Using Docker

```bash
# Build image
docker build -t appdirect-workshop:latest .

# Run container
docker run -p 8080:8080 \
  -e SUBSCOLLECTION_ID=workshop-2024 \
  -e ADMIN_PASSWORD=your-password \
  -e PORT=8080 \
  -e CORS_ORIGIN=https://your-domain.com \
  appdirect-workshop:latest
```

### Manual Build

#### Backend

```bash
cd backend
go build -o app main.go
./app
```

#### Frontend

```bash
cd frontend
npm run build
# Serve dist/ directory with your preferred web server
```

## Google Cloud Run Deployment

The application is configured to use Application Default Credentials (ADC) on Cloud Run, eliminating the need for service account files.

### Prerequisites

1. Google Cloud Project with Firestore API enabled
2. Cloud Run API enabled
3. Service account with Firestore permissions

### Deployment Steps

1. **Build and push Docker image:**

```bash
# Set your project ID
export PROJECT_ID=your-project-id
export SERVICE_NAME=appdirect-workshop

# Build and push
gcloud builds submit --tag gcr.io/$PROJECT_ID/$SERVICE_NAME

# Or use Artifact Registry
gcloud builds submit --tag $REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$SERVICE_NAME
```

2. **Deploy to Cloud Run:**

```bash
gcloud run deploy $SERVICE_NAME \
  --image gcr.io/$PROJECT_ID/$SERVICE_NAME \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars SUBSCOLLECTION_ID=workshop-2024,ADMIN_PASSWORD=your-secure-password,CORS_ORIGIN=https://your-service-url.run.app \
  --service-account your-service-account@your-project.iam.gserviceaccount.com
```

3. **Set IAM permissions:**

The Cloud Run service account needs Firestore access:
- Cloud Datastore User
- Or custom role with Firestore read/write permissions

### Environment Variables for Cloud Run

- `SUBSCOLLECTION_ID` - Required: Firestore subcollection identifier
- `ADMIN_PASSWORD` - Required: Admin dashboard password
- `PORT` - Optional: Cloud Run sets this automatically
- `CORS_ORIGIN` - Required: Your Cloud Run service URL
- `GOOGLE_CLOUD_PROJECT` - Optional: Auto-detected on Cloud Run
- `K_SERVICE` - Optional: Auto-set by Cloud Run (triggers ADC mode)

**Note:** `FIREBASE_SERVICE_ACCOUNT` is NOT required on Cloud Run - the application uses Application Default Credentials automatically.

## Testing

### Run All Tests

```bash
make test
```

### Run Specific Test Suites

```bash
# Unit tests
make test-unit

# Integration tests
make test-integration

# Individual packages
cd backend
go test ./internal/handlers -v
go test ./internal/middleware -v
go test ./internal/config -v
```

## Security Notes

- Never commit `.env` files or service account JSON files
- Use environment variables for all sensitive data
- Admin password should be strong and kept secure
- CORS origin should be configured for production
- On Cloud Run, use IAM service accounts instead of service account files

## API Documentation

### Public Endpoints

- `POST /api/register` - Register for event
- `GET /api/registrations/count` - Get registration count
- `GET /api/speakers` - List speakers
- `GET /api/sessions` - List sessions

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

## Health Check

The Dockerfile includes a health check endpoint:
- `GET /api/registrations/count` - Used for container health checks
