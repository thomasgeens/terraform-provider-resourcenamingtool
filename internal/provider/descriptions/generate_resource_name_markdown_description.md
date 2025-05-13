# Resource Name Generator

Generates standardized, consistent resource names for cloud resources following industry best practices.

## Overview

This function applies cloud provider-specific naming patterns based on resource type and allows customizing all aspects of the name. It supports:

- Microsoft's Cloud Adoption Framework (CAF) for Azure
- AWS Well-Architected Framework (WAF) for AWS
- Google Cloud's recommended naming conventions

## Name Structure Diagram

```text
┌─────────────────────────────────────────────────────────────────┐
│                      RESOURCE NAME STRUCTURE                    │
├───────────┬───────────┬─────────┬───────────┬─────────┬─────────┤
│ Provider/ │ Resource  │         │           │         │         │
│ Org Prefix│   Type    │ Basename│Environment│ Region  │Instance │
│           │ (short)   │         │  (short)  │ (short) │         │
├───────────┼───────────┼─────────┼───────────┼─────────┼─────────┤
│    az     │    vm     │   app   │    prod   │   eus   │   01    │
└───────────┴───────────┴─────────┴───────────┴─────────┴─────────┘
      ↓           ↓          ↓          ↓          ↓         ↓
      └───────────┴──────────┴──────────┴──────────┴─────────┘
                                  ↓
                            azvmappprodeus01
```

## Component Representation

Each component provides information about a specific aspect of the resource naming and can be represented in three ways:

| Representation | Format              | Example (Environment) | Required | Description                                                                                                   |
|----------------|---------------------|----------------------|-----------|---------------------------------------------------------------------------------------------------------------|
| **fullname**   | "{component}"       | production           | Yes       | Complete/long form                                                                                            |
| **shortcode**  | "{component:short}" | prod                 | No        | Abbreviated form, if not provided will be constructed from the first 3 characters of the fullname             |
| **char**       | "{component:char}"  | p                    | No        | Single character representation, if not provided will be constructed from the first character of the fullname |

The naming pattern determines which representation is used.

## Supported Components

### Core Resource Components

| Component        | Description                                      | Example (fullname/shortcode/char) |
|------------------|--------------------------------------------------|----------------------------------|
| **resource_type**| Type of resource being created                   | virtual_machine / vm / v         |
| **resource_prefix** | Prefix before resource type abbreviation      | azurerm / az / a                 |
| **basename**     | Core identifying name for the resource           | webapp / web / w                 |
| **environment**  | Deployment environment                           | production / prod / p            |
| **region**       | Cloud region where resource is deployed          | eastus / eus / e                 |
| **instance**     | Instance identifier for multiple instances       | 001 / 01 / 1                     |

### Organization Components

| Component        | Description                                      | Example (fullname/shortcode/char) |
|------------------|--------------------------------------------------|----------------------------------|
| **organization** | Organization owning the resource                 | contoso / cto / c                |
| **business_unit**| Business unit within the organization            | finance / fin / f                |
| **cost_center**  | Financial cost center                            | marketing / mkt / m              |
| **project**      | Project associated with the resource             | website / web / w                |
| **application**  | Application using the resource                   | inventory / inv / i              |
| **workload**     | Function or purpose of the workload              | api / api / a                    |

### Provider-Specific Components

| Component        | Description                                      | Example (fullname/shortcode/char) |
|------------------|--------------------------------------------------|----------------------------------|
| **subscription** | Subscription context (primarily for Azure)       | subscription01 / sub01 / s       |
| **location**     | Alternative to region for some providers         | westeurope / weu / w             |
| **domain**       | Domain name or business domain                   | contoso.com / cto / c            |
| **criticality**  | Importance of the resource                       | critical / crit / c              |

### Initiative/Solution Components

| Component        | Description                                      | Example (fullname/shortcode/char) |
|------------------|--------------------------------------------------|----------------------------------|
| **initiative**   | Business initiative the resource belongs to      | digital_transformation / digitx / d |
| **solution**     | Solution name, architecture or pattern           | microservices / micro / m        |

## Extension Points

### additional_naming_patterns

A map of custom naming patterns for specific resource types.

