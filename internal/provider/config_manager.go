// Copyright (c) Thomas Geens

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gofrs/flock"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Global variables for the provider
var (
	providerSuffixPath = ".terraform/providers/registry.terraform.io/thomasgeens/resourcenamingtool"
	fileLockTimeout    = 10 * time.Second
	lockRetryInterval  = 50 * time.Millisecond
	globalConfigMutex  = &sync.Mutex{} // Memory-level lock for in-process synchronization
	// Define builtin default naming patterns following cloud provider best practices:
	// - Azure: Following Microsoft's Cloud Adoption Framework (CAF) naming conventions
	// - AWS: Following AWS Well-Architected Framework (WAF) and AWS service-specific naming guidelines
	// - GCP: Following Google Cloud's recommended naming conventions
	builtin_NamingPatterns = map[string]string{
		// Azure RM Resources - Based on CAF naming conventions
		// https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-naming
		// // Azure Core Resources
		// "azurerm_resource_group":         "rg-{basename}-{environment:short}-{region:short}",
		// "azurerm_virtual_network":        "vnet-{basename}-{environment:short}-{region:short}",
		// "azurerm_subnet":                 "snet-{basename}-{environment:short}-{instance}",
		// "azurerm_network_security_group": "nsg-{basename}-{environment:short}-{region:short}",
		// "azurerm_route_table":            "rt-{basename}-{environment:short}-{region:short}",

		// // Azure Compute Resources
		// "azurerm_virtual_machine":    "vm-{basename}-{environment:short}-{region:short}-{instance}",
		// "azurerm_availability_set":   "avs-{basename}-{environment:short}-{region:short}",
		// "azurerm_vm_scale_set":       "vmss-{basename}-{environment:short}-{region:short}",
		// "azurerm_kubernetes_cluster": "aks-{basename}-{environment:short}-{region:short}",

		// // Azure Storage Resources
		// "azurerm_storage_account":   "{basename}{environment:char}{region:char}{instance}",
		// "azurerm_storage_container": "sc-{basename}-{environment:short}",

		// // Azure Database Resources
		// "azurerm_sql_server":        "sql-{basename}-{environment:short}-{region:short}",
		// "azurerm_sql_database":      "sqldb-{basename}-{environment:short}",
		// "azurerm_cosmosdb_account":  "cosmos-{basename}-{environment:short}-{region:short}",
		// "azurerm_mysql_server":      "mysql-{basename}-{environment:short}-{region:short}",
		// "azurerm_postgresql_server": "psql-{basename}-{environment:short}-{region:short}",

		// // Azure App Resources
		// "azurerm_app_service":      "app-{basename}-{environment:short}-{region:short}",
		// "azurerm_app_service_plan": "plan-{basename}-{environment:short}-{region:short}",
		// "azurerm_function_app":     "func-{basename}-{environment:short}-{region:short}",

		// // Azure Security Resources
		// "azurerm_key_vault": "kv-{basename}-{environment:short}-{region:short}",

		// // Azure Integration Resources
		// "azurerm_servicebus_namespace": "sb-{basename}-{environment:short}-{region:short}",
		// "azurerm_eventhub_namespace":   "evh-{basename}-{environment:short}-{region:short}",
		// "azurerm_eventgrid_topic":      "evg-{basename}-{environment:short}-{region:short}",
		// "azurerm_logic_app_workflow":   "logic-{basename}-{environment:short}",

		// // Azure Container Resources
		// "azurerm_container_registry": "acr{basename}{environment:char}{region:char}",
		// "azurerm_container_group":    "aci-{basename}-{environment:short}",

		// // Azure Analytics Resources
		// "azurerm_log_analytics_workspace": "log-{basename}-{environment:short}-{region:short}",
		// "azurerm_application_insights":    "appi-{basename}-{environment:short}-{region:short}",

		// // Azure Network Resources
		// "azurerm_public_ip":           "pip-{basename}-{environment:short}-{region:short}",
		// "azurerm_lb":                  "lb-{basename}-{environment:short}-{region:short}",
		// "azurerm_application_gateway": "agw-{basename}-{environment:short}-{region:short}",
		// "azurerm_network_interface":   "nic-{basename}-{environment:short}",
		// "azurerm_private_endpoint":    "pe-{basename}-{environment:short}",

		// // Azure Identity Resources
		// "azurerm_user_assigned_identity": "id-{basename}-{environment:short}-{region:short}",

		// // Azure Monitor Resources
		// "azurerm_monitor_action_group": "ag-{basename}-{environment:short}",
		// "azurerm_monitor_metric_alert": "ar-{basename}-{environment:short}",

		// // AWS Resources - Following AWS Well-Architected Framework and service-specific naming guidelines
		// // https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
		// // https://docs.aws.amazon.com/whitepapers/latest/tagging-best-practices/tagging-best-practices.html
		// // AWS Compute Resources
		// "aws_ec2_instance":       "ec2-{basename}-{environment:short}-{region:short}-{instance}",
		// "aws_auto_scaling_group": "asg-{basename}-{environment:short}-{region:short}",
		// "aws_launch_template":    "lt-{basename}-{environment:short}-{region:short}",

		// // AWS Storage Resources
		// "aws_s3_bucket":       "{basename}-{environment:short}-{region:short}-{instance}",
		// "aws_efs_file_system": "efs-{basename}-{environment:short}-{region:short}",

		// // AWS Database Resources
		// "aws_rds_instance":   "rds-{basename}-{environment:short}-{region:short}",
		// "aws_rds_cluster":    "rdsc-{basename}-{environment:short}-{region:short}",
		// "aws_dynamodb_table": "ddb-{basename}-{environment:short}-{region:short}",
		// "aws_elasticache":    "ec-{basename}-{environment:short}-{region:short}",

		// // AWS Network Resources
		// "aws_vpc":            "vpc-{basename}-{environment:short}-{region:short}",
		// "aws_subnet":         "snet-{basename}-{environment:short}-{region:short}-{instance}",
		// "aws_security_group": "sg-{basename}-{environment:short}-{region:short}",
		// "aws_route_table":    "rt-{basename}-{environment:short}-{region:short}",
		// "aws_elastic_ip":     "eip-{basename}-{environment:short}",
		// "aws_nat_gateway":    "nat-{basename}-{environment:short}-{region:short}",
		// "aws_load_balancer":  "lb-{basename}-{environment:short}-{region:short}",
		// "aws_target_group":   "tg-{basename}-{environment:short}-{region:short}",

		// // AWS Lambda Resources
		// "aws_lambda_function": "lambda-{basename}-{environment:short}-{region:short}",
		// "aws_layer":           "layer-{basename}-{environment:short}-{region:short}",

		// // AWS Container Resources
		// "aws_ecr_repository": "ecr-{basename}-{environment:short}-{region:short}",
		// "aws_ecs_cluster":    "ecs-{basename}-{environment:short}-{region:short}",
		// "aws_eks_cluster":    "eks-{basename}-{environment:short}-{region:short}",

		// // AWS IAM Resources
		// "aws_iam_role":   "role-{basename}-{environment:short}",
		// "aws_iam_policy": "pol-{basename}-{environment:short}",
		// "aws_iam_user":   "usr-{basename}-{environment:short}",
		// "aws_iam_group":  "grp-{basename}-{environment:short}",

		// // AWS Monitoring Resources
		// "aws_cloudwatch_alarm": "cwa-{basename}-{environment:short}",
		// "aws_log_group":        "log-{basename}-{environment:short}",
		// "aws_sns_topic":        "sns-{basename}-{environment:short}-{region:short}",
		// "aws_sqs_queue":        "sqs-{basename}-{environment:short}-{region:short}",

		// // AWS Application Resources
		// "aws_api_gateway":   "api-{basename}-{environment:short}-{region:short}",
		// "aws_step_function": "sf-{basename}-{environment:short}-{region:short}",
		// "aws_cloudfront":    "cf-{basename}-{environment:short}",

		// // AWS Route53 Resources
		// "aws_hosted_zone": "hz-{basename}-{environment:short}",
		// "aws_record_set":  "rs-{basename}-{environment:short}",

		// // Google Cloud Platform Resources - Following Google Cloud's naming conventions
		// // https://cloud.google.com/resource-manager/docs/creating-managing-projects#identifying_projects
		// // https://cloud.google.com/architecture/best-practices-naming-resources
		// // GCP Compute Resources
		// "google_compute_instance":          "vm-{basename}-{environment:short}-{region:short}-{instance}",
		// "google_compute_instance_group":    "ig-{basename}-{environment:short}-{region:short}",
		// "google_compute_instance_template": "it-{basename}-{environment:short}-{region:short}",
		// "google_compute_disk":              "disk-{basename}-{environment:short}-{region:short}",
		// "google_compute_snapshot":          "snap-{basename}-{environment:short}-{region:short}",
		// "google_compute_image":             "img-{basename}-{environment:short}",

		// // GCP Kubernetes Resources
		// "google_container_cluster":   "gke-{basename}-{environment:short}-{region:short}",
		// "google_container_node_pool": "np-{basename}-{environment:short}-{region:short}",

		// // GCP Storage Resources
		// "google_storage_bucket":     "{basename}-{environment:short}-{region:short}",
		// "google_filestore_instance": "fs-{basename}-{environment:short}-{region:short}",

		// // GCP Network Resources
		// "google_compute_network":            "vpc-{basename}-{environment:short}",
		// "google_compute_subnetwork":         "subnet-{basename}-{environment:short}-{region:short}",
		// "google_compute_firewall":           "fw-{basename}-{environment:short}",
		// "google_compute_router":             "router-{basename}-{environment:short}-{region:short}",
		// "google_compute_address":            "addr-{basename}-{environment:short}-{region:short}",
		// "google_compute_global_address":     "gaddr-{basename}-{environment:short}",
		// "google_compute_forwarding_rule":    "fr-{basename}-{environment:short}-{region:short}",
		// "google_compute_target_http_proxy":  "http-proxy-{basename}-{environment:short}",
		// "google_compute_target_https_proxy": "https-proxy-{basename}-{environment:short}",
		// "google_compute_ssl_certificate":    "cert-{basename}-{environment:short}",
		// "google_compute_url_map":            "url-map-{basename}-{environment:short}",
		// "google_compute_backend_service":    "bes-{basename}-{environment:short}",

		// // GCP Database Resources
		// "google_sql_database_instance": "sql-{basename}-{environment:short}-{region:short}",
		// "google_sql_database":          "db-{basename}-{environment:short}",
		// "google_bigtable_instance":     "bt-{basename}-{environment:short}-{region:short}",
		// "google_bigtable_table":        "bt-tbl-{basename}-{environment:short}",
		// "google_spanner_instance":      "spanner-{basename}-{environment:short}",
		// "google_spanner_database":      "spanner-db-{basename}-{environment:short}",
		// "google_firestore_database":    "fs-db-{basename}-{environment:short}",

		// // GCP Serverless Resources
		// "google_cloudfunctions_function":         "func-{basename}-{environment:short}-{region:short}",
		// "google_cloud_run_service":               "run-{basename}-{environment:short}-{region:short}",
		// "google_app_engine_application":          "app-{basename}-{environment:short}",
		// "google_app_engine_standard_app_version": "app-{basename}-{environment:short}-{version}",

		// // GCP Data Analytics Resources
		// "google_bigquery_dataset":    "bq-ds-{basename}-{environment:short}",
		// "google_bigquery_table":      "bq-tbl-{basename}-{environment:short}",
		// "google_dataflow_job":        "df-{basename}-{environment:short}",
		// "google_dataproc_cluster":    "dp-{basename}-{environment:short}-{region:short}",
		// "google_pubsub_topic":        "ps-topic-{basename}-{environment:short}",
		// "google_pubsub_subscription": "ps-sub-{basename}-{environment:short}",

		// // GCP IAM Resources
		// "google_service_account":         "sa-{basename}-{environment:short}",
		// "google_project_iam_custom_role": "role-{basename}-{environment:short}",

		// // GCP Security Resources
		// "google_kms_key_ring":          "kr-{basename}-{environment:short}-{region:short}",
		// "google_kms_crypto_key":        "kms-{basename}-{environment:short}",
		// "google_secret_manager_secret": "secret-{basename}-{environment:short}",

		// // GCP Monitoring Resources
		// "google_monitoring_alert_policy":         "alert-{basename}-{environment:short}",
		// "google_logging_metric":                  "log-{basename}-{environment:short}",
		// "google_monitoring_notification_channel": "notif-{basename}-{environment:short}",
		// "google_monitoring_dashboard":            "dash-{basename}-{environment:short}",
	}
)

