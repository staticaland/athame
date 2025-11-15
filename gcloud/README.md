# gcloud

A Dagger module for Google Cloud SDK (gcloud) commands.

## Setup

### Creating a service account key

1. Go to the [GCP Console Service Accounts page](https://console.cloud.google.com/iam-admin/serviceaccounts)
2. Select your project
3. Click "Create Service Account"
4. Give it a name and grant it the following roles:
   - **Artifact Registry Reader** - allows pulling container images from Artifact Registry
   - **Cloud Run Admin** - allows deploying and managing Cloud Run services
   - **Service Account User** - allows specifying which service account the deployed Cloud Run service will run as. When you deploy to Cloud Run, two service accounts are involved:
     - Your deployment service account (this one, with the key file) that performs the deployment
     - The runtime service account that the Cloud Run service will run as (defaults to the project's default compute service account if not specified)
     - You need this role so your deployment service account has permission to assign the runtime service account to the Cloud Run service
5. Click "Keys" → "Add Key" → "Create new key" → "JSON"
6. Save the downloaded JSON file (e.g., `$HOME/.config/gcloud/service-account-key.json`)
7. Minify the JSON to avoid secret redaction issues:
   ```bash
   jq -c . < input.json > output.json
   ```

## Usage

Deploy a container to Google Cloud Run:

```bash
dagger call --mod ./gcloud deploy \
  --service=my-service \
  --image=us-central1-docker.pkg.dev/my-project/my-repo/my-image:latest \
  --project=my-project \
  --region=us-central1 \
  --allow-unauthenticated=true \
  --service-account-key=file://$HOME/.config/gcloud/service-account-key.json
```

### Using 1Password

If your service account key is stored in 1Password as a document:

```bash
dagger call --mod ./gcloud deploy \
  --service=mkdocs-demo \
  --image=europe-north2-docker.pkg.dev/apps-477608/ghcr/staticaland/athame/mkdocs-demo:latest \
  --project=apps-477608 \
  --region=europe-west1 \
  --allow-unauthenticated=true \
  --service-account-key='cmd:op document get "Google Cloud - Service account key"'
```

### Private service

To deploy a service that requires IAM authentication for incoming requests:

```bash
dagger call --mod ./gcloud deploy \
  --service=mkdocs-demo \
  --image=europe-north2-docker.pkg.dev/apps-477608/ghcr/staticaland/athame/mkdocs-demo:latest \
  --project=apps-477608 \
  --region=europe-west1 \
  --service-account-key=file://$HOME/.config/gcloud/service-account-key.json
```

## Functions

- `base` - Returns the base container with Google Cloud SDK installed
- `deploy` - Deploys a container to Google Cloud Run