```hcl
additional_naming_patterns = {
  "custom_resource" = "{resource_type:short}-{basename}-{environment:short}"
  "azurerm_logic_app_workflow" = "logic-{basename}-{environment:short}-custom"
}
```

### additional_components

A map of custom component values that can be used in naming patterns.

```hcl
additional_components = {
  "department.fullname" = "engineering"
  "department.shortcode" = "eng"
  "department.char" = "e"
"team.fullname" = "platform"
  "team.shortcode" = "plf"
  "team.char" = "p"
}
```

## Complete Example
Generating a resource name with component values, custom patterns, and custom components:

```hcl
resource "arm_storage_account" "backup_prod" {
  name                     = provider::resourcenamingtool::generate_resource_name(
    [
      {
        resource_type = {
          "fullname"  = "storage_account"
          "shortcode" = "st"
          "char"      = "s"
        },
        additional_components = {
          "slug.fullname"  = "backup"
          "slug.shortcode" = "bkp"
          "slug.char"      = "b"
          "slot.fullname"  = "prod"
          "slot.shortcode" = "prd"
          "slot.char"      = "p"
        },
        additional_naming_patterns = {
          "storage_account" = "{resource_type:short}{slug:short}{region:short}{slot:char}{instance}"
        }
      }
    ]
  )
  resource_group_name      = provider::resourcenamingtool::generate_resource_name(
    [
      {
        resource_type = {
          "fullname"  = "resource_group"
        }
      }
    ]
  )
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
```

# Cloud Provider Resource Type Examples

## Azure Resource Management (ARM)

| Resource Type | Naming Pattern | Example Output | Max Length |
|---------------|----------------|----------------|------------|
| azurerm_resource_group | rg-{basename}-{environment:short}-{region:short} | rg-app-prod-eus | 90 |
| azurerm_virtual_network | vnet-{basename}-{environment:short}-{region:short} | vnet-app-prod-eus | 64 |
| azurerm_subnet | snet-{basename}-{environment:short}-{instance} | snet-app-prod-01 | 80 |
| azurerm_network_security_group | nsg-{basename}-{environment:short}-{region:short} | nsg-app-prod-eus | 80 |
| azurerm_route_table | rt-{basename}-{environment:short}-{region:short} | rt-app-prod-eus | 80 |
| azurerm_virtual_machine | vm-{basename}-{environment:short}-{region:short}-{instance} | vm-app-prod-eus-01 | 64 |
| azurerm_availability_set | avs-{basename}-{environment:short}-{region:short} | avs-app-prod-eus | 80 |
| azurerm_vm_scale_set | vmss-{basename}-{environment:short}-{region:short} | vmss-app-prod-eus | 64 |
| azurerm_kubernetes_cluster | aks-{basename}-{environment:short}-{region:short} | aks-app-prod-eus | 63 |
| azurerm_storage_account | {basename}{environment:char}{region:char}{instance} | apppeus01 | 24 |
| azurerm_storage_container | sc-{basename}-{environment:short} | sc-app-prod | 63 |
| azurerm_sql_server | sql-{basename}-{environment:short}-{region:short} | sql-app-prod-eus | 63 |
| azurerm_sql_database | sqldb-{basename}-{environment:short} | sqldb-app-prod | 128 |
| azurerm_cosmosdb_account | cosmos-{basename}-{environment:short}-{region:short} | cosmos-app-prod-eus | 44 |
| azurerm_mysql_server | mysql-{basename}-{environment:short}-{region:short} | mysql-app-prod-eus | 63 |
| azurerm_postgresql_server | psql-{basename}-{environment:short}-{region:short} | psql-app-prod-eus | 63 |
| azurerm_app_service | app-{basename}-{environment:short}-{region:short} | app-api-prod-eus | 60 |
| azurerm_app_service_plan | plan-{basename}-{environment:short}-{region:short} | plan-app-prod-eus | 40 |
| azurerm_function_app | func-{basename}-{environment:short}-{region:short} | func-app-prod-eus | 60 |
| azurerm_key_vault | kv-{basename}-{environment:short}-{region:short} | kv-app-prod-eus | 24 |
| azurerm_servicebus_namespace | sb-{basename}-{environment:short}-{region:short} | sb-app-prod-eus | 50 |
| azurerm_eventhub_namespace | evh-{basename}-{environment:short}-{region:short} | evh-app-prod-eus | 50 |
| azurerm_eventgrid_topic | evg-{basename}-{environment:short}-{region:short} | evg-app-prod-eus | 50 |
| azurerm_logic_app_workflow | logic-{basename}-{environment:short} | logic-app-prod | 80 |
| azurerm_container_registry | acr{basename}{environment:char}{region:char} | acrappeust | 50 |
| azurerm_container_group | aci-{basename}-{environment:short} | aci-app-prod | 63 |
| azurerm_log_analytics_workspace | log-{basename}-{environment:short}-{region:short} | log-app-prod-eus | 63 |
| azurerm_application_insights | appi-{basename}-{environment:short}-{region:short} | appi-app-prod-eus | 260 |
| azurerm_public_ip | pip-{basename}-{environment:short}-{region:short} | pip-app-prod-eus | 80 |
| azurerm_lb | lb-{basename}-{environment:short}-{region:short} | lb-app-prod-eus | 80 |
| azurerm_application_gateway | agw-{basename}-{environment:short}-{region:short} | agw-app-prod-eus | 80 |
| azurerm_network_interface | nic-{basename}-{environment:short} | nic-app-prod | 80 |
| azurerm_private_endpoint | pe-{basename}-{environment:short} | pe-app-prod | 80 |
| azurerm_user_assigned_identity | id-{basename}-{environment:short}-{region:short} | id-app-prod-eus | 128 |
| azurerm_monitor_action_group | ag-{basename}-{environment:short} | ag-app-prod | 260 |
| azurerm_monitor_metric_alert | ar-{basename}-{environment:short} | ar-app-prod | 260 |

