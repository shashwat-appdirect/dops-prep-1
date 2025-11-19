# Cloud Run Deployment Debugging Guide

## Common Issues and Solutions

### Issue: Blank Page on Cloud Run

**Symptoms:**
- Page loads but shows blank/white screen
- Browser console shows errors
- Network tab shows failed requests

**Possible Causes:**

1. **API URL Configuration**
   - ✅ Fixed: Frontend now uses relative URLs (empty baseURL)
   - The frontend will automatically use the same domain as the page

2. **Static Files Not Served**
   - Check Cloud Run logs for "Static files directory found" message
   - Verify assets are being served at `/assets/` path
   - Check if index.html is being served at root `/`

3. **CORS Issues**
   - Ensure CORS_ORIGIN matches your Cloud Run URL exactly (without trailing slash)
   - Current: `https://dops-prep-1-1041941408881.asia-south1.run.app/` (has trailing slash)
   - Should be: `https://dops-prep-1-1041941408881.asia-south1.run.app` (no trailing slash)

4. **Firestore Connection**
   - Verify service account has Firestore permissions
   - Check Cloud Run logs for Firestore connection errors
   - Ensure GOOGLE_CLOUD_PROJECT is set correctly

## Debugging Steps

### 1. Check Cloud Run Logs

```bash
gcloud run services logs read dops-prep-1 --region asia-south1 --limit 50
```

Look for:
- "Static files directory found" - confirms static files are present
- "Server starting on port" - confirms server started
- Any Firestore connection errors
- Any panic/error messages

### 2. Test API Endpoints

```bash
# Test registration count
curl https://dops-prep-1-1041941408881.asia-south1.run.app/api/registrations/count

# Test speakers
curl https://dops-prep-1-1041941408881.asia-south1.run.app/api/speakers

# Test sessions
curl https://dops-prep-1-1041941408881.asia-south1.run.app/api/sessions
```

### 3. Test Static Files

```bash
# Test index.html
curl https://dops-prep-1-1041941408881.asia-south1.run.app/

# Test assets (replace with actual filename from your build)
curl https://dops-prep-1-1041941408881.asia-south1.run.app/assets/index-*.js
curl https://dops-prep-1-1041941408881.asia-south1.run.app/assets/index-*.css
```

### 4. Check Browser Console

Open browser DevTools (F12) and check:
- Console tab for JavaScript errors
- Network tab for failed requests (404, 500, CORS errors)
- Check if assets are loading (status 200)

### 5. Verify Environment Variables

In Cloud Run console, verify:
- `SUBSCOLLECTION_ID` = `shashwat.rawat`
- `ADMIN_PASSWORD` = `TechAdmin`
- `GOOGLE_CLOUD_PROJECT` = `india-tech-meetup-2025`
- `CORS_ORIGIN` = `https://dops-prep-1-1041941408881.asia-south1.run.app` (NO trailing slash)

## Quick Fixes

### Fix 1: Remove Trailing Slash from CORS_ORIGIN

Update CORS_ORIGIN in Cloud Run:
- Change from: `https://dops-prep-1-1041941408881.asia-south1.run.app/`
- Change to: `https://dops-prep-1-1041941408881.asia-south1.run.app`

### Fix 2: Rebuild and Redeploy

```bash
# Rebuild with latest changes
docker build -t appdirect-workshop:latest .

# Push to registry
gcloud builds submit --tag gcr.io/india-tech-meetup-2025/appdirect-workshop

# Redeploy
gcloud run deploy dops-prep-1 \
  --image gcr.io/india-tech-meetup-2025/appdirect-workshop \
  --region asia-south1 \
  --update-env-vars CORS_ORIGIN=https://dops-prep-1-1041941408881.asia-south1.run.app
```

### Fix 3: Check Firestore Permissions

Ensure Cloud Run service account has:
- Cloud Datastore User role
- Or Firestore read/write permissions

```bash
# Check service account
gcloud run services describe dops-prep-1 --region asia-south1 --format="value(spec.template.spec.serviceAccountName)"

# Grant Firestore permissions
gcloud projects add-iam-policy-binding india-tech-meetup-2025 \
  --member="serviceAccount:SERVICE_ACCOUNT_EMAIL" \
  --role="roles/datastore.user"
```

## Expected Behavior

After fixes:
1. Root URL (`/`) should serve index.html
2. Assets should load from `/assets/`
3. API calls should work at `/api/*`
4. No CORS errors in browser console
5. Frontend should render correctly

