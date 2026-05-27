# Copyright (c) Thomas Geens

terraform {
  required_version = ">= 1.8.0"
  required_providers {
    resourcenamingtool = {
      source  = "thomasgeens/resourcenamingtool"
      version = "~> 0.1.0"
    }
  }
}

# An example of using the `resourcenamingtool_status` data source to retrieve the status of the Resource Naming Tool provider.
data "resourcenamingtool_status" "example" {}

# Output the provider version
output "resourcenamingtool_provider_version" {
  value = data.resourcenamingtool_status.example.provider_version
}

# Output the Go version
output "resourcenamingtool_go_version" {
  value = data.resourcenamingtool_status.example.go_version
}