### Azure-Specific Naming Recommendations

It is recommended to follow Microsoft's Cloud Adoption Framework best practices for Azure resource naming:

- Resource names are consistent and follow a predictable pattern
- Names include necessary context for identification (environment, location, etc.)
- Names comply with Azure resource-specific length and character limitations
- Storage accounts and certain resources use all lowercase with no separators
- Resource types that don't support hyphens use camelCase or no separators

Example Azure Provider Configuration:

```hcl
provider "resourcenamingtool" {
  default_resource_type = {
    fullname  = "azurerm_resource_group"
    shortcode = "rg"
    char      = "r"
  }

  default_environment = {
    fullname  = "production"
    shortcode = "prod"
    char      = "p"
  }

  default_region = {
    fullname  = "eastus"
    shortcode = "eus"
    char      = "e"
  }

  additional_naming_patterns = {
	// Azure Core Resources
    "azurerm_resource_group"         = "rg-{basename}-{environment:short}-{region:short}"
    "azurerm_virtual_network"        = "vnet-{basename}-{environment:short}-{region:short}"
    "azurerm_subnet"                 = "snet-{basename}-{environment:short}-{instance}"
    "azurerm_network_security_group" = "nsg-{basename}-{environment:short}-{region:short}"
    "azurerm_route_table"            = "rt-{basename}-{environment:short}-{region:short}"

    // Azure Compute Resources
    "azurerm_virtual_machine"        = "vm-{basename}-{environment:short}-{region:short}-{instance}"
	"azurerm_availability_set"       = "avs-{basename}-{environment:short}-{region:short}"
    "azurerm_vm_scale_set"           = "vmss-{basename}-{environment:short}-{region:short}"
    "azurerm_kubernetes_cluster"     = "aks-{basename}-{environment:short}-{region:short}"

    // Azure Storage Resources
    "azurerm_storage_account"        = "{basename}{environment:char}{region:char}{instance}"
	"azurerm_storage_container"      = "sc-{basename}-{environment:short}"

    // Azure Database Resources
    "azurerm_sql_server"             = "sql-{basename}-{environment:short}-{region:short}"
	"azurerm_sql_database"           = "sqldb-{basename}-{environment:short}"
    "azurerm_cosmosdb_account"       = "cosmos-{basename}-{environment:short}-{region:short}"
	"azurerm_mysql_server"           = "mysql-{basename}-{environment:short}-{region:short}"
    "azurerm_postgresql_server"      = "psql-{basename}-{environment:short}-{region:short}"

    // Azure App Resources
    "azurerm_app_service"            = "app-{basename}-{environment:short}-{region:short}"
	"azurerm_app_service_plan"       = "plan-{basename}-{environment:short}-{region:short}"
    "azurerm_function_app"           = "func-{basename}-{environment:short}-{region:short}"

    // Azure Security Resources
    "azurerm_key_vault"              = "kv-{basename}-{environment:short}-{region:short}"
    "azurerm_container_registry"     = "acr{basename}{environment:char}{region:char}"
    "azurerm_private_endpoint"       = "pe-{basename}-{environment:short}"

	// Azure Integration Resources
    "azurerm_servicebus_namespace"   = "sb-{basename}-{environment:short}-{region:short}"
    "azurerm_eventhub_namespace"     = "evh-{basename}-{environment:short}-{region:short}"
    "azurerm_eventgrid_topic"        = "evg-{basename}-{environment:short}-{region:short}"
    "azurerm_logic_app_workflow"     = "logic-{basename}-{environment:short}"

    // Azure Container Resources
    "azurerm_container_registry"     = "acr{basename}{environment:char}{region:char}"
    "azurerm_container_group"        = "aci-{basename}-{environment:short}"

    // Azure Analytics Resources
    "azurerm_log_analytics_workspace" = "log-{basename}-{environment:short}-{region:short}"
    "azurerm_application_insights"    = "appi-{basename}-{environment:short}-{region:short}"

    // Azure Network Resources
    "azurerm_public_ip"              = "pip-{basename}-{environment:short}-{region:short}"
    "azurerm_lb"                     = "lb-{basename}-{environment:short}-{region:short}"
    "azurerm_application_gateway"    = "agw-{basename}-{environment:short}-{region:short}"
    "azurerm_network_interface"      = "nic-{basename}-{environment:short}"
    "azurerm_private_endpoint"       = "pe-{basename}-{environment:short}"

    // Azure Identity Resources
    "azurerm_user_assigned_identity" = "id-{basename}-{environment:short}-{region:short}"

    // Azure Monitor Resources
    "azurerm_monitor_action_group"   = "ag-{basename}-{environment:short}"
    "azurerm_monitor_metric_alert"   = "ar-{basename}-{environment:short}"
  }
}
```