// processComponentFromMap extracts component values from a map and creates a ComponentValueObject
// This reduces duplication in the code that processes component values
func processComponentFromMap(ctx context.Context, componentMap map[string]interface{}) (ComponentValueObject, bool) {
	fullname, _ := componentMap["fullname"].(string)
	shortcode, _ := componentMap["shortcode"].(string)
	char, _ := componentMap["char"].(string)
	componentValue, diags := CreateComponentValueObjectFromParts(ctx, fullname, shortcode, char)
	return componentValue, !diags.HasError()
}

// getConfigPath returns the standard configuration file path
// This ensures consistent path resolution across all functions
func getConfigPath(ctx context.Context) string {
	logDebug(ctx, "Invoking getConfigPath")
	workDir, err := os.Getwd()
	configDir := ""
	if err != nil {
		logDebug(ctx, "getConfigPath: Failed to get working directory path, default to temp directory path: %s", err)
		configDir = filepath.Join(os.TempDir(), providerSuffixPath)
	} else {
		logDebug(ctx, "getConfigPath: Working directory path: %s", workDir)
		configDir = filepath.Join(workDir, providerSuffixPath)
	}
	logDebug(ctx, "getConfigPath: Config directory path: %s", configDir)
	return filepath.Join(configDir, "provider-config.json")
}

