# Copyright (c) Thomas Geens

# Define requirements for the provider and Terraform version
terraform {
  required_version = ">= 1.8.0" // Required for Terraform Provider Defined Functions
  required_providers {
    resourcenamingtool = {
      source  = "thomasgeens/resourcenamingtool"
      version = "~> 0.1.0" // Specify the version of the provider
    }
  }
}

# Define the provider configuration
# This configuration is used to set default values for the resource naming tool
provider "resourcenamingtool" {
  default_resource_type = {
    fullname  = "azurerm_resource_group"
    shortcode = "rg"
    char      = "r"
  }

  default_environment = {
    fullname  = "production"
    shortcode = "prd"
    char      = "p"
  }

  default_instance = {
    fullname  = "00001"
    shortcode = "001"
    char      = "1"
  }

  default_basename = {
    fullname  = "example"
    shortcode = "ex"
    char      = "e"
  }

  default_subscription = {
    fullname  = "prod-01"
    shortcode = "p01"
    char      = "p"
  }

  default_region = {
    fullname  = "westeurope"
    shortcode = "we"
    char      = "w"
  }

  additional_naming_patterns = {
    // Azure Core Resources
    "azurerm_resource_group" : "rg-{basename}-{environment:short}-{region:short}",
    "azurerm_virtual_network" : "vnet-{basename}-{environment:short}-{region:short}",
    "azurerm_subnet" : "snet-{basename}-{environment:short}-{instance}",
    "azurerm_network_security_group" : "nsg-{basename}-{environment:short}-{region:short}",
    "azurerm_route_table" : "rt-{basename}-{environment:short}-{region:short}",

    // Azure Compute Resources
    "azurerm_virtual_machine" : "vm-{basename}-{environment:short}-{region:short}-{instance}",
    "azurerm_availability_set" : "avs-{basename}-{environment:short}-{region:short}",
    "azurerm_vm_scale_set" : "vmss-{basename}-{environment:short}-{region:short}",
    "azurerm_kubernetes_cluster" : "aks-{basename}-{environment:short}-{region:short}",

    // Azure Storage Resources
    "azurerm_storage_account" : "{basename}{environment:char}{region:char}{instance}",
    "azurerm_storage_container" : "sc-{basename}-{environment:short}",

    // Azure Database Resources
    "azurerm_sql_server" : "sql-{basename}-{environment:short}-{region:short}",
    "azurerm_sql_database" : "sqldb-{basename}-{environment:short}",
    "azurerm_cosmosdb_account" : "cosmos-{basename}-{environment:short}-{region:short}",
    "azurerm_mysql_server" : "mysql-{basename}-{environment:short}-{region:short}",
    "azurerm_postgresql_server" : "psql-{basename}-{environment:short}-{region:short}",

    // Azure App Resources
    "azurerm_app_service" : "app-{basename}-{environment:short}-{region:short}",
    "azurerm_app_service_plan" : "plan-{basename}-{environment:short}-{region:short}",
    "azurerm_function_app" : "func-{basename}-{environment:short}-{region:short}",

    // Azure Security Resources
    "azurerm_key_vault" : "kv-{basename}-{environment:short}-{region:short}",

    // Azure Integration Resources
    "azurerm_servicebus_namespace" : "sb-{basename}-{environment:short}-{region:short}",
    "azurerm_eventhub_namespace" : "evh-{basename}-{environment:short}-{region:short}",
    "azurerm_eventgrid_topic" : "evg-{basename}-{environment:short}-{region:short}",
    "azurerm_logic_app_workflow" : "logic-{basename}-{environment:short}",

    // Azure Container Resources
    "azurerm_container_registry" : "acr{basename}{environment:char}{region:char}",
    "azurerm_container_group" : "aci-{basename}-{environment:short}",

    // Azure Analytics Resources
    "azurerm_log_analytics_workspace" : "log-{basename}-{environment:short}-{region:short}",
    "azurerm_application_insights" : "appi-{basename}-{environment:short}-{region:short}",

    // Azure Network Resources
    "azurerm_public_ip" : "pip-{basename}-{environment:short}-{region:short}",
    "azurerm_lb" : "lb-{basename}-{environment:short}-{region:short}",
    "azurerm_application_gateway" : "agw-{basename}-{environment:short}-{region:short}",
    "azurerm_network_interface" : "nic-{basename}-{environment:short}",
    "azurerm_private_endpoint" : "pe-{basename}-{environment:short}",

    // Azure Identity Resources
    "azurerm_user_assigned_identity" : "id-{basename}-{environment:short}-{region:short}",

    // Azure Monitor Resources
    "azurerm_monitor_action_group" : "ag-{basename}-{environment:short}",
    "azurerm_monitor_metric_alert" : "ar-{basename}-{environment:short}",

    // AWS Compute Resources
    "aws_ec2_instance"       = "ec2-{basename}-{environment:short}-{region:short}-{instance}",
    "aws_auto_scaling_group" = "asg-{basename}-{environment:short}-{region:short}",
    "aws_launch_template"    = "lt-{basename}-{environment:short}-{region:short}",

    // AWS Storage Resources
    "aws_s3_bucket"       = "{basename}-{environment:short}-{region:short}-{instance}",
    "aws_efs_file_system" = "efs-{basename}-{environment:short}-{region:short}",

    // AWS Database Resources
    "aws_rds_instance"   = "rds-{basename}-{environment:short}-{region:short}",
    "aws_rds_cluster"    = "rdsc-{basename}-{environment:short}-{region:short}",
    "aws_dynamodb_table" = "ddb-{basename}-{environment:short}-{region:short}",
    "aws_elasticache"    = "ec-{basename}-{environment:short}-{region:short}",

    // AWS Network Resources
    "aws_vpc"            = "vpc-{basename}-{environment:short}-{region:short}",
    "aws_subnet"         = "snet-{basename}-{environment:short}-{region:short}-{instance}",
    "aws_security_group" = "sg-{basename}-{environment:short}-{region:short}",
    "aws_route_table"    = "rt-{basename}-{environment:short}-{region:short}",
    "aws_elastic_ip"     = "eip-{basename}-{environment:short}",
    "aws_nat_gateway"    = "nat-{basename}-{environment:short}-{region:short}",
    "aws_load_balancer"  = "lb-{basename}-{environment:short}-{region:short}",
    "aws_target_group"   = "tg-{basename}-{environment:short}-{region:short}",

    // AWS Lambda Resources
    "aws_lambda_function" = "lambda-{basename}-{environment:short}-{region:short}",
    "aws_layer"           = "layer-{basename}-{environment:short}-{region:short}",

    // AWS Container Resources
    "aws_ecr_repository" = "ecr-{basename}-{environment:short}-{region:short}",
    "aws_ecs_cluster"    = "ecs-{basename}-{environment:short}-{region:short}",
    "aws_eks_cluster"    = "eks-{basename}-{environment:short}-{region:short}",

    // AWS IAM Resources
    "aws_iam_role"   = "role-{basename}-{environment:short}",
    "aws_iam_policy" = "pol-{basename}-{environment:short}",
    "aws_iam_user"   = "usr-{basename}-{environment:short}",
    "aws_iam_group"  = "grp-{basename}-{environment:short}",

    // AWS Monitoring Resources
    "aws_cloudwatch_alarm" = "cwa-{basename}-{environment:short}",
    "aws_log_group"        = "log-{basename}-{environment:short}",
    "aws_sns_topic"        = "sns-{basename}-{environment:short}-{region:short}",
    "aws_sqs_queue"        = "sqs-{basename}-{environment:short}-{region:short}",

    // AWS Application Resources
    "aws_api_gateway"   = "api-{basename}-{environment:short}-{region:short}",
    "aws_step_function" = "sf-{basename}-{environment:short}-{region:short}",
    "aws_cloudfront"    = "cf-{basename}-{environment:short}",

    // AWS Route53 Resources
    "aws_hosted_zone" = "hz-{basename}-{environment:short}",
    "aws_record_set"  = "rs-{basename}-{environment:short}",

    // GCP Compute Resources
    "google_compute_instance"          = "vm-{basename}-{environment:short}-{region:short}-{instance}",
    "google_compute_instance_group"    = "ig-{basename}-{environment:short}-{region:short}",
    "google_compute_instance_template" = "it-{basename}-{environment:short}-{region:short}",
    "google_compute_disk"              = "disk-{basename}-{environment:short}-{region:short}",
    "google_compute_snapshot"          = "snap-{basename}-{environment:short}-{region:short}",
    "google_compute_image"             = "img-{basename}-{environment:short}",

    // GCP Kubernetes Resources
    "google_container_cluster"   = "gke-{basename}-{environment:short}-{region:short}",
    "google_container_node_pool" = "np-{basename}-{environment:short}-{region:short}",

    // GCP Storage Resources
    "google_storage_bucket"     = "{basename}-{environment:short}-{region:short}",
    "google_filestore_instance" = "fs-{basename}-{environment:short}-{region:short}",

    // GCP Network Resources
    "google_compute_network"            = "vpc-{basename}-{environment:short}",
    "google_compute_subnetwork"         = "subnet-{basename}-{environment:short}-{region:short}",
    "google_compute_firewall"           = "fw-{basename}-{environment:short}",
    "google_compute_router"             = "router-{basename}-{environment:short}-{region:short}",
    "google_compute_address"            = "addr-{basename}-{environment:short}-{region:short}",
    "google_compute_global_address"     = "gaddr-{basename}-{environment:short}",
    "google_compute_forwarding_rule"    = "fr-{basename}-{environment:short}-{region:short}",
    "google_compute_target_http_proxy"  = "http-proxy-{basename}-{environment:short}",
    "google_compute_target_https_proxy" = "https-proxy-{basename}-{environment:short}",
    "google_compute_ssl_certificate"    = "cert-{basename}-{environment:short}",
    "google_compute_url_map"            = "url-map-{basename}-{environment:short}",
    "google_compute_backend_service"    = "bes-{basename}-{environment:short}",

    // GCP Database Resources
    "google_sql_database_instance" = "sql-{basename}-{environment:short}-{region:short}",
    "google_sql_database"          = "db-{basename}-{environment:short}",
    "google_bigtable_instance"     = "bt-{basename}-{environment:short}-{region:short}",
    "google_bigtable_table"        = "bt-tbl-{basename}-{environment:short}",
    "google_spanner_instance"      = "spanner-{basename}-{environment:short}",
    "google_spanner_database"      = "spanner-db-{basename}-{environment:short}",
    "google_firestore_database"    = "fs-db-{basename}-{environment:short}",

    // GCP Serverless Resources
    "google_cloudfunctions_function"         = "func-{basename}-{environment:short}-{region:short}",
    "google_cloud_run_service"               = "run-{basename}-{environment:short}-{region:short}",
    "google_app_engine_application"          = "app-{basename}-{environment:short}",
    "google_app_engine_standard_app_version" = "app-{basename}-{environment:short}-{version}",

    // GCP Data Analytics Resources
    "google_bigquery_dataset"    = "bq-ds-{basename}-{environment:short}",
    "google_bigquery_table"      = "bq-tbl-{basename}-{environment:short}",
    "google_dataflow_job"        = "df-{basename}-{environment:short}",
    "google_dataproc_cluster"    = "dp-{basename}-{environment:short}-{region:short}",
    "google_pubsub_topic"        = "ps-topic-{basename}-{environment:short}",
    "google_pubsub_subscription" = "ps-sub-{basename}-{environment:short}",

    // GCP IAM Resources
    "google_service_account"         = "sa-{basename}-{environment:short}",
    "google_project_iam_custom_role" = "role-{basename}-{environment:short}",

    // GCP Security Resources
    "google_kms_key_ring"          = "kr-{basename}-{environment:short}-{region:short}",
    "google_kms_crypto_key"        = "kms-{basename}-{environment:short}",
    "google_secret_manager_secret" = "secret-{basename}-{environment:short}",

    // GCP Monitoring Resources
    "google_monitoring_alert_policy"         = "alert-{basename}-{environment:short}",
    "google_logging_metric"                  = "log-{basename}-{environment:short}",
    "google_monitoring_notification_channel" = "notif-{basename}-{environment:short}",
    "google_monitoring_dashboard"            = "dash-{basename}-{environment:short}"

  }
}

# REMARK: Required initialization step to ensure the provider configuration is loaded during the ValidateConfig RPC
data "resourcenamingtool_status" "init" {}