## Amazon Web Services (AWS)

| Resource Type | Naming Pattern | Example Output | Max Length |
|---------------|----------------|----------------|------------|
| aws_ec2_instance | ec2-{basename}-{environment:short}-{region:short}-{instance} | ec2-app-prod-usea-01 | 255 |
| aws_auto_scaling_group | asg-{basename}-{environment:short}-{region:short} | asg-app-prod-usea | 255 |
| aws_launch_template | lt-{basename}-{environment:short}-{region:short} | lt-app-prod-usea | 128 |
| aws_s3_bucket | {basename}-{environment:short}-{region:short}-{instance} | app-prod-usea-01 | 63 |
| aws_efs_file_system | efs-{basename}-{environment:short}-{region:short} | efs-app-prod-usea | 255 |
| aws_rds_instance | rds-{basename}-{environment:short}-{region:short} | rds-app-prod-usea | 63 |
| aws_rds_cluster | rdsc-{basename}-{environment:short}-{region:short} | rdsc-app-prod-usea | 63 |
| aws_dynamodb_table | ddb-{basename}-{environment:short}-{region:short} | ddb-app-prod-usea | 255 |
| aws_elasticache | ec-{basename}-{environment:short}-{region:short} | ec-app-prod-usea | 40 |
| aws_vpc | vpc-{basename}-{environment:short}-{region:short} | vpc-app-prod-usea | 255 |
| aws_subnet | snet-{basename}-{environment:short}-{region:short}-{instance} | snet-app-prod-usea-01 | 255 |
| aws_security_group | sg-{basename}-{environment:short}-{region:short} | sg-app-prod-usea | 255 |
| aws_route_table | rt-{basename}-{environment:short}-{region:short} | rt-app-prod-usea | 255 |
| aws_elastic_ip | eip-{basename}-{environment:short} | eip-app-prod | 255 |
| aws_nat_gateway | nat-{basename}-{environment:short}-{region:short} | nat-app-prod-usea | 255 |
| aws_load_balancer | lb-{basename}-{environment:short}-{region:short} | lb-app-prod-usea | 32 |
| aws_target_group | tg-{basename}-{environment:short}-{region:short} | tg-app-prod-usea | 32 |
| aws_lambda_function | lambda-{basename}-{environment:short}-{region:short} | lambda-app-prod-usea | 64 |
| aws_layer | layer-{basename}-{environment:short}-{region:short} | layer-app-prod-usea | 64 |
| aws_ecr_repository | ecr-{basename}-{environment:short}-{region:short} | ecr-app-prod-usea | 256 |
| aws_ecs_cluster | ecs-{basename}-{environment:short}-{region:short} | ecs-app-prod-usea | 255 |
| aws_eks_cluster | eks-{basename}-{environment:short}-{region:short} | eks-app-prod-usea | 100 |
| aws_iam_role | role-{basename}-{environment:short} | role-app-prod | 64 |
| aws_iam_policy | pol-{basename}-{environment:short} | pol-app-prod | 128 |
| aws_iam_user | usr-{basename}-{environment:short} | usr-app-prod | 64 |
| aws_iam_group | grp-{basename}-{environment:short} | grp-app-prod | 128 |
| aws_cloudwatch_alarm | cwa-{basename}-{environment:short} | cwa-app-prod | 255 |
| aws_log_group | log-{basename}-{environment:short} | log-app-prod | 512 |
| aws_sns_topic | sns-{basename}-{environment:short}-{region:short} | sns-app-prod-usea | 256 |
| aws_sqs_queue | sqs-{basename}-{environment:short}-{region:short} | sqs-app-prod-usea | 80 |
| aws_api_gateway | api-{basename}-{environment:short}-{region:short} | api-app-prod-usea | 128 |
| aws_step_function | sf-{basename}-{environment:short}-{region:short} | sf-app-prod-usea | 80 |
| aws_cloudfront | cf-{basename}-{environment:short} | cf-app-prod | 128 |
| aws_hosted_zone | hz-{basename}-{environment:short} | hz-app-prod | 255 |
| aws_record_set | rs-{basename}-{environment:short} | rs-app-prod | 255 |

