terraform {
  required_version = ">= 1.0"
}

# Simple hello world resource using null_resource
resource "null_resource" "hello" {
  triggers = {
    message = "Hello, World!"
  }

  provisioner "local-exec" {
    command = "echo ${self.triggers.message}"
  }
}

# Random pet name generator
resource "random_pet" "example" {
  length    = 2
  separator = "-"
}

# Local value example
locals {
  greeting = "Hello from Terraform!"
  timestamp = timestamp()
}

# Output examples
output "greeting" {
  value = local.greeting
}

output "pet_name" {
  value = random_pet.example.id
}

output "message" {
  value = "Terraform configuration validated successfully"
}
