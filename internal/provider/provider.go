// Copyright (c) Thomas Geens

package provider

import (
	"context"
	_ "embed" // Import the embed package
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ provider.Provider = &ResourceNamingToolProvider{}

// New returns a new instance of the resourcenamingtool provider
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &resourcenamingtoolFunctionsProvider{version: version}
	}
}

type resourcenamingtoolFunctionsProvider struct {
	provider.ProviderWithValidateConfig
	provider.ProviderWithFunctions
	version string
	config  *resourcenamingtoolProviderModel
}

func (p *resourcenamingtoolFunctionsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	logDebug(ctx, "Loading resourcenamingtool metadata...")
	// Debug log all provider metadata
	logDebugWithFields(ctx, "Provider metadata", map[string]interface{}{
		"MetadataRequest": req,
	})
	resp.TypeName = "resourcenamingtool"
	resp.Version = p.version
}

//go:embed descriptions/provider_description.txt
var providerDescription string

//go:embed descriptions/provider_markdown_description.md
var providerMarkdownDescription string

func (p *resourcenamingtoolFunctionsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// Provider instance identification
			"provider_instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "A unique identifier for this provider instance. Used to avoid configuration file conflicts when using multiple provider instances in the same Terraform configuration.",
			},

			// Core components
			"default_resource_type": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default resource type to use when not provided in the function call. This corresponds to the type of resource being created (e.g., 'virtual_machine', 'storage_account'). The resource type determines which naming pattern is used.",
			},
			"default_resource_prefix": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default resource prefix to use when not provided in the function call. This is an optional prefix that goes before the resource type abbreviation in the name pattern (e.g., 'shared', 'core').",
			},
			"default_basename": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default basename to use when not provided in the function call. This is the core identifying name of the resource that will be part of all resource names (e.g., 'webapp', 'payroll').",
			},
			"default_environment": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default environment to use when not provided in the function call. Represents the deployment environment (e.g., 'dev', 'test', 'prod'). Used to distinguish resources across different environments.",
			},
			"default_region": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default region to use when not provided in the function call. Represents the cloud region where the resource is deployed (e.g., 'eastus', 'westeurope', 'us-west-2'). Often used in naming patterns to distinguish resources across regions.",
			},
			"default_instance": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default instance identifier to use when not provided in the function call. Used to distinguish between multiple instances of the same resource type (e.g., '01', '02'). Commonly used for resources that are deployed in multiples.",
			},

			// Organization related components
			"default_organization": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default organization to use when not provided in the function call. Identifies the overall organization owning the resource (e.g., 'contoso', 'fabrikam').",
			},
			"default_business_unit": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default business unit to use when not provided in the function call. Identifies the business unit within the organization (e.g., 'finance', 'hr', 'it').",
			},
			"default_cost_center": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default cost center to use when not provided in the function call. Identifies the financial cost center associated with the resource (e.g., 'cc123', 'marketing').",
			},
			"default_project": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default project to use when not provided in the function call. Identifies the project associated with the resource (e.g., 'website-redesign', 'data-migration').",
			},
			"default_application": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default application name to use when not provided in the function call. Identifies the application using the resource (e.g., 'inventory-system', 'crm').",
			},
			"default_workload": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default workload type to use when not provided in the function call. Describes the function or purpose of the workload (e.g., 'api', 'web', 'batch').",
			},

			// Provider specific components
			"default_subscription": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default subscription to use when not provided in the function call. Primarily used with Azure to identify the subscription context (e.g., 'prod', 'dev', 'subscription-name').",
			},
			"default_location": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default location to use when not provided in the function call. Alternative to region, used in some cloud providers (e.g., 'eastus', 'westeurope').",
			},
			"default_domain": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default domain to use when not provided in the function call. Used for resources that require a domain name (e.g., 'contoso.com', 'fabrikam.net').",
			},
			"default_criticality": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default criticality level to use when not provided in the function call. Indicates the importance or criticality of the resource (e.g., 'high', 'medium', 'low', 'mission-critical').",
			},

			// Initiative/solution related
			"default_initiative": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default initiative to use when not provided in the function call. Identifies a broader business initiative the resource belongs to (e.g., 'cloud-migration', 'security-enhancement').",
			},
			"default_solution": schema.ObjectAttribute{
				CustomType:  NewComponentValueType(),
				Optional:    true,
				Description: "Default solution to use when not provided in the function call. Identifies the solution architecture or pattern the resource is part of (e.g., 'microservices', 'data-lake').",
			},

			// Extension points
			"additional_components": schema.MapAttribute{
				ElementType: NewComponentValueType(),
				Optional:    true,
				Description: "Additional custom components to use in resource name generation. Keys must be wrapped in curly braces to be used in patterns (e.g., {custom_component1}, {department}). These can be used in custom naming patterns.",
			},
			"additional_naming_patterns": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Additional naming patterns for specific resource types. Overrides built-in patterns. Keys should match resource types and values should contain component placeholders (e.g., \"my_custom_resource\": \"prefix-{basename}-{environment:short}\").",
			},
		},
		Description:         providerDescription,
		MarkdownDescription: providerMarkdownDescription,
	}
}