### AWS-Specific Naming Recommendations

AWS resources follow Well-Architected Framework best practices:

- Consistent naming patterns for each resource type
- Consider service-specific constraints (like S3 bucket naming restrictions)
- Use lowercase letters, numbers, and hyphens for most resources
- Consider global uniqueness requirements for certain resources
- Include environment and region context in resource names
- Add instance identifiers for resources that may have multiple instances

Example AWS Provider Configuration:

```hcl
provider "resourcenamingtool" {
  default_resource_type = {
    fullname  = "aws_ec2_instance"
    shortcode = "ec2"
    char      = "e"
  }

  default_environment = {
    fullname  = "production"
    shortcode = "prod"
    char      = "p"
  }

  default_region = {
    fullname  = "us-east-1"
    shortcode = "usea"
    char      = "e"
  }

  additional_naming_patterns = {
    // AWS Compute Resources
    "aws_ec2_instance"        = "ec2-{basename}-{environment:short}-{region:short}-{instance}"
    "aws_auto_scaling_group"  = "asg-{basename}-{environment:short}-{region:short}"
    "aws_launch_template"     = "lt-{basename}-{environment:short}-{region:short}"

    // AWS Storage Resources
    "aws_s3_bucket"           = "{basename}-{environment:short}-{region:short}-{instance}"
    "aws_efs_file_system"     = "efs-{basename}-{environment:short}-{region:short}"

    // AWS Database Resources
    "aws_rds_instance"        = "rds-{basename}-{environment:short}-{region:short}"
    "aws_rds_cluster"         = "rdsc-{basename}-{environment:short}-{region:short}"
    "aws_dynamodb_table"      = "ddb-{basename}-{environment:short}-{region:short}"
    "aws_elasticache"         = "ec-{basename}-{environment:short}-{region:short}"

    // AWS Network Resources
    "aws_vpc"                 = "vpc-{basename}-{environment:short}-{region:short}"
    "aws_subnet"              = "snet-{basename}-{environment:short}-{region:short}-{instance}"
    "aws_security_group"      = "sg-{basename}-{environment:short}-{region:short}"
    "aws_route_table"         = "rt-{basename}-{environment:short}-{region:short}"
    "aws_elastic_ip"          = "eip-{basename}-{environment:short}"
    "aws_nat_gateway"         = "nat-{basename}-{environment:short}-{region:short}"
    "aws_load_balancer"       = "lb-{basename}-{environment:short}-{region:short}"
    "aws_target_group"        = "tg-{basename}-{environment:short}-{region:short}"

    // AWS Lambda Resources
    "aws_lambda_function"     = "lambda-{basename}-{environment:short}-{region:short}"
    "aws_layer"               = "layer-{basename}-{environment:short}-{region:short}"

    // AWS Container Resources
    "aws_ecr_repository"      = "ecr-{basename}-{environment:short}-{region:short}"
    "aws_ecs_cluster"         = "ecs-{basename}-{environment:short}-{region:short}"
    "aws_eks_cluster"         = "eks-{basename}-{environment:short}-{region:short}"

    // AWS IAM Resources
    "aws_iam_role"            = "role-{basename}-{environment:short}"
    "aws_iam_policy"          = "pol-{basename}-{environment:short}"
    "aws_iam_user"            = "usr-{basename}-{environment:short}"
    "aws_iam_group"           = "grp-{basename}-{environment:short}"

    // AWS Monitoring Resources
    "aws_cloudwatch_alarm"    = "cwa-{basename}-{environment:short}"
    "aws_log_group"           = "log-{basename}-{environment:short}"
    "aws_sns_topic"           = "sns-{basename}-{environment:short}-{region:short}"
    "aws_sqs_queue"           = "sqs-{basename}-{environment:short}-{region:short}"

    // AWS Application Resources
    "aws_api_gateway"         = "api-{basename}-{environment:short}-{region:short}"
    "aws_step_function"       = "sf-{basename}-{environment:short}-{region:short}"
    "aws_cloudfront"          = "cf-{basename}-{environment:short}"

    // AWS Route53 Resources
    "aws_hosted_zone"         = "hz-{basename}-{environment:short}"
    "aws_record_set"          = "rs-{basename}-{environment:short}"
  }
}
```

