// Copyright (c) Thomas Geens

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"resourcenamingtool": providerserver.NewProtocol6WithError(New("test")()),
}

// providerConfig is a shared configuration to combine with the actual
// test configuration so the provider is properly configured.
const providerConfig = `
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
		"azurerm_resource_group":         "rg-{basename}-{environment:short}-{region:short}",
		"azurerm_virtual_network":        "vnet-{basename}-{environment:short}-{region:short}",
		"azurerm_subnet":                 "snet-{basename}-{environment:short}-{instance}",
		"azurerm_network_security_group": "nsg-{basename}-{environment:short}-{region:short}",
		"azurerm_route_table":            "rt-{basename}-{environment:short}-{region:short}",

		// Azure Compute Resources
		"azurerm_virtual_machine":    "vm-{basename}-{environment:short}-{region:short}-{instance}",
		"azurerm_availability_set":   "avs-{basename}-{environment:short}-{region:short}",
		"azurerm_vm_scale_set":       "vmss-{basename}-{environment:short}-{region:short}",
		"azurerm_kubernetes_cluster": "aks-{basename}-{environment:short}-{region:short}",

		// Azure Storage Resources
		"azurerm_storage_account":   "{basename}{environment:char}{region:char}{instance}",
		"azurerm_storage_container": "sc-{basename}-{environment:short}",

		// Azure Database Resources
		"azurerm_sql_server":        "sql-{basename}-{environment:short}-{region:short}",
		"azurerm_sql_database":      "sqldb-{basename}-{environment:short}",
		"azurerm_cosmosdb_account":  "cosmos-{basename}-{environment:short}-{region:short}",
		"azurerm_mysql_server":      "mysql-{basename}-{environment:short}-{region:short}",
		"azurerm_postgresql_server": "psql-{basename}-{environment:short}-{region:short}",

		// Azure App Resources
		"azurerm_app_service":      "app-{basename}-{environment:short}-{region:short}",
		"azurerm_app_service_plan": "plan-{basename}-{environment:short}-{region:short}",
		"azurerm_function_app":     "func-{basename}-{environment:short}-{region:short}",

		// Azure Security Resources
		"azurerm_key_vault": "kv-{basename}-{environment:short}-{region:short}",

		// Azure Integration Resources
		"azurerm_servicebus_namespace": "sb-{basename}-{environment:short}-{region:short}",
		"azurerm_eventhub_namespace":   "evh-{basename}-{environment:short}-{region:short}",
		"azurerm_eventgrid_topic":      "evg-{basename}-{environment:short}-{region:short}",
		"azurerm_logic_app_workflow":   "logic-{basename}-{environment:short}",

		// Azure Container Resources
		"azurerm_container_registry": "acr{basename}{environment:char}{region:char}",
		"azurerm_container_group":    "aci-{basename}-{environment:short}",

		// Azure Analytics Resources
		"azurerm_log_analytics_workspace": "log-{basename}-{environment:short}-{region:short}",
		"azurerm_application_insights":    "appi-{basename}-{environment:short}-{region:short}",

		// Azure Network Resources
		"azurerm_public_ip":           "pip-{basename}-{environment:short}-{region:short}",
		"azurerm_lb":                  "lb-{basename}-{environment:short}-{region:short}",
		"azurerm_application_gateway": "agw-{basename}-{environment:short}-{region:short}",
		"azurerm_network_interface":   "nic-{basename}-{environment:short}",
		"azurerm_private_endpoint":    "pe-{basename}-{environment:short}",

		// Azure Identity Resources
		"azurerm_user_assigned_identity": "id-{basename}-{environment:short}-{region:short}",

		// Azure Monitor Resources
		"azurerm_monitor_action_group": "ag-{basename}-{environment:short}",
		"azurerm_monitor_metric_alert": "ar-{basename}-{environment:short}",
  }
}

data "resourcenamingtool_status" "init" {} // Required initialization step to ensure the provider configuration is loaded

`