// ensureConfigDirExists ensures the configuration directory exists
// This is needed before any lock operations can be performed
func ensureConfigDirExists(configPath string) error {
	dirPath := filepath.Dir(configPath)
	return os.MkdirAll(dirPath, 0750)
}

// getLockFilePath returns the path to the lock file for the configuration
func getLockFilePath(configPath string) string {
	return configPath + ".lock"
}

// GetSharedProviderConfig retrieves the provider configuration from the file
// This allows functions to access the provider configuration between different process invocations
func GetSharedProviderConfig(ctx context.Context) *resourcenamingtoolProviderModel {
	// Use the provided context rather than creating a new one
	logDebug(ctx, "GetSharedProviderConfig: Starting configuration retrieval")

	// Get the standard config path
	configPath := getConfigPath(ctx)

	logDebug(ctx, "GetSharedProviderConfig: Checking file path: %s", configPath)

	// Ensure the config directory exists before attempting to acquire a lock
	if err := ensureConfigDirExists(configPath); err != nil {
		logError(ctx, "GetSharedProviderConfig: Error creating config directory: %s", err)
		return nil
	}

	// Use in-memory mutex first (for same-process synchronization)
	globalConfigMutex.Lock()
	// Using helper function to unlock and log when the function returns
	defer unlockMutexAndLog(globalConfigMutex, ctx, "GetSharedProviderConfig")

	// Create a file lock for cross-process synchronization
	fileLock := flock.New(getLockFilePath(configPath))
	locked, err := tryLockWithRetries(fileLock, fileLockTimeout, lockRetryInterval)
	if err != nil {
		logError(ctx, "GetSharedProviderConfig: Error acquiring file lock: %s", err)
		return nil
	}
	if !locked {
		logError(ctx, "GetSharedProviderConfig: Could not acquire lock within timeout period (%s)", fileLockTimeout)
		return nil
	}
	// Using helper function to unlock and log when the function returns
	defer unlockAndLog(fileLock, ctx, "GetSharedProviderConfig")

	fileConfig := loadProviderConfigFromFile(ctx)

	if fileConfig != nil {
		logDebug(ctx, "GetSharedProviderConfig: Successfully loaded config from file")
		return fileConfig
	}

	logDebug(ctx, "GetSharedProviderConfig: No shared provider configuration found")
	return nil
}