## Google Cloud Platform (GCP)

| Resource Type | Naming Pattern | Example Output | Max Length |
|---------------|----------------|----------------|------------|
| google_compute_instance | vm-{basename}-{environment:short}-{region:short}-{instance} | vm-app-prod-usea1-01 | 63 |
| google_compute_instance_group | ig-{basename}-{environment:short}-{region:short} | ig-app-prod-usea1 | 63 |
| google_compute_instance_template | it-{basename}-{environment:short}-{region:short} | it-app-prod-usea1 | 63 |
| google_compute_disk | disk-{basename}-{environment:short}-{region:short} | disk-app-prod-usea1 | 63 |
| google_compute_snapshot | snap-{basename}-{environment:short}-{region:short} | snap-app-prod-usea1 | 63 |
| google_compute_image | img-{basename}-{environment:short} | img-app-prod | 63 |
| google_container_cluster | gke-{basename}-{environment:short}-{region:short} | gke-app-prod-usea1 | 40 |
| google_container_node_pool | np-{basename}-{environment:short}-{region:short} | np-app-prod-usea1 | 63 |
| google_storage_bucket | {basename}-{environment:short}-{region:short} | app-prod-usea1 | 63 |
| google_filestore_instance | fs-{basename}-{environment:short}-{region:short} | fs-app-prod-usea1 | 63 |
| google_compute_network | vpc-{basename}-{environment:short} | vpc-app-prod | 63 |
| google_compute_subnetwork | subnet-{basename}-{environment:short}-{region:short} | subnet-app-prod-usea1 | 63 |
| google_compute_firewall | fw-{basename}-{environment:short} | fw-app-prod | 63 |
| google_compute_router | router-{basename}-{environment:short}-{region:short} | router-app-prod-usea1 | 63 |
| google_compute_address | addr-{basename}-{environment:short}-{region:short} | addr-app-prod-usea1 | 63 |
| google_compute_global_address | gaddr-{basename}-{environment:short} | gaddr-app-prod | 63 |
| google_compute_forwarding_rule | fr-{basename}-{environment:short}-{region:short} | fr-app-prod-usea1 | 63 |
| google_compute_target_http_proxy | http-proxy-{basename}-{environment:short} | http-proxy-app-prod | 63 |
| google_compute_target_https_proxy | https-proxy-{basename}-{environment:short} | https-proxy-app-prod | 63 |
| google_compute_ssl_certificate | cert-{basename}-{environment:short} | cert-app-prod | 63 |
| google_compute_url_map | url-map-{basename}-{environment:short} | url-map-app-prod | 63 |
| google_compute_backend_service | bes-{basename}-{environment:short} | bes-app-prod | 63 |
| google_sql_database_instance | sql-{basename}-{environment:short}-{region:short} | sql-app-prod-usea1 | 98 |
| google_sql_database | db-{basename}-{environment:short} | db-app-prod | 128 |
| google_bigtable_instance | bt-{basename}-{environment:short}-{region:short} | bt-app-prod-usea1 | 63 |
| google_bigtable_table | bt-tbl-{basename}-{environment:short} | bt-tbl-app-prod | 50 |
| google_spanner_instance | spanner-{basename}-{environment:short} | spanner-app-prod | 30 |
| google_spanner_database | spanner-db-{basename}-{environment:short} | spanner-db-app-prod | 30 |
| google_firestore_database | fs-db-{basename}-{environment:short} | fs-db-app-prod | 63 |
| google_cloudfunctions_function | func-{basename}-{environment:short}-{region:short} | func-app-prod-usea1 | 63 |
| google_cloud_run_service | run-{basename}-{environment:short}-{region:short} | run-app-prod-usea1 | 63 |
| google_app_engine_application | app-{basename}-{environment:short} | app-app-prod | 63 |
| google_app_engine_standard_app_version | app-{basename}-{environment:short}-{version} | app-app-prod-v1 | 100 |
| google_bigquery_dataset | bq-ds-{basename}-{environment:short} | bq-ds-app-prod | 1024 |
| google_bigquery_table | bq-tbl-{basename}-{environment:short} | bq-tbl-app-prod | 1024 |
| google_dataflow_job | df-{basename}-{environment:short} | df-app-prod | 1024 |
| google_dataproc_cluster | dp-{basename}-{environment:short}-{region:short} | dp-app-prod-usea1 | 52 |
| google_pubsub_topic | ps-topic-{basename}-{environment:short} | ps-topic-app-prod | 255 |
| google_pubsub_subscription | ps-sub-{basename}-{environment:short} | ps-sub-app-prod | 255 |
| google_service_account | sa-{basename}-{environment:short} | sa-app-prod | 30 |
| google_project_iam_custom_role | role-{basename}-{environment:short} | role-app-prod | 64 |
| google_kms_key_ring | kr-{basename}-{environment:short}-{region:short} | kr-app-prod-usea1 | 63 |
| google_kms_crypto_key | kms-{basename}-{environment:short} | kms-app-prod | 63 |
| google_secret_manager_secret | secret-{basename}-{environment:short} | secret-app-prod | 255 |
| google_monitoring_alert_policy | alert-{basename}-{environment:short} | alert-app-prod | 100 |
| google_logging_metric | log-{basename}-{environment:short} | log-app-prod | 100 |
| google_monitoring_notification_channel | notif-{basename}-{environment:short} | notif-app-prod | 100 |
| google_monitoring_dashboard | dash-{basename}-{environment:short} | dash-app-prod | 100 |