// Implement provider data model
type resourcenamingtoolProviderModel struct {
	// Instance identifier for the provider
	ProviderInstanceID types.String `tfsdk:"provider_instance_id" json:"provider_instance_id,omitempty"`

	// Core components
	DefaultResourceType   ComponentValueObject `tfsdk:"default_resource_type" json:"DefaultResourceType,omitempty"`
	DefaultResourcePrefix ComponentValueObject `tfsdk:"default_resource_prefix" json:"DefaultResourcePrefix,omitempty"`
	DefaultBasename       ComponentValueObject `tfsdk:"default_basename" json:"DefaultBasename,omitempty"`
	DefaultEnvironment    ComponentValueObject `tfsdk:"default_environment" json:"DefaultEnvironment,omitempty"`
	DefaultRegion         ComponentValueObject `tfsdk:"default_region" json:"DefaultRegion,omitempty"`
	DefaultInstance       ComponentValueObject `tfsdk:"default_instance" json:"DefaultInstance,omitempty"`

	// Organization related components
	DefaultOrganization ComponentValueObject `tfsdk:"default_organization" json:"DefaultOrganization,omitempty"`
	DefaultBusinessUnit ComponentValueObject `tfsdk:"default_business_unit" json:"DefaultBusinessUnit,omitempty"`
	DefaultCostCenter   ComponentValueObject `tfsdk:"default_cost_center" json:"DefaultCostCenter,omitempty"`
	DefaultProject      ComponentValueObject `tfsdk:"default_project" json:"DefaultProject,omitempty"`
	DefaultApplication  ComponentValueObject `tfsdk:"default_application" json:"DefaultApplication,omitempty"`
	DefaultWorkload     ComponentValueObject `tfsdk:"default_workload" json:"DefaultWorkload,omitempty"`

	// Provider specific components
	DefaultSubscription ComponentValueObject `tfsdk:"default_subscription" json:"DefaultSubscription,omitempty"`
	DefaultLocation     ComponentValueObject `tfsdk:"default_location" json:"DefaultLocation,omitempty"`
	DefaultDomain       ComponentValueObject `tfsdk:"default_domain" json:"DefaultDomain,omitempty"`
	DefaultCriticality  ComponentValueObject `tfsdk:"default_criticality" json:"DefaultCriticality,omitempty"`

	// Initiative/solution related
	DefaultInitiative ComponentValueObject `tfsdk:"default_initiative" json:"DefaultInitiative,omitempty"`
	DefaultSolution   ComponentValueObject `tfsdk:"default_solution" json:"DefaultSolution,omitempty"`

	// Extension points
	AdditionalComponents     types.Map `tfsdk:"additional_components" json:"AdditionalComponents,omitempty"`
	AdditionalNamingPatterns types.Map `tfsdk:"additional_naming_patterns" json:"AdditionalNamingPatterns,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for resourcenamingtoolProviderModel
// This is needed to properly marshal types.Map fields which aren't automatically handled by the standard JSON marshaller
func (m resourcenamingtoolProviderModel) MarshalJSON() ([]byte, error) {
	// Create a map to store all the serialized fields
	output := make(map[string]interface{})

	// Handle provider instance ID
	if !m.ProviderInstanceID.IsNull() && !m.ProviderInstanceID.IsUnknown() {
		output["provider_instance_id"] = m.ProviderInstanceID.ValueString()
	}

	// Handle all component value objects
	// Core components
	if !m.DefaultResourceType.IsNull() && !m.DefaultResourceType.IsUnknown() {
		output["DefaultResourceType"] = m.DefaultResourceType
	}
	if !m.DefaultResourcePrefix.IsNull() && !m.DefaultResourcePrefix.IsUnknown() {
		output["DefaultResourcePrefix"] = m.DefaultResourcePrefix
	}
	if !m.DefaultBasename.IsNull() && !m.DefaultBasename.IsUnknown() {
		output["DefaultBasename"] = m.DefaultBasename
	}
	if !m.DefaultEnvironment.IsNull() && !m.DefaultEnvironment.IsUnknown() {
		output["DefaultEnvironment"] = m.DefaultEnvironment
	}
	if !m.DefaultRegion.IsNull() && !m.DefaultRegion.IsUnknown() {
		output["DefaultRegion"] = m.DefaultRegion
	}
	if !m.DefaultInstance.IsNull() && !m.DefaultInstance.IsUnknown() {
		output["DefaultInstance"] = m.DefaultInstance
	}

	// Organization related components
	if !m.DefaultOrganization.IsNull() && !m.DefaultOrganization.IsUnknown() {
		output["DefaultOrganization"] = m.DefaultOrganization
	}
	if !m.DefaultBusinessUnit.IsNull() && !m.DefaultBusinessUnit.IsUnknown() {
		output["DefaultBusinessUnit"] = m.DefaultBusinessUnit
	}
	if !m.DefaultCostCenter.IsNull() && !m.DefaultCostCenter.IsUnknown() {
		output["DefaultCostCenter"] = m.DefaultCostCenter
	}
	if !m.DefaultProject.IsNull() && !m.DefaultProject.IsUnknown() {
		output["DefaultProject"] = m.DefaultProject
	}
	if !m.DefaultApplication.IsNull() && !m.DefaultApplication.IsUnknown() {
		output["DefaultApplication"] = m.DefaultApplication
	}
	if !m.DefaultWorkload.IsNull() && !m.DefaultWorkload.IsUnknown() {
		output["DefaultWorkload"] = m.DefaultWorkload
	}

	// Provider specific components
	if !m.DefaultSubscription.IsNull() && !m.DefaultSubscription.IsUnknown() {
		output["DefaultSubscription"] = m.DefaultSubscription
	}
	if !m.DefaultLocation.IsNull() && !m.DefaultLocation.IsUnknown() {
		output["DefaultLocation"] = m.DefaultLocation
	}
	if !m.DefaultDomain.IsNull() && !m.DefaultDomain.IsUnknown() {
		output["DefaultDomain"] = m.DefaultDomain
	}
	if !m.DefaultCriticality.IsNull() && !m.DefaultCriticality.IsUnknown() {
		output["DefaultCriticality"] = m.DefaultCriticality
	}

	// Initiative/solution related
	if !m.DefaultInitiative.IsNull() && !m.DefaultInitiative.IsUnknown() {
		output["DefaultInitiative"] = m.DefaultInitiative
	}
	if !m.DefaultSolution.IsNull() && !m.DefaultSolution.IsUnknown() {
		output["DefaultSolution"] = m.DefaultSolution
	}

	// Handle AdditionalComponents map
	if !m.AdditionalComponents.IsNull() && !m.AdditionalComponents.IsUnknown() {
		componentsMap := make(map[string]interface{})
		for key, value := range m.AdditionalComponents.Elements() {
			// We expect these to be ComponentValueObject
			if compObj, ok := value.(ComponentValueObject); ok && !compObj.IsNull() && !compObj.IsUnknown() {
				componentsMap[key] = compObj
			} else if strVal, ok := value.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
				// Handle string values (usually from flattened components)
				componentsMap[key] = strVal.ValueString()
			}
		}
		output["AdditionalComponents"] = componentsMap
	}

	// Handle AdditionalNamingPatterns map
	if !m.AdditionalNamingPatterns.IsNull() && !m.AdditionalNamingPatterns.IsUnknown() {
		patternsMap := make(map[string]interface{})
		for key, value := range m.AdditionalNamingPatterns.Elements() {
			// We expect these to be types.String
			if strVal, ok := value.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
				patternsMap[key] = strVal.ValueString()
			}
		}
		output["AdditionalNamingPatterns"] = patternsMap
	}

	return json.Marshal(output)
}

// Configure prepares the provider for data sources and resources
func (p *resourcenamingtoolFunctionsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	logInfo(ctx, "Configuring resourcenamingtool provider...")

	// Load configuration from the file that was saved during ValidateConfig
	// This simplifies our model - ValidateConfig is the source of truth for configuration
	logDebug(ctx, "Loading configuration from file that was saved during ValidateConfig")
	config := loadProviderConfigFromFile(ctx)

	if config == nil {
		// If no configuration found in file, this is unexpected since ValidateConfig should have created it
		logError(ctx, "No configuration file found, this is unexpected as ValidateConfig should have created it")
		resp.Diagnostics.AddError(
			"Configuration Load Error",
			"Unable to load provider configuration from file. ValidateConfig should have created this file.",
		)
		return
	}

	logDebug(ctx, "Successfully loaded configuration from file")

	// Store the configuration in the provider struct
	p.config = config

	logDebug(ctx, "Provider configuration complete")
}

// Resources returns the resources to register for this provider
func (p *resourcenamingtoolFunctionsProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

// DataSources returns the data sources to register for this provider
func (p *resourcenamingtoolFunctionsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProviderStatusDataSource,
	}
}

// Functions returns the functions to register for this provider
func (p *resourcenamingtoolFunctionsProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		func() function.Function {
			return NewGenerateResourceNameFunction(p.config)
		},
	}
}

// ValidateConfig validates the provider configuration values and sets up defaults
// This is called before Configure to allow validating configuration values
// This is our primary location for processing and caching configuration
func (p *resourcenamingtoolFunctionsProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	logDebug(ctx, "Validating resourcenamingtool provider configuration...")

	// Retrieve provider data from configuration
	var config resourcenamingtoolProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate additional components if provided
	logDebug(ctx, "Validating additional components...")
	// Check if additional components are provided and validate them
	if !config.AdditionalComponents.IsNull() && !config.AdditionalComponents.IsUnknown() {
		logDebug(ctx, "Validating additional components: %s", config.AdditionalComponents.String())
		for key, value := range config.AdditionalComponents.Elements() {
			// Validate that keys in additional components follow the expected pattern format
			if !strings.HasPrefix(key, "{") || !strings.HasSuffix(key, "}") {
				resp.Diagnostics.AddAttributeError(
					path.Root("additional_components"),
					"Invalid Component Key Format",
					fmt.Sprintf("Component key %q must be wrapped in curly braces, e.g. {component_name}", key),
				)
				logDebug(ctx, "Invalid component key format: %s", key)
				continue
			}

			// Validate that values are not null or empty strings
			if value.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("additional_components"),
					"Null Component Value",
					fmt.Sprintf("Component value for key %q cannot be null", key),
				)
				logDebug(ctx, "Null component value for key: %s", key)
			} else {
				logDebug(ctx, "Valid component value for key: %s", key)
			}
		}
	} else {
		logDebug(ctx, "No additional components provided or they are unknown")
	}

	// Validate additional naming patterns if provided
	logDebug(ctx, "Validating additional naming patterns...")
	// Check if additional naming patterns are provided and validate them
	if !config.AdditionalNamingPatterns.IsNull() && !config.AdditionalNamingPatterns.IsUnknown() {
		logDebug(ctx, "Validating additional naming patterns: %s", config.AdditionalNamingPatterns.String())
		for key, value := range config.AdditionalNamingPatterns.Elements() {
			// Check that pattern values are not null
			if value.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("additional_naming_patterns"),
					"Null Pattern Value",
					fmt.Sprintf("Naming pattern for resource type %q cannot be null", key),
				)
				logDebug(ctx, "Null pattern value for resource type: %s", key)
				continue
			}

			// Check that pattern values contain at least one component placeholder
			patternStr, ok := value.(types.String)
			if ok && !strings.Contains(patternStr.ValueString(), "{") {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("additional_naming_patterns"),
					"Invalid Naming Pattern",
					fmt.Sprintf("Naming pattern for resource type %q does not contain any component placeholders", key),
				)
				logDebug(ctx, "Invalid naming pattern for resource type: %s", key)
			} else {
				logDebug(ctx, "Valid naming pattern for resource type: %s", key)
			}
		}
	} else {
		logDebug(ctx, "No additional naming patterns provided or they are unknown")
	}

	// If component values are provided, validate them
	validateComponentIfProvided := func(comp ComponentValueObject, attrName string) {
		logDebug(ctx, "Validating component: %s", attrName)
		if comp.IsNull() || comp.IsUnknown() {
			logDebug(ctx, "Component is null or unknown: %s", attrName)
			return
		} else {
			logDebug(ctx, "Component is not null or unknown: %s", attrName)
		}

		// At least one of fullname, shortcode, or char should be provided
		fullname, diagFull := comp.GetFullname(ctx)
		shortcode, diagShort := comp.GetShortcode(ctx)
		char, diagChar := comp.GetChar(ctx)

		if (diagFull.HasError() || fullname == "") &&
			(diagShort.HasError() || shortcode == "") &&
			(diagChar.HasError() || char == "") {
			resp.Diagnostics.AddAttributeError(
				path.Root(attrName),
				"Invalid Component Configuration",
				fmt.Sprintf("At least one of fullname, shortcode, or char must be provided for %s", attrName),
			)
			logDebug(ctx, "Invalid component configuration for %s", attrName)
		} else {
			// Print valid component configuration's fullname, shortcode, and char
			logDebug(ctx, "Valid component configuration for %s: fullname: %s shortcode: %s char: %s",
				attrName, fullname, shortcode, char)
		}
	}

	// Validate core components if provided
	validateComponentIfProvided(config.DefaultResourceType, "default_resource_type")
	validateComponentIfProvided(config.DefaultResourcePrefix, "default_resource_prefix")
	validateComponentIfProvided(config.DefaultBasename, "default_basename")
	validateComponentIfProvided(config.DefaultEnvironment, "default_environment")
	validateComponentIfProvided(config.DefaultRegion, "default_region")
	validateComponentIfProvided(config.DefaultInstance, "default_instance")

	// Validate organization components if provided
	validateComponentIfProvided(config.DefaultOrganization, "default_organization")
	validateComponentIfProvided(config.DefaultBusinessUnit, "default_business_unit")
	validateComponentIfProvided(config.DefaultCostCenter, "default_cost_center")
	validateComponentIfProvided(config.DefaultProject, "default_project")
	validateComponentIfProvided(config.DefaultApplication, "default_application")
	validateComponentIfProvided(config.DefaultWorkload, "default_workload")

	// Validate provider specific components if provided
	validateComponentIfProvided(config.DefaultSubscription, "default_subscription")
	validateComponentIfProvided(config.DefaultLocation, "default_location")
	validateComponentIfProvided(config.DefaultDomain, "default_domain")
	validateComponentIfProvided(config.DefaultCriticality, "default_criticality")

	// Validate initiative/solution components if provided
	validateComponentIfProvided(config.DefaultInitiative, "default_initiative")
	validateComponentIfProvided(config.DefaultSolution, "default_solution")

	if resp.Diagnostics.HasError() {
		logError(ctx, "Validation errors detected, not saving configuration")
		return
	}

	// Store the validated configuration in the provider struct
	p.config = &config

	// Save to file for cross-process sharing - this is the primary way functions will access the config
	logDebug(ctx, "Saving configuration to file for cross-process sharing")
	if err := saveProviderConfigToFile(ctx, &config); err != nil {
		logErrorWithFields(ctx, "Failed to save configuration to file", map[string]interface{}{
			"error": err.Error(),
		})
		resp.Diagnostics.AddError(
			"Configuration Persistence Error",
			fmt.Sprintf("Failed to save configuration to file: %s", err.Error()),
		)
	} else {
		logDebug(ctx, "Successfully saved configuration to file for cross-process sharing")
	}
}

// ResourceNamingToolProvider is an alias for resourcenamingtoolFunctionsProvider
// for backward compatibility with older code
type ResourceNamingToolProvider = resourcenamingtoolFunctionsProvider