// SaveSharedProviderConfig saves the provider configuration to a file
// This allows it to be shared between different process invocations
func SaveSharedProviderConfig(ctx context.Context, config *resourcenamingtoolProviderModel) error {
	// Use the provided context rather than creating a new one
	logDebug(ctx, "SaveSharedProviderConfig: Starting save operation")

	// Get the standard config path
	configPath := getConfigPath(ctx)

	logDebug(ctx, "SaveSharedProviderConfig: Saving to path: %s", configPath)

	// Ensure the config directory exists before attempting to use a lock
	if err := ensureConfigDirExists(configPath); err != nil {
		logError(ctx, "SaveSharedProviderConfig: Error creating config directory: %s", err)
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Use in-memory mutex first (for same-process synchronization)
	globalConfigMutex.Lock()
	// Using helper function to unlock and log when the function returns
	defer unlockMutexAndLog(globalConfigMutex, ctx, "SaveSharedProviderConfig")

	// Create a file lock for cross-process synchronization
	fileLock := flock.New(getLockFilePath(configPath))
	locked, err := tryLockWithRetries(fileLock, fileLockTimeout, lockRetryInterval)
	if err != nil {
		logError(ctx, "SaveSharedProviderConfig: Error acquiring file lock: %s", err)
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !locked {
		logError(ctx, "SaveSharedProviderConfig: Could not acquire lock within timeout period (%s)", fileLockTimeout)
		return fmt.Errorf("timeout acquiring lock")
	}
	// Using helper function to unlock and log when the function returns
	defer unlockAndLog(fileLock, ctx, "SaveSharedProviderConfig")

	return saveProviderConfigToFile(ctx, config)
}

// saveProviderConfigToFile persists the provider configuration to a file
// so it can be shared across different process invocations
func saveProviderConfigToFile(ctx context.Context, config *resourcenamingtoolProviderModel) error {
	logDebug(ctx, "Invoking saveProviderConfigToFile")

	if config == nil {
		logError(ctx, "cannot save nil configuration")
		return fmt.Errorf("cannot save nil configuration")
	}

	// Get the standard config path
	configPath := getConfigPath(ctx)
	tempDir := filepath.Dir(configPath)

	// Ensure the config directory exists before attempting to acquire a lock
	if err := ensureConfigDirExists(configPath); err != nil {
		logError(ctx, "saveProviderConfigToFile: Error creating config directory: %s", err.Error())
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Use in-memory mutex first (for same-process synchronization)
	globalConfigMutex.Lock()
	// Using helper function to unlock and log when the function returns
	defer unlockMutexAndLog(globalConfigMutex, ctx, "saveProviderConfigToFile")

	// Create a file lock for cross-process synchronization
	fileLock := flock.New(getLockFilePath(configPath))
	locked, err := tryLockWithRetries(fileLock, fileLockTimeout, lockRetryInterval)
	if err != nil {
		logError(ctx, "saveProviderConfigToFile: Error acquiring file lock: %s", err.Error())
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !locked {
		logError(ctx, "saveProviderConfigToFile: Could not acquire lock within timeout period (%s)", fileLockTimeout)
		return fmt.Errorf("timeout acquiring lock")
	}
	// Using helper function to unlock and log when the function returns
	defer unlockAndLog(fileLock, ctx, "saveProviderConfigToFile")

	logDebug(ctx, "Attempting to save provider config to: %s", configPath)
	logDebug(ctx, "Configuration directory path: %s", tempDir)

	// Marshal the configuration to JSON using the struct tags
	configJson, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		logError(ctx, "Failed to marshal configuration to JSON: %s", err.Error())
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Log the JSON content (truncated for readability)
	jsonPreview := string(configJson)
	if len(jsonPreview) > 2000 {
		jsonPreview = jsonPreview[:2000] + "..."
	}
	logDebug(ctx, "Configuration JSON to be written (preview): %s", jsonPreview)

	// Save to a file
	logDebug(ctx, "Writing configuration to file: %s", configPath)
	if err := os.WriteFile(configPath, configJson, 0600); err != nil {
		logError(ctx, "Failed to write configuration to file %s: %s", configPath, err.Error())
		return err
	}

	// Verify the file was written by reading it back
	// Ignore gosec G304: Potential file inclusion via variable
	//#nosec G304
	if readBytes, err := os.ReadFile(configPath); err == nil {
		logDebug(ctx, "Successfully verified file was written: size=%d, matches=%v",
			len(readBytes), len(readBytes) == len(configJson))
	} else {
		logError(ctx, "Failed to verify file was written: %s", err.Error())
	}

	logDebug(ctx, "Successfully wrote configuration to file: %s", configPath)
	return nil
}

// loadProviderConfigFromFile loads the provider configuration from a file
// allowing it to be shared across different process invocations
func loadProviderConfigFromFile(ctx context.Context) *resourcenamingtoolProviderModel {
	logDebug(ctx, "Invoking loadProviderConfigFromFile")

	// Get the standard config path
	configPath := getConfigPath(ctx)
	logDebug(ctx, "Attempting to load provider config from: %s", configPath)

	// Check if the file exists
	if fileInfo, err := os.Stat(configPath); os.IsNotExist(err) {
		logDebug(ctx, "Config file does not exist: %s (%s)", configPath, err.Error())
		return nil
	} else if err != nil {
		logError(ctx, "Error checking config file %s: %s", configPath, err.Error())
		return nil
	} else {
		logDebug(ctx, "Found config file: %s (size=%d, modified=%s, permissions=%s)",
			configPath, fileInfo.Size(), fileInfo.ModTime().String(), fileInfo.Mode().String())
	}

	// Read the file
	// Ignore gosec G304: Potential file inclusion via variable
	//#nosec G304
	configJson, err := os.ReadFile(configPath)
	if err != nil {
		logError(ctx, "Failed to read config file %s: %s", configPath, err.Error())
		return nil
	}

	// Debug log the JSON content (truncated for readability)
	jsonPreview := string(configJson)
	if len(jsonPreview) > 200 {
		jsonPreview = jsonPreview[:200] + "..." // Truncate if too long
	}
	logDebug(ctx, "Read configuration JSON from file (length=%d): %s",
		len(configJson), jsonPreview)

	// Unmarshal the JSON directly to the provider model using the struct tags
	config := &resourcenamingtoolProviderModel{}
	if err := json.Unmarshal(configJson, config); err != nil {
		logError(ctx, "Failed to unmarshal configuration from JSON: %s", err.Error())
		return nil
	}

	logDebug(ctx, "Successfully unmarshaled provider config using JSON struct tags")

	// Handle the AdditionalComponents and AdditionalNamingPatterns maps specially
	// since they require additional conversion
	logDebug(ctx, "Converting maps from JSON representation to internal types...")

	// Parse the JSON into a map to extract AdditionalComponents and AdditionalNamingPatterns
	var rawConfig map[string]interface{}
	if err := json.Unmarshal(configJson, &rawConfig); err != nil {
		logError(ctx, "Failed to parse raw config: %s", err.Error())
		return nil
	}

	// Handle AdditionalComponents
	if components, ok := rawConfig["AdditionalComponents"].(map[string]interface{}); ok && len(components) > 0 {
		elements := make(map[string]attr.Value)
		for k, v := range components {
			if component, ok := v.(map[string]interface{}); ok {
				if componentValue, ok := processComponentFromMap(ctx, component); ok {
					elements[k] = componentValue
				}
			} else {
				logError(ctx, "Additional component '%s' is not a map[string]interface{}: %T", k, v)
			}
		}
		componentsMap, diags := types.MapValueFrom(ctx, NewComponentValueType(), elements)
		if !diags.HasError() {
			config.AdditionalComponents = componentsMap
		}
	}

	// Handle the AdditionalNamingPatterns map
	// Always initialize the map even if empty from the JSON
	elements := make(map[string]attr.Value)

	if patterns, ok := rawConfig["AdditionalNamingPatterns"].(map[string]interface{}); ok {
		logDebug(ctx, "Found AdditionalNamingPatterns in config JSON with %d entries", len(patterns))

		for k, v := range patterns {
			if strVal, ok := v.(string); ok {
				elements[k] = types.StringValue(strVal)
				logDebug(ctx, "Loaded naming pattern for resource type: %s = %s", k, strVal)
			}
		}
	} else {
		logDebug(ctx, "AdditionalNamingPatterns was not found in config JSON or is not a map")
	}

	// Always create the map even if empty
	patternsMap, diags := types.MapValueFrom(ctx, types.StringType, elements)
	if !diags.HasError() {
		config.AdditionalNamingPatterns = patternsMap
		logDebug(ctx, "Set AdditionalNamingPatterns with %d elements", len(elements))
	} else {
		logError(ctx, "Failed to create AdditionalNamingPatterns map: %s", diags)
	}

	// Log configuration loading complete
	logDebug(ctx, "Successfully loaded provider config")
	return config
}

// unlockAndLog unlocks the file lock and logs a message when the unlock happens
func unlockAndLog(lock *flock.Flock, ctx context.Context, functionName string) {
	if err := lock.Unlock(); err != nil {
		logError(ctx, "%s: Error releasing file lock: %s", functionName, err.Error())
	} else {
		logDebug(ctx, "%s: File lock has been released", functionName)
	}
}

// unlockMutexAndLog unlocks the mutex and logs a message when the unlock happens
func unlockMutexAndLog(mutex *sync.Mutex, ctx context.Context, functionName string) {
	mutex.Unlock()
	logDebug(ctx, "%s: Memory mutex has been released", functionName)
}

// tryLockWithRetries attempts to acquire a file lock with retries
// It will retry at the specified interval until it gets the lock or the timeout expires
func tryLockWithRetries(fileLock *flock.Flock, timeout time.Duration, retryInterval time.Duration) (bool, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Try to acquire the lock with retries until the context timeout
	startTime := time.Now()
	for {
		// Check if the context is done (timed out)
		select {
		case <-ctx.Done():
			return false, nil // Timed out
		default:
			// Try to acquire the lock
			locked, err := fileLock.TryLock()
			if err != nil {
				return false, err // Error acquiring lock
			}
			if locked {
				return true, nil // Successfully acquired lock
			}

			// Calculate elapsed time for logging
			elapsed := time.Since(startTime)
			logDebug(ctx, "Lock attempt failed, retrying. Elapsed: %v, Timeout: %v", elapsed, timeout)

			// Wait before retrying
			select {
			case <-ctx.Done():
				return false, nil // Timed out during wait
			case <-time.After(retryInterval):
				// Continue with retry
			}
		}
	}
}
