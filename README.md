# terraform-provider-resourcenamingtool

[![CI](https://github.com/thomasgeens/terraform-provider-resourcenamingtool/actions/workflows/test.yaml/badge.svg)](https://github.com/thomasgeens/terraform-provider-resourcenamingtool/actions/workflows/test.yaml)
[![Registry](https://img.shields.io/badge/Terraform_Registry-thomasgeens%2Fresourcenamingtool-blueviolet)](https://registry.terraform.io/providers/thomasgeens/resourcenamingtool/latest)

Terraform provider that generates consistent, cloud-compliant resource names following Microsoft CAF, AWS WAF, and GCP naming conventions.

## Requirements

- Terraform >= 1.8.0 (required for [Provider Functions](https://developer.hashicorp.com/terraform/plugin/framework/functions))

## Installation

```hcl
terraform {
  required_version = ">= 1.8.0"
  required_providers {
    resourcenamingtool = {
      source  = "thomasgeens/resourcenamingtool"
      version = "~> 1.0"
    }
  }
}
```

## Quick Start

```hcl
provider "resourcenamingtool" {
  default_environment = {
    fullname  = "production"
    shortcode = "prod"
    char      = "p"
  }
  default_region = {
    fullname  = "westeurope"
    shortcode = "we"
    char      = "w"
  }
  default_basename = {
    fullname  = "myapp"
    shortcode = "app"
    char      = "a"
  }
  additional_naming_patterns = {
    "azurerm_resource_group"   = "rg-{basename}-{environment:short}-{region:short}"
    "azurerm_storage_account"  = "{basename}{environment:char}{region:char}{instance}"
  }
}

# Required: loads provider config during ValidateConfig
data "resourcenamingtool_status" "init" {}

output "resource_group_name" {
  value = provider::resourcenamingtool::generate_resource_name([{
    resource_type = { fullname = "azurerm_resource_group" }
  }])
  # → "rg-myapp-prod-we"
}

output "storage_account_name" {
  value = provider::resourcenamingtool::generate_resource_name([{
    resource_type = { fullname = "azurerm_storage_account" }
    additional_components = {
      "instance.fullname"  = "001"
      "instance.shortcode" = "01"
      "instance.char"      = "1"
    }
  }])
  # → "myapppw1"
}
```

> **Note:** The `data "resourcenamingtool_status" "init" {}` block is required to ensure the provider configuration is loaded during the `ValidateConfig` RPC. Add it once in your `provider.tf`.

## Documentation

Full schema, component reference, and cloud-provider examples:
[registry.terraform.io/providers/thomasgeens/resourcenamingtool](https://registry.terraform.io/providers/thomasgeens/resourcenamingtool/latest/docs)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, commit conventions, and the release process.
