# LocalStack Demo

A Dagger module demonstrating how to use LocalStack for local AWS development with both AWS CLI and Terraform.

## Functions

### create-bucket

Creates an S3 bucket in LocalStack using the AWS CLI.

```bash
dagger call create-bucket
```

**Output:**

```
make_bucket: demo-bucket
```

You can specify a custom bucket name:

```bash
dagger call create-bucket --bucket-name my-custom-bucket
```

### test-localstack

Verifies that LocalStack is running by making HTTP requests to the health endpoint.

```bash
dagger call test-localstack
```

### terraform-apply

Applies Terraform configuration against LocalStack using `terraform-local` (tflocal). This example creates an S3 bucket and outputs its details.

```bash
dagger call terraform-apply
```

**Output:**

```
Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # aws_s3_bucket.demo will be created
  + resource "aws_s3_bucket" "demo" {
      + acceleration_status         = (known after apply)
      + acl                         = (known after apply)
      + arn                         = (known after apply)
      + bucket                      = "demo-bucket"
      + bucket_domain_name          = (known after apply)
      + bucket_prefix               = (known after apply)
      + bucket_regional_domain_name = (known after apply)
      + force_destroy               = false
      + hosted_zone_id              = (known after apply)
      + id                          = (known after apply)
      + object_lock_enabled         = (known after apply)
      + policy                      = (known after apply)
      + region                      = (known after apply)
      + request_payer               = (known after apply)
      + tags_all                    = (known after apply)
      + website_domain              = (known after apply)
      + website_endpoint            = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + bucket_arn  = (known after apply)
  + bucket_name = "demo-bucket"

aws_s3_bucket.demo: Creating...
aws_s3_bucket.demo: Creation complete after 0s [id=demo-bucket]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

bucket_arn = "arn:aws:s3:::demo-bucket"
bucket_name = "demo-bucket"
```

You can specify a custom working directory containing Terraform files:

```bash
dagger call terraform-apply --workdir path/to/terraform
```