### GCP-Specific Naming Recommendations

Google Cloud Platform resources follow Google's recommended naming conventions:

- Resource names should be consistent and follow a predictable pattern
- Use lowercase letters, numbers, and hyphens for most resources
- Certain resources have specific character limitations (e.g., service accounts, spanner instances)
- Include necessary context (project, environment, region) in resource names
- Consider global uniqueness requirements for bucket names and other global resources
- Follow Google's resource-specific guidelines for character constraints

Example GCP Provider Configuration:

```hcl
provider "resourcenamingtool" {
  default_resource_type = {
    fullname  = "google_compute_instance"
    shortcode = "vm"
    char      = "v"
  }

  default_environment = {
    fullname  = "production"
    shortcode = "prod"
    char      = "p"
  }

  default_region = {
    fullname  = "us-east1"
    shortcode = "usea1"
    char      = "e"
  }

  additional_naming_patterns = {
    // GCP Compute Resources
    "google_compute_instance"          = "vm-{basename}-{environment:short}-{region:short}-{instance}"
    "google_compute_instance_group"    = "ig-{basename}-{environment:short}-{region:short}"
    "google_compute_instance_template" = "it-{basename}-{environment:short}-{region:short}"
    "google_compute_disk"              = "disk-{basename}-{environment:short}-{region:short}"
    "google_compute_snapshot"          = "snap-{basename}-{environment:short}-{region:short}"
    "google_compute_image"             = "img-{basename}-{environment:short}"

    // GCP Kubernetes Resources
    "google_container_cluster"         = "gke-{basename}-{environment:short}-{region:short}"
    "google_container_node_pool"       = "np-{basename}-{environment:short}-{region:short}"

    // GCP Storage Resources
    "google_storage_bucket"            = "{basename}-{environment:short}-{region:short}"
    "google_filestore_instance"        = "fs-{basename}-{environment:short}-{region:short}"

    // GCP Network Resources
    "google_compute_network"           = "vpc-{basename}-{environment:short}"
    "google_compute_subnetwork"        = "subnet-{basename}-{environment:short}-{region:short}"
    "google_compute_firewall"          = "fw-{basename}-{environment:short}"
    "google_compute_router"            = "router-{basename}-{environment:short}-{region:short}"
    "google_compute_address"           = "addr-{basename}-{environment:short}-{region:short}"
    "google_compute_global_address"    = "gaddr-{basename}-{environment:short}"
    "google_compute_forwarding_rule"   = "fr-{basename}-{environment:short}-{region:short}"
    "google_compute_target_http_proxy" = "http-proxy-{basename}-{environment:short}"
    "google_compute_target_https_proxy" = "https-proxy-{basename}-{environment:short}"
    "google_compute_ssl_certificate"   = "cert-{basename}-{environment:short}"
    "google_compute_url_map"           = "url-map-{basename}-{environment:short}"
    "google_compute_backend_service"   = "bes-{basename}-{environment:short}"

    // GCP Database Resources
    "google_sql_database_instance"     = "sql-{basename}-{environment:short}-{region:short}"
    "google_sql_database"              = "db-{basename}-{environment:short}"
    "google_bigtable_instance"         = "bt-{basename}-{environment:short}-{region:short}"
    "google_bigtable_table"            = "bt-tbl-{basename}-{environment:short}"
    "google_spanner_instance"          = "spanner-{basename}-{environment:short}"
    "google_spanner_database"          = "spanner-db-{basename}-{environment:short}"
    "google_firestore_database"        = "fs-db-{basename}-{environment:short}"

    // GCP Serverless Resources
    "google_cloudfunctions_function"   = "func-{basename}-{environment:short}-{region:short}"
    "google_cloud_run_service"         = "run-{basename}-{environment:short}-{region:short}"
    "google_app_engine_application"    = "app-{basename}-{environment:short}"
    "google_app_engine_standard_app_version" = "app-{basename}-{environment:short}-{version}"

    // GCP Data Analytics Resources
    "google_bigquery_dataset"          = "bq-ds-{basename}-{environment:short}"
    "google_bigquery_table"            = "bq-tbl-{basename}-{environment:short}"
    "google_dataflow_job"              = "df-{basename}-{environment:short}"
    "google_dataproc_cluster"          = "dp-{basename}-{environment:short}-{region:short}"
    "google_pubsub_topic"              = "ps-topic-{basename}-{environment:short}"
    "google_pubsub_subscription"       = "ps-sub-{basename}-{environment:short}"

    // GCP IAM Resources
    "google_service_account"           = "sa-{basename}-{environment:short}"
    "google_project_iam_custom_role"   = "role-{basename}-{environment:short}"

    // GCP Security Resources
    "google_kms_key_ring"              = "kr-{basename}-{environment:short}-{region:short}"
    "google_kms_crypto_key"            = "kms-{basename}-{environment:short}"
    "google_secret_manager_secret"     = "secret-{basename}-{environment:short}"

    // GCP Monitoring Resources
    "google_monitoring_alert_policy"   = "alert-{basename}-{environment:short}"
    "google_logging_metric"            = "log-{basename}-{environment:short}"
    "google_monitoring_notification_channel" = "notif-{basename}-{environment:short}"
    "google_monitoring_dashboard"      = "dash-{basename}-{environment:short}"
  }
}
```
