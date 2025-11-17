# Deployment Guide

## Google Cloud Run Deployment

This application is ready for deployment to Google Cloud Run.

### Prerequisites

1. Google Cloud Project with billing enabled
2. Firebase project with Firestore enabled
3. Firebase service account JSON file
4. Docker installed locally (for building image)

### Step 1: Build and Push Docker Image

```bash
# Set your project ID
export PROJECT_ID=your-project-id
export REGION=us-central1
export SERVICE_NAME=appdirect-workshop

# Build the image
docker build -t gcr.io/${PROJECT_ID}/${SERVICE_NAME}:latest .

# Push to Google Container Registry
docker push gcr.io/${PROJECT_ID}/${SERVICE_NAME}:latest
```

Or use Artifact Registry:

```bash
# Create Artifact Registry repository
gcloud artifacts repositories create ${SERVICE_NAME} \
  --repository-format=docker \
  --location=${REGION}

# Configure Docker to use gcloud
gcloud auth configure-docker ${REGION}-docker.pkg.dev

# Build and push
docker build -t ${REGION}-docker.pkg.dev/${PROJECT_ID}/${SERVICE_NAME}/${SERVICE_NAME}:latest .
docker push ${REGION}-docker.pkg.dev/${PROJECT_ID}/${SERVICE_NAME}/${SERVICE_NAME}:latest
```

### Step 2: Prepare Firebase Service Account

Encode your Firebase service account JSON to base64:

```bash
# Option 1: Base64 encode the file
cat firebase-service-account.json | base64 | tr -d '\n' > service-account-base64.txt

# Option 2: Use the base64 string directly
FIREBASE_SA_B64=$(cat firebase-service-account.json | base64 | tr -d '\n')
```

### Step 3: Deploy to Cloud Run

```bash
gcloud run deploy ${SERVICE_NAME} \
  --image gcr.io/${PROJECT_ID}/${SERVICE_NAME}:latest \
  --platform managed \
  --region ${REGION} \
  --allow-unauthenticated \
  --set-env-vars="FIREBASE_SERVICE_ACCOUNT=base64:${FIREBASE_SA_B64}" \
  --set-env-vars="SUBSCOLLECTION_ID=workshop-2024" \
  --set-env-vars="ADMIN_PASSWORD=your-secure-password" \
  --set-env-vars="CORS_ORIGIN=https://${SERVICE_NAME}-xxxxx-${REGION}.a.run.app" \
  --memory 512Mi \
  --cpu 1 \
  --min-instances 0 \
  --max-instances 10 \
  --timeout 300
```

**Note:** Update `CORS_ORIGIN` with your actual Cloud Run URL after first deployment.

### Step 4: Update CORS Origin

After deployment, get your service URL:

```bash
SERVICE_URL=$(gcloud run services describe ${SERVICE_NAME} --region=${REGION} --format='value(status.url)')
echo "Service URL: ${SERVICE_URL}"
```

Then update the CORS_ORIGIN:

```bash
gcloud run services update ${SERVICE_NAME} \
  --region ${REGION} \
  --update-env-vars="CORS_ORIGIN=${SERVICE_URL}"
```

### Environment Variables

Required environment variables for Cloud Run:

- `FIREBASE_SERVICE_ACCOUNT`: Base64 encoded Firebase service account JSON (prefixed with `base64:`)
- `SUBSCOLLECTION_ID`: Firestore subcollection identifier
- `ADMIN_PASSWORD`: Password for admin dashboard access
- `PORT`: Automatically set by Cloud Run (defaults to 8080)
- `CORS_ORIGIN`: Your Cloud Run service URL

### Health Check

The application includes a health check endpoint at `/api/registrations/count` which Cloud Run uses to verify the service is running.

### Security Best Practices

1. Use Google Secret Manager for sensitive values:
   ```bash
   # Store admin password in Secret Manager
   echo -n "your-password" | gcloud secrets create admin-password --data-file=-
   
   # Reference in Cloud Run
   gcloud run services update ${SERVICE_NAME} \
     --update-secrets=ADMIN_PASSWORD=admin-password:latest
   ```

2. Enable Cloud Armor for DDoS protection
3. Use IAM for access control if needed
4. Regularly rotate admin password

### Monitoring

- View logs: `gcloud run services logs read ${SERVICE_NAME} --region=${REGION}`
- Monitor metrics in Cloud Console
- Set up alerts for errors and latency

### Scaling

The service is configured with:
- Min instances: 0 (scales to zero when idle)
- Max instances: 10 (adjust based on expected load)
- Memory: 512Mi (increase if needed)
- CPU: 1 (increase for better performance)

Adjust these values based on your requirements.

