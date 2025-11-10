# gcloud

A Dagger module for Google Cloud SDK (gcloud) commands.

## Setup

### Creating a service account key

1. Go to the [GCP Console Service Accounts page](https://console.cloud.google.com/iam-admin/serviceaccounts)
2. Select your project
3. Click "Create Service Account"
4. Give it a name and grant it the "Cloud Run Admin" role
5. Click "Keys" → "Add Key" → "Create new key" → "JSON"
6. Save the downloaded JSON file (e.g., `$HOME/.config/gcloud/service-account-key.json`)

## Usage

Deploy a container to Google Cloud Run:

```bash
dagger call --mod ./gcloud deploy \
  --service=my-service \
  --image=us-central1-docker.pkg.dev/my-project/my-repo/my-image:latest \
  --region=us-central1 \
  --allow-unauthenticated=true \
  --service-account-key=file://$HOME/.config/gcloud/service-account-key.json
```

### Using 1Password

If your service account key is stored in 1Password as a document:

```bash
dagger call --mod ./gcloud deploy \
  --service=my-service \
  --image=us-central1-docker.pkg.dev/my-project/my-repo/my-image:latest \
  --region=us-central1 \
  --allow-unauthenticated=true \
  --service-account-key='cmd:op document get "GCP Service Account Key"'
```

### Private service

To deploy a service that requires IAM authentication for incoming requests:

```bash
dagger call --mod ./gcloud deploy \
  --service=my-service \
  --image=us-central1-docker.pkg.dev/my-project/my-repo/my-image:latest \
  --region=us-central1 \
  --service-account-key=file://$HOME/.config/gcloud/service-account-key.json
```

## Functions

- `base` - Returns the base container with Google Cloud SDK installed
- `deploy` - Deploys a container to Google Cloud Run
