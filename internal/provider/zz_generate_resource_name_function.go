// Copyright (c) Thomas Geens

package provider

import (
	"context"
	_ "embed" // Import the embed package
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

//go:embed descriptions/generate_resource_name_description.txt
var generateResourceNameDescription string

//go:embed descriptions/generate_resource_name_markdown_description.md
var generateResourceNameMarkdownDescription string

// GenerateResourceNameFunction implements function.Function with provider access
type GenerateResourceNameFunction struct {
	// Store the provider configuration pointer itself
	config *resourcenamingtoolProviderModel
}

// NewGenerateResourceNameFunction creates a new instance with the provider config
func NewGenerateResourceNameFunction(config *resourcenamingtoolProviderModel) function.Function {
	return &GenerateResourceNameFunction{
		config: config, // Store the pointer directly, don't dereference
	}
}

func (f *GenerateResourceNameFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "generate_resource_name"
}

func (f *GenerateResourceNameFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	logDebug(ctx, "Defining generate_resource_name function")

	resp.Definition = function.Definition{
		Summary:             "Generate a resource name based on resource name patterns and parameters.",
		Description:         generateResourceNameDescription,
		MarkdownDescription: generateResourceNameMarkdownDescription,
		Parameters: []function.Parameter{
			function.SetParameter{
				Name:        "parameters",
				Description: "A set of parameters used to generate the resource name.",
				ElementType: types.MapType{
					ElemType: types.MapType{
						ElemType: types.StringType,
					},
				},
				AllowNullValue:     true,
				AllowUnknownValues: true,
			},
		},
		Return: function.StringReturn{},
	}

	logDebugWithFields(ctx, "Function definition completed", map[string]interface{}{
		"parameters_count": len(resp.Definition.Parameters),
		"parameter_type":   fmt.Sprintf("%T", resp.Definition.Parameters[0]),
	})
}

// ComponentParameterType is a custom type for component parameters
type ComponentParameterType struct {
	basetypes.ObjectType
}

// NewComponentParameterType creates a new ComponentParameterType
func NewComponentParameterType() ComponentParameterType {
	return ComponentParameterType{}
}

// Equal returns true if the given type is equivalent
func (t ComponentParameterType) Equal(o attr.Type) bool {
	other, ok := o.(ComponentParameterType)
	if !ok {
		return false
	}
	return t.ObjectType.Equal(other.ObjectType)
}

// String returns a human-readable representation of the type
func (t ComponentParameterType) String() string {
	return "ComponentParameterType"
}

// TerraformType returns the tftypes.Type that should be used to represent this type
func (t ComponentParameterType) TerraformType(ctx context.Context) tftypes.Type {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":      tftypes.String,
			"fullname":  tftypes.String,
			"shortcode": tftypes.String,
			"char":      tftypes.String,
		},
		OptionalAttributes: map[string]struct{}{
			"fullname":  {},
			"shortcode": {},
			"char":      {},
		},
	}
}

// ValueFromTerraform transforms a tftypes.Value into the appropriate Go type
func (t ComponentParameterType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	// Use basic logging to help with debugging
	logDebugWithFields(ctx, "ComponentParameterType.ValueFromTerraform called", map[string]interface{}{
		"value_type":   in.Type().String(),
		"is_null":      in.IsNull(),
		"is_known":     in.IsKnown(),
		"debug_string": fmt.Sprintf("%#v", in),
	})

	// If the value is null or unknown, create a null or unknown value
	if !in.IsKnown() {
		logDebug(ctx, "Value is unknown, returning unknown object")
		return basetypes.NewObjectUnknown(map[string]attr.Type{
			"name":      types.StringType,
			"fullname":  types.StringType,
			"shortcode": types.StringType,
			"char":      types.StringType,
		}), nil
	}

	if in.IsNull() {
		logDebug(ctx, "Value is null, returning null object")
		return basetypes.NewObjectNull(map[string]attr.Type{
			"name":      types.StringType,
			"fullname":  types.StringType,
			"shortcode": types.StringType,
			"char":      types.StringType,
		}), nil
	}

	// Extract the values from the tftypes.Value
	var objMap map[string]tftypes.Value
	err := in.As(&objMap)
	if err != nil {
		logErrorWithFields(ctx, "Error extracting values from tftypes.Value", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	// Create the attribute types and values
	attrTypes := map[string]attr.Type{
		"name":      types.StringType,
		"fullname":  types.StringType,
		"shortcode": types.StringType,
		"char":      types.StringType,
	}

	attrs := map[string]attr.Value{
		"name":      types.StringNull(),
		"fullname":  types.StringNull(),
		"shortcode": types.StringNull(),
		"char":      types.StringNull(),
	}

	// Extract the name value
	if nameVal, ok := objMap["name"]; ok && !nameVal.IsNull() && nameVal.IsKnown() {
		var nameStr string
		if err := nameVal.As(&nameStr); err == nil {
			attrs["name"] = types.StringValue(nameStr)
		}
	}

	// Extract the fullname value
	if fullnameVal, ok := objMap["fullname"]; ok && !fullnameVal.IsNull() && fullnameVal.IsKnown() {
		var fullnameStr string
		if err := fullnameVal.As(&fullnameStr); err == nil {
			attrs["fullname"] = types.StringValue(fullnameStr)
		}
	}

	// Extract the shortcode value
	if shortcodeVal, ok := objMap["shortcode"]; ok && !shortcodeVal.IsNull() && shortcodeVal.IsKnown() {
		var shortcodeStr string
		if err := shortcodeVal.As(&shortcodeStr); err == nil {
			attrs["shortcode"] = types.StringValue(shortcodeStr)
		}
	}

	// Extract the char value
	if charVal, ok := objMap["char"]; ok && !charVal.IsNull() && charVal.IsKnown() {
		var charStr string
		if err := charVal.As(&charStr); err == nil {
			attrs["char"] = types.StringValue(charStr)
		}
	}

	// Create the object value
	return types.ObjectValueMust(attrTypes, attrs), nil
}

// ValueType returns the value type for this type
func (t ComponentParameterType) ValueType(ctx context.Context) attr.Value {
	return ComponentParameterValue{}
}

// ComponentParameterValue is a custom value type for component parameters
type ComponentParameterValue struct {
	basetypes.ObjectValue
}

// Type returns the type of the value
func (v ComponentParameterValue) Type(ctx context.Context) attr.Type {
	return ComponentParameterType{}
}

func (f *GenerateResourceNameFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	logDebug(ctx, "Invoking GenerateResourceNameFunction...")

	// Parse the incoming parameters as a set
	var parametersSet types.Set
	diags := req.Arguments.Get(ctx, &parametersSet)

	if diags != nil {
		logError(ctx, "Failed to parse parameters set: %s", diags.Error())
		resp.Error = function.NewFuncError("Failed to parse parameters set: " + diags.Error())
		return
	}

	logDebugWithFields(ctx, "Successfully parsed parameters set", map[string]interface{}{
		"is_null":    parametersSet.IsNull(),
		"is_unknown": parametersSet.IsUnknown(),
		"count":      len(parametersSet.Elements()),
		"elements":   fmt.Sprintf("%#v", parametersSet.Elements()),
	})

	// Convert the set to a ResourceNamingParametersValue for use with the existing generateResourceName function
	resourceParams, err := setWithNestedMapsToResourceNamingParametersValue(ctx, parametersSet)
	if err != nil {
		logError(ctx, "Failed to convert parameters: %s", err.Error())
		resp.Error = function.NewFuncError("Failed to convert parameters: " + err.Error())
		return
	} else {
		logDebugWithFields(ctx, "Successfully converted parameters to ResourceNamingParametersValue", map[string]interface{}{
			"parameters": resourceParams,
			"length":     len(resourceParams.Attributes()),
			"elements":   fmt.Sprintf("%#v", resourceParams.Attributes()),
		})
	}

	// Get configuration - use the shared provider config, potentially loading from file
	var config *resourcenamingtoolProviderModel

	// Try to get the shared provider configuration, which will now check both the
	// in-memory atomic variable and the file-based storage
	sharedConfig := GetSharedProviderConfig(ctx)
	// Show sharedConfig in debug
	logDebugWithFields(ctx, "Shared provider configuration", map[string]interface{}{
		"config": sharedConfig,
	})
	if sharedConfig != nil {
		// Create a deep copy to avoid modifying the shared config
		config = &resourcenamingtoolProviderModel{}
		*config = *sharedConfig
		logDebug(ctx, "Using shared provider configuration")
	} else {
		// No shared config available, create a safe empty config to avoid nil pointer dereference
		logDebug(ctx, "No shared provider configuration found, creating empty config")
		config = &resourcenamingtoolProviderModel{}
	}

	// If function has a local config, use it to override specific values
	if f.config != nil {
		logDebug(ctx, "Found function-specific configuration")

		// Only override values that are not null in the local config
		if !f.config.DefaultResourceType.IsNull() {
			config.DefaultResourceType = f.config.DefaultResourceType
		}
		if !f.config.DefaultResourcePrefix.IsNull() {
			config.DefaultResourcePrefix = f.config.DefaultResourcePrefix
		}
		if !f.config.DefaultBasename.IsNull() {
			config.DefaultBasename = f.config.DefaultBasename
		}
		if !f.config.DefaultEnvironment.IsNull() {
			config.DefaultEnvironment = f.config.DefaultEnvironment
		}
		if !f.config.DefaultRegion.IsNull() {
			config.DefaultRegion = f.config.DefaultRegion
		}
		if !f.config.DefaultInstance.IsNull() {
			config.DefaultInstance = f.config.DefaultInstance
		}
		if !f.config.DefaultOrganization.IsNull() {
			config.DefaultOrganization = f.config.DefaultOrganization
		}
		if !f.config.DefaultProject.IsNull() {
			config.DefaultProject = f.config.DefaultProject
		}
		if !f.config.DefaultBusinessUnit.IsNull() {
			config.DefaultBusinessUnit = f.config.DefaultBusinessUnit
		}
		if !f.config.DefaultCostCenter.IsNull() {
			config.DefaultCostCenter = f.config.DefaultCostCenter
		}
		if !f.config.DefaultApplication.IsNull() {
			config.DefaultApplication = f.config.DefaultApplication
		}
		if !f.config.DefaultWorkload.IsNull() {
			config.DefaultWorkload = f.config.DefaultWorkload
		}
		if !f.config.DefaultSubscription.IsNull() {
			config.DefaultSubscription = f.config.DefaultSubscription
		}
		if !f.config.DefaultLocation.IsNull() {
			config.DefaultLocation = f.config.DefaultLocation
		}
		if !f.config.DefaultDomain.IsNull() {
			config.DefaultDomain = f.config.DefaultDomain
		}
		if !f.config.DefaultCriticality.IsNull() {
			config.DefaultCriticality = f.config.DefaultCriticality
		}
		if !f.config.DefaultInitiative.IsNull() {
			config.DefaultInitiative = f.config.DefaultInitiative
		}
		if !f.config.DefaultSolution.IsNull() {
			config.DefaultSolution = f.config.DefaultSolution
		}
	}

	// Variables to hold diagnostics and result
	var result string
	var resultDiags diag.Diagnostics

	// Try to get resource_type from parameters
	resourceTypeComp, diagResType := resourceParams.GetComponentValue(ctx, "resource_type")

	// If we don't have a valid resource_type in parameters
	if diagResType.HasError() || resourceTypeComp.IsNull() {
		// No resource_type provided, we need to add the default one from the provider
		componentAttrs := resourceParams.Attributes()
		if componentAttrs == nil {
			componentAttrs = make(map[string]attr.Value)
		}

		// Create a default resource_type component if none was provided
		componentAttrs["resource_type"] = config.DefaultResourceType

		// Create a new params object with the default resource_type
		attrTypes := make(map[string]attr.Type)
		for k, v := range componentAttrs {
			attrTypes[k] = v.Type(ctx)
		}

		newParams, _ := types.ObjectValue(attrTypes, componentAttrs)
		updatedParams := ResourceNamingParametersValue{
			ObjectValue: newParams,
		}

		// Generate the resource name using the updated parameters
		result, resultDiags = generateResourceName(ctx, updatedParams, *config)
	} else {
		// We have a resource_type, use the parameters as provided
		result, resultDiags = generateResourceName(ctx, resourceParams, *config)
	}

	// Check if there are any error diagnostics
	if resultDiags.HasError() {
		// Collect all error messages into a single error message
		var errorMessages strings.Builder

		for _, d := range resultDiags {
			if d.Severity() == diag.SeverityError {
				if errorMessages.Len() > 0 {
					errorMessages.WriteString("; ")
				}
				errorMessages.WriteString(d.Summary())
				if d.Detail() != "" {
					errorMessages.WriteString(": ")
					errorMessages.WriteString(d.Detail())
				}

				logErrorWithFields(ctx, "Error generating resource name", map[string]interface{}{
					"summary": d.Summary(),
					"detail":  d.Detail(),
				})
			}
		}

		if errorMessages.Len() > 0 {
			// Create a function error that will cause Terraform to fail
			resp.Error = function.NewFuncError(errorMessages.String())
			return
		}
	}

	// Log the result
	logDebugWithFields(ctx, "Generated resource name", map[string]interface{}{
		"result": result,
	})

	// Set the result
	resp.Error = resp.Result.Set(ctx, result)
}

// setWithNestedMapsToResourceNamingParametersValue converts a set of nested maps to a ResourceNamingParametersValue
func setWithNestedMapsToResourceNamingParametersValue(ctx context.Context, parametersSet types.Set) (ResourceNamingParametersValue, error) {
	// Create a new ResourceNamingParametersValue
	attrTypes := map[string]attr.Type{}
	attributes := map[string]attr.Value{}

	// For collecting flattened additional_components entries
	additionalComponentsMap := make(map[string]map[string]string)
	hasAdditionalComponents := false

	// Process each element in the set
	for _, element := range parametersSet.Elements() {
		// Each element should be a map where keys are component names and values are maps of attributes
		if elemMap, ok := element.(types.Map); ok {
			logDebugWithFields(ctx, "Processing map element", map[string]interface{}{
				"element_type": fmt.Sprintf("%T", element),
				"keys_count":   len(elemMap.Elements()),
			})

			// Each map in the set represents a component with its attributes
			for componentName, componentValue := range elemMap.Elements() {
				logDebugWithFields(ctx, "Processing component", map[string]interface{}{
					"component_name":  componentName,
					"component_type":  fmt.Sprintf("%T", componentValue),
					"component_value": fmt.Sprintf("%v", componentValue),
				})

				// Special handling for additional_components
				if componentName == "additional_components" {
					logDebugWithFields(ctx, "Found additional_components map", map[string]interface{}{
						"component_value": fmt.Sprintf("%v", componentValue),
					})

					// The component value should be a map with dotted notation keys
					if componentsMap, ok := componentValue.(types.Map); ok {
						logDebugWithFields(ctx, "additional_components is a map", map[string]interface{}{
							"components_map_count": len(componentsMap.Elements()),
							"components_map_keys":  fmt.Sprintf("%v", componentsMap.Elements()),
						})

						// Process each entry in the additional_components map
						for keyDotted, valueString := range componentsMap.Elements() {
							// Parse the dotted key to extract component name and attribute (format: "componentName.attribute")
							parts := strings.Split(keyDotted, ".")
							if len(parts) == 2 {
								compName := parts[0]
								attrName := parts[1]

								// Get or create map for this component
								if _, exists := additionalComponentsMap[compName]; !exists {
									additionalComponentsMap[compName] = make(map[string]string)
								}

								// Add value to the component's map
								if strVal, ok := valueString.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
									additionalComponentsMap[compName][attrName] = strVal.ValueString()
									logDebugWithFields(ctx, "Added additional component attribute from map", map[string]interface{}{
										"component": compName,
										"attribute": attrName,
										"value":     strVal.ValueString(),
									})
								}

								hasAdditionalComponents = true
							} else {
								logErrorWithFields(ctx, "Invalid key format in additional_components", map[string]interface{}{
									"key":           keyDotted,
									"expected_form": "componentName.attribute",
								})
							}
						}
						continue
					}
				}

				// Special handling for additional_naming_patterns
				if componentName == "additional_naming_patterns" {
					logDebugWithFields(ctx, "Found additional_naming_patterns map", map[string]interface{}{
						"component_value": fmt.Sprintf("%v", componentValue),
					})

					// The component value should be a map with string values
					if valueMap, ok := componentValue.(types.Map); ok {
						logDebugWithFields(ctx, "additional_naming_patterns is a map", map[string]interface{}{
							"pattern_map_count": len(valueMap.Elements()),
							"pattern_map_keys":  fmt.Sprintf("%v", valueMap.Elements()),
						})

						// Convert elements to ensure they're all strings
						stringElements := make(map[string]attr.Value)
						for patternKey, patternValue := range valueMap.Elements() {
							if strVal, ok := patternValue.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
								stringElements[patternKey] = strVal
								logDebugWithFields(ctx, "Added naming pattern", map[string]interface{}{
									"key":   patternKey,
									"value": strVal.ValueString(),
								})
							} else {
								// Try to convert to string if not already a string
								patternStr := patternValue.String()
								stringElements[patternKey] = types.StringValue(patternStr)
								logDebugWithFields(ctx, "Converted naming pattern value to string", map[string]interface{}{
									"key":             patternKey,
									"original_type":   fmt.Sprintf("%T", patternValue),
									"converted_value": patternStr,
								})
							}
						}

						// Create the map value for additional_naming_patterns
						patternMap, diags := types.MapValue(types.StringType, stringElements)
						if diags.HasError() {
							logErrorWithFields(ctx, "Failed to create map of strings for additional_naming_patterns", map[string]interface{}{
								"error": diags.Errors()[0].Summary(),
							})
							return ResourceNamingParametersValue{}, fmt.Errorf("failed to create additional_naming_patterns map: %s", diags.Errors()[0].Summary())
						}

						// Add to attributes
						attrTypes["additional_naming_patterns"] = types.MapType{ElemType: types.StringType}
						attributes["additional_naming_patterns"] = patternMap
						logDebugWithFields(ctx, "Added additional_naming_patterns to ResourceNamingParametersValue", map[string]interface{}{
							"pattern_count": len(stringElements),
						})
						continue
					}
				}

				// The component value should be a map with attribute values
				if valueMap, ok := componentValue.(types.Map); ok {
					logDebugWithFields(ctx, "Component is a map", map[string]interface{}{
						"component_name":      componentName,
						"component_map_count": len(valueMap.Elements()),
						"component_map_keys":  fmt.Sprintf("%v", valueMap.Elements()),
					})

					// Extract attribute values from the map
					attrs := make(map[string]attr.Value)
					attrType := map[string]attr.Type{
						"name":      types.StringType,
						"fullname":  types.StringType,
						"shortcode": types.StringType,
						"char":      types.StringType,
					}

					// Initialize with null values
					attrs["name"] = types.StringNull()
					attrs["fullname"] = types.StringNull()
					attrs["shortcode"] = types.StringNull()
					attrs["char"] = types.StringNull()

					// Extract values from the map
					for attrKey, attrValue := range valueMap.Elements() {
						if strVal, ok := attrValue.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
							logDebugWithFields(ctx, "Found attribute value", map[string]interface{}{
								"attribute": attrKey,
								"value":     strVal.ValueString(),
							})

							switch attrKey {
							case "name":
								attrs["name"] = types.StringValue(strVal.ValueString())
							case "fullname":
								attrs["fullname"] = types.StringValue(strVal.ValueString())
							case "shortcode":
								attrs["shortcode"] = types.StringValue(strVal.ValueString())
							case "char":
								attrs["char"] = types.StringValue(strVal.ValueString())
							}
						}
					}

					// Create an object value for this component
					objValue, diags := types.ObjectValue(attrType, attrs)
					if diags.HasError() {
						logErrorWithFields(ctx, "Failed to create object value for component", map[string]interface{}{
							"component": componentName,
							"error":     diags.Errors()[0].Summary(),
						})
						return ResourceNamingParametersValue{}, fmt.Errorf("failed to create component %s object: %s", componentName, diags.Errors()[0].Summary())
					}

					// Add the component to the attributes map
					attrTypes[componentName] = types.ObjectType{
						AttrTypes: attrType,
					}
					attributes[componentName] = objValue

					logDebugWithFields(ctx, "Added component to ResourceNamingParametersValue", map[string]interface{}{
						"component":     componentName,
						"attributes":    fmt.Sprintf("%#v", attrs),
						"element_type":  fmt.Sprintf("%T", element),
						"element_value": fmt.Sprintf("%#v", componentValue),
					})
				} else {
					// Component value is not a map - this is unexpected
					logErrorWithFields(ctx, "Component value is not a map", map[string]interface{}{
						"component":      componentName,
						"component_type": fmt.Sprintf("%T", componentValue),
					})
					return ResourceNamingParametersValue{}, fmt.Errorf("component %s value is not a map: %T", componentName, componentValue)
				}
			}
		} else {
			// Element is not a map - this is unexpected
			logErrorWithFields(ctx, "Set element is not a map", map[string]interface{}{
				"element_type": fmt.Sprintf("%T", element),
			})
			return ResourceNamingParametersValue{}, fmt.Errorf("set element is not a map: %T", element)
		}
	}

	// Process collected additional_components after all elements are processed
	if hasAdditionalComponents {
		logDebugWithFields(ctx, "Processing flattened additional components", map[string]interface{}{
			"component_count": len(additionalComponentsMap),
			"components":      fmt.Sprintf("%v", additionalComponentsMap),
		})

		// Create a flattened map for additional_components
		flattenedComponents := make(map[string]attr.Value)

		// Add each component attribute as a flattened string
		for compName, attrs := range additionalComponentsMap {
			for attrName, attrValue := range attrs {
				flattenedKey := fmt.Sprintf("%s.%s", compName, attrName)
				flattenedComponents[flattenedKey] = types.StringValue(attrValue)
				logDebugWithFields(ctx, "Added flattened component attribute", map[string]interface{}{
					"key":   flattenedKey,
					"value": attrValue,
				})
			}
		}

		// Create the map value for additional_components as a map of strings
		componentsMapVal, mapDiags := types.MapValue(types.StringType, flattenedComponents)
		if mapDiags.HasError() {
			logErrorWithFields(ctx, "Failed to create map of strings for additional_components", map[string]interface{}{
				"error": mapDiags.Errors()[0].Summary(),
			})
		} else {
			// Add to attributes with the correct type - map of strings
			attrTypes["additional_components"] = types.MapType{ElemType: types.StringType}
			attributes["additional_components"] = componentsMapVal
			logDebugWithFields(ctx, "Added additional_components to ResourceNamingParametersValue", map[string]interface{}{
				"component_count": len(flattenedComponents),
			})
		}
	}

	// Create the final object value
	objVal, diags := types.ObjectValue(attrTypes, attributes)
	if diags.HasError() {
		logErrorWithFields(ctx, "Failed to create ResourceNamingParametersValue object", map[string]interface{}{
			"error": diags.Errors()[0].Summary(),
		})
		return ResourceNamingParametersValue{}, fmt.Errorf("failed to create ResourceNamingParametersValue: %s", diags.Errors()[0].Summary())
	}

	logDebugWithFields(ctx, "Successfully created ResourceNamingParametersValue", map[string]interface{}{
		"attribute_count": len(attributes),
	})

	return ResourceNamingParametersValue{
		ObjectValue: objVal,
	}, nil
}

// Generate a resource name by replacing all placeholders in a pattern with actual values
func generateResourceName(ctx context.Context, params ResourceNamingParametersValue, config resourcenamingtoolProviderModel) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	logDebug(ctx, "Starting generateResourceName function")

	// Show the parameters received
	logDebugWithFields(ctx, "Received parameters", map[string]interface{}{
		"length":     len(params.Attributes()),
		"parameters": fmt.Sprintf("%#v", params),
		"elements":   fmt.Sprintf("%#v", params.Attributes()),
	})

	// We no longer need to get a components map since we've flattened the structure
	// Get directly the resource type, which is required
	resourceTypeComp, diagResType := params.GetComponentValue(ctx, "resource_type")
	if diagResType.HasError() {
		logErrorWithFields(ctx, "Error getting resource_type component", map[string]interface{}{
			"error": diagResType.Errors()[0].Summary(),
		})
		diags.Append(diagResType...)
		return "", diags
	}

	resourceTypeFull, diagResTypeFull := resourceTypeComp.GetFullname(ctx)
	if diagResTypeFull.HasError() || resourceTypeFull == "" {
		// Show resourceTypeFull
		logDebug(ctx, "Resource type from parameters is empty or has error, trying default from provider")
		// Use default from provider if not specified
		var localDiag diag.Diagnostics
		resourceTypeFull, localDiag = config.DefaultResourceType.GetFullname(ctx)
		if localDiag.HasError() || resourceTypeFull == "" {
			logErrorWithFields(ctx, "Missing required resource_type parameter", map[string]interface{}{
				"param_is_null":      resourceTypeComp.IsNull(),
				"param_is_unknown":   resourceTypeComp.IsUnknown(),
				"default_is_null":    config.DefaultResourceType.IsNull(),
				"default_is_unknown": config.DefaultResourceType.IsUnknown(),
			})
			diags.AddError("Missing Required Parameter", "Resource type is required")
			return "", diags
		}
		logDebugWithFields(ctx, "Using default resource_type from provider", map[string]interface{}{
			"resource_type": resourceTypeFull,
		})
	} else {
		logDebugWithFields(ctx, "Using resource_type from parameters", map[string]interface{}{
			"resource_type": resourceTypeFull,
		})
	}

	// Build a map of placeholder patterns to actual values
	placeholders := make(map[string]string)
	logDebug(ctx, "Building placeholders map")

	// Create a consolidated map of naming patterns - start with built-in patterns
	patternElements := make(map[string]attr.Value)
	for key, value := range builtin_NamingPatterns {
		patternElements[key] = types.StringValue(value)
		logDebugWithFields(ctx, "Added built-in pattern", map[string]interface{}{
			"key":   key,
			"value": value,
		})
	}

	// Add patterns from provider config
	if !config.AdditionalNamingPatterns.IsNull() && !config.AdditionalNamingPatterns.IsUnknown() {
		for k, v := range config.AdditionalNamingPatterns.Elements() {
			if strVal, ok := v.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
				patternElements[k] = strVal
				logDebugWithFields(ctx, "Added pattern from provider config", map[string]interface{}{
					"key":   k,
					"value": strVal.ValueString(),
				})
			}
		}
	} else {
		logDebug(ctx, "No additional naming patterns found in provider config")
	}

	// Add patterns from function parameters
	if additionalPatterns, patternDiags := params.GetAdditionalNamingPatterns(ctx); !patternDiags.HasError() && !additionalPatterns.IsNull() && !additionalPatterns.IsUnknown() {
		for k, v := range additionalPatterns.Elements() {
			if strVal, ok := v.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
				patternElements[k] = strVal
				logDebugWithFields(ctx, "Added pattern from function parameters", map[string]interface{}{
					"key":   k,
					"value": strVal.ValueString(),
				})
			}
		}

		// Special handling for resource-specific pattern
		if patternVal, exists := additionalPatterns.Elements()[resourceTypeFull]; exists {
			if strVal, ok := patternVal.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
				patternElements[resourceTypeFull] = strVal
				logDebugWithFields(ctx, "Added resource-specific pattern for current resource type", map[string]interface{}{
					"resource_type": resourceTypeFull,
					"pattern":       strVal.ValueString(),
				})
			}
		}
	} else if patternDiags.HasError() {
		logErrorWithFields(ctx, "Error getting additional naming patterns", map[string]interface{}{
			"error": patternDiags.Errors()[0].Summary(),
		})
	}

	// Create the final namingPatterns map
	namingPatterns, diags := types.MapValue(types.StringType, patternElements)
	if diags.HasError() {
		logErrorWithFields(ctx, "Failed to create consolidated naming patterns map", map[string]interface{}{
			"error": diags.Errors()[0].Summary(),
		})
	}

	logDebugWithFields(ctx, "Final naming patterns", map[string]interface{}{
		"naming_patterns": namingPatterns.String(),
		"is_null":         namingPatterns.IsNull(),
		"is_unknown":      namingPatterns.IsUnknown(),
		"length":          len(patternElements),
	})

	// Get the naming pattern for the resource type
	patternValue, ok := namingPatterns.Elements()[resourceTypeFull]
	if !ok {
		logErrorWithFields(ctx, "No naming pattern found for resource type", map[string]interface{}{
			"resource_type": resourceTypeFull,
			"available_patterns": func() []string {
				patterns := make([]string, 0, len(namingPatterns.Elements()))
				for k := range namingPatterns.Elements() {
					patterns = append(patterns, k)
				}
				return patterns
			}(),
		})
		diags.AddError("Missing Pattern", fmt.Sprintf("No naming pattern found for resource type: %s", resourceTypeFull))
		return "", diags
	}

	pattern := ""
	if str, ok := patternValue.(types.String); ok {
		pattern = str.ValueString()
	} else {
		pattern = patternValue.String()
	}

	logDebugWithFields(ctx, "Using naming pattern", map[string]interface{}{
		"resource_type": resourceTypeFull,
		"pattern":       pattern,
	})

	// Define the component types we'll search for in the pattern
	componentTypes := []struct {
		Name      string   // Component name
		Patterns  []string // Placeholders to search for in pattern
		ValueType string   // Which value type to use by default: "full", "short", "char"
	}{
		// Core components
		{"resource_type", []string{"{resource_type}", "{resource_type:full}", "{resource_type:short}", "{resource_type:char}"}, "full"},
		{"resource_prefix", []string{"{resource_prefix}", "{resource_prefix:full}", "{resource_prefix:short}", "{resource_prefix:char}"}, "full"},
		{"basename", []string{"{basename}", "{basename:full}", "{basename:short}", "{basename:char}"}, "full"},
		{"environment", []string{"{environment}", "{environment:full}", "{environment:short}", "{environment:char}", "{env}", "{e}"}, "full"},
		{"region", []string{"{region}", "{region:full}", "{region:short}", "{region:char}", "{location}", "{loc}", "{r}"}, "full"},
		{"instance", []string{"{instance}", "{instance:full}", "{instance:short}", "{instance:char}", "{inst}", "{i}"}, "full"},

		// Organization components
		{"organization", []string{"{organization}", "{organization:full}", "{organization:short}", "{organization:char}", "{org}", "{o}"}, "full"},
		{"business_unit", []string{"{business_unit}", "{business_unit:full}", "{business_unit:short}", "{business_unit:char}", "{bu}"}, "full"},
		{"cost_center", []string{"{cost_center}", "{cost_center:full}", "{cost_center:short}", "{cost_center:char}", "{cc}"}, "full"},
		{"project", []string{"{project}", "{project:full}", "{project:short}", "{project:char}", "{proj}", "{p}"}, "full"},
		{"application", []string{"{application}", "{application:full}", "{application:short}", "{application:char}", "{app}", "{a}"}, "full"},
		{"workload", []string{"{workload}", "{workload:full}", "{workload:short}", "{workload:char}", "{wl}", "{w}"}, "full"},

		// Provider-specific components
		{"subscription", []string{"{subscription}", "{subscription:full}", "{subscription:short}", "{subscription:char}", "{sub}", "{s}"}, "full"},
		{"domain", []string{"{domain}", "{domain:full}", "{domain:short}", "{domain:char}", "{d}"}, "full"},
		{"criticality", []string{"{criticality}", "{criticality:full}", "{criticality:short}", "{criticality:char}", "{crit}", "{c}"}, "full"},

		// Initiative/solution components
		{"initiative", []string{"{initiative}", "{initiative:full}", "{initiative:short}", "{initiative:char}", "{init}"}, "full"},
		{"solution", []string{"{solution}", "{solution:full}", "{solution:short}", "{solution:char}", "{sol}"}, "full"},
	}

	// For each component, check if it's used in the pattern and add the appropriate value
	for _, compType := range componentTypes {
		// For each pattern variation of this component
		for _, placeholder := range compType.Patterns {
			if strings.Contains(pattern, placeholder) {
				logDebugWithFields(ctx, "Found placeholder in pattern", map[string]interface{}{
					"component":   compType.Name,
					"placeholder": placeholder,
					"pattern":     pattern,
				})

				// Determine which format (full, short, char) to use
				format := compType.ValueType // Default

				// Check for explicit format in placeholder
				if strings.Contains(placeholder, ":full") {
					format = "full"
				} else if strings.Contains(placeholder, ":short") {
					format = "short"
				} else if strings.Contains(placeholder, ":char") {
					format = "char"
				} else {
					// For abbreviated patterns, use appropriate format
					if len(placeholder) <= 5 { // Like {env}, {r}, {i}, etc.
						format = "short"
					}
					if len(placeholder) <= 4 { // Like {e}, {r}, {i}, etc.
						format = "char"
					}
				}

				// Get the component value
				compValue, diagComp := params.GetComponentValue(ctx, compType.Name)
				if diagComp.HasError() {
					logDebugWithFields(ctx, "Error getting component value, will use default", map[string]interface{}{
						"component": compType.Name,
						"error":     diagComp.Errors()[0].Summary(),
					})
					// Skip if error - we'll use provider defaults later
					continue
				}

				// Get the appropriate representation based on format
				var value string
				var localDiag diag.Diagnostics

				switch format {
				case "full":
					value, localDiag = compValue.GetFullname(ctx)
				case "short":
					value, localDiag = compValue.GetShortcode(ctx)
					// Fallback to fullname if shortcode is empty
					if localDiag.HasError() || value == "" {
						value, _ = compValue.GetFullname(ctx)
					}
				case "char":
					value, localDiag = compValue.GetChar(ctx)
					// Fallback to first character of fullname if char is empty
					if localDiag.HasError() || value == "" {
						fullValue, _ := compValue.GetFullname(ctx)
						if len(fullValue) > 0 {
							value = string(fullValue[0])
						}
					}
				}

				logDebugWithFields(ctx, "Retrieved component value", map[string]interface{}{
					"component": compType.Name,
					"format":    format,
					"value":     value,
					"has_error": localDiag.HasError(),
				})

				// Use default from provider config if value is empty
				if value == "" {
					logDebugWithFields(ctx, "Component value is empty, trying provider default", map[string]interface{}{
						"component": compType.Name,
					})

					var defaultValue string
					switch compType.Name {
					case "resource_type":
						defaultValue, localDiag = config.DefaultResourceType.GetFullname(ctx)
					case "resource_prefix":
						defaultValue, localDiag = config.DefaultResourcePrefix.GetFullname(ctx)
					case "basename":
						defaultValue, localDiag = config.DefaultBasename.GetFullname(ctx)
					case "environment":
						defaultValue, localDiag = config.DefaultEnvironment.GetFullname(ctx)
					case "region":
						defaultValue, localDiag = config.DefaultRegion.GetFullname(ctx)
					case "instance":
						defaultValue, localDiag = config.DefaultInstance.GetFullname(ctx)
					case "organization":
						defaultValue, localDiag = config.DefaultOrganization.GetFullname(ctx)
					case "business_unit":
						defaultValue, localDiag = config.DefaultBusinessUnit.GetFullname(ctx)
					case "cost_center":
						defaultValue, localDiag = config.DefaultCostCenter.GetFullname(ctx)
					case "project":
						defaultValue, localDiag = config.DefaultProject.GetFullname(ctx)
					case "application":
						defaultValue, localDiag = config.DefaultApplication.GetFullname(ctx)
					case "workload":
						defaultValue, localDiag = config.DefaultWorkload.GetFullname(ctx)
					case "subscription":
						defaultValue, localDiag = config.DefaultSubscription.GetFullname(ctx)
					case "domain":
						defaultValue, localDiag = config.DefaultDomain.GetFullname(ctx)
					case "criticality":
						defaultValue, localDiag = config.DefaultCriticality.GetFullname(ctx)
					case "initiative":
						defaultValue, localDiag = config.DefaultInitiative.GetFullname(ctx)
					case "solution":
						defaultValue, localDiag = config.DefaultSolution.GetFullname(ctx)
					}

					// If there are diagnostics errors, log them but continue
					if localDiag.HasError() {
						logDebugWithFields(ctx, "Error getting default value", map[string]interface{}{
							"component": compType.Name,
							"error":     localDiag.Errors()[0].Summary(),
						})
						continue
					}

					logDebugWithFields(ctx, "Using default value from provider", map[string]interface{}{
						"component": compType.Name,
						"value":     defaultValue,
					})

					// Get the appropriate representation based on format
					switch format {
					case "full":
						value = defaultValue
					case "short":
						// If component is not null, try to get shortcode
						switch compType.Name {
						case "resource_type":
							if !config.DefaultResourceType.IsNull() {
								value, _ = config.DefaultResourceType.GetShortcode(ctx)
							}
						case "resource_prefix":
							if !config.DefaultResourcePrefix.IsNull() {
								value, _ = config.DefaultResourcePrefix.GetShortcode(ctx)
							}
						case "basename":
							if !config.DefaultBasename.IsNull() {
								value, _ = config.DefaultBasename.GetShortcode(ctx)
							}
						case "environment":
							if !config.DefaultEnvironment.IsNull() {
								value, _ = config.DefaultEnvironment.GetShortcode(ctx)
							}
						case "region":
							if !config.DefaultRegion.IsNull() {
								value, _ = config.DefaultRegion.GetShortcode(ctx)
							}
						case "instance":
							if !config.DefaultInstance.IsNull() {
								value, _ = config.DefaultInstance.GetShortcode(ctx)
							}
						case "organization":
							if !config.DefaultOrganization.IsNull() {
								value, _ = config.DefaultOrganization.GetShortcode(ctx)
							}
						case "business_unit":
							if !config.DefaultBusinessUnit.IsNull() {
								value, _ = config.DefaultBusinessUnit.GetShortcode(ctx)
							}
						case "cost_center":
							if !config.DefaultCostCenter.IsNull() {
								value, _ = config.DefaultCostCenter.GetShortcode(ctx)
							}
						case "project":
							if !config.DefaultProject.IsNull() {
								value, _ = config.DefaultProject.GetShortcode(ctx)
							}
						case "application":
							if !config.DefaultApplication.IsNull() {
								value, _ = config.DefaultApplication.GetShortcode(ctx)
							}
						case "workload":
							if !config.DefaultWorkload.IsNull() {
								value, _ = config.DefaultWorkload.GetShortcode(ctx)
							}
						case "subscription":
							if !config.DefaultSubscription.IsNull() {
								value, _ = config.DefaultSubscription.GetShortcode(ctx)
							}
						case "domain":
							if !config.DefaultDomain.IsNull() {
								value, _ = config.DefaultDomain.GetShortcode(ctx)
							}
						case "criticality":
							if !config.DefaultCriticality.IsNull() {
								value, _ = config.DefaultCriticality.GetShortcode(ctx)
							}
						case "initiative":
							if !config.DefaultInitiative.IsNull() {
								value, _ = config.DefaultInitiative.GetShortcode(ctx)
							}
						case "solution":
							if !config.DefaultSolution.IsNull() {
								value, _ = config.DefaultSolution.GetShortcode(ctx)
							}
						}

						// Fallback to first 3 characters of fullname if shortcode is empty
						if value == "" {
							logDebug(ctx, "Shortcode is empty, using maximum first 3 characters of fullname")
							if len(defaultValue) > 3 {
								value = defaultValue[:3]
							} else if len(defaultValue) > 0 {
								// If defaultValue is less than 3 characters, use it as is
								value = defaultValue
							}
						}
					case "char":
						// If component is not null, try to get char
						switch compType.Name {
						case "resource_type":
							if !config.DefaultResourceType.IsNull() {
								value, _ = config.DefaultResourceType.GetChar(ctx)
							}
						case "resource_prefix":
							if !config.DefaultResourcePrefix.IsNull() {
								value, _ = config.DefaultResourcePrefix.GetChar(ctx)
							}
						case "basename":
							if !config.DefaultBasename.IsNull() {
								value, _ = config.DefaultBasename.GetChar(ctx)
							}
						case "environment":
							if !config.DefaultEnvironment.IsNull() {
								value, _ = config.DefaultEnvironment.GetChar(ctx)
							}
						case "region":
							if !config.DefaultRegion.IsNull() {
								value, _ = config.DefaultRegion.GetChar(ctx)
							}
						case "instance":
							if !config.DefaultInstance.IsNull() {
								value, _ = config.DefaultInstance.GetChar(ctx)
							}
						case "organization":
							if !config.DefaultOrganization.IsNull() {
								value, _ = config.DefaultOrganization.GetChar(ctx)
							}
						case "business_unit":
							if !config.DefaultBusinessUnit.IsNull() {
								value, _ = config.DefaultBusinessUnit.GetChar(ctx)
							}
						case "cost_center":
							if !config.DefaultCostCenter.IsNull() {
								value, _ = config.DefaultCostCenter.GetChar(ctx)
							}
						case "project":
							if !config.DefaultProject.IsNull() {
								value, _ = config.DefaultProject.GetChar(ctx)
							}
						case "application":
							if !config.DefaultApplication.IsNull() {
								value, _ = config.DefaultApplication.GetChar(ctx)
							}
						case "workload":
							if !config.DefaultWorkload.IsNull() {
								value, _ = config.DefaultWorkload.GetChar(ctx)
							}
						case "subscription":
							if !config.DefaultSubscription.IsNull() {
								value, _ = config.DefaultSubscription.GetChar(ctx)
							}
						case "domain":
							if !config.DefaultDomain.IsNull() {
								value, _ = config.DefaultDomain.GetChar(ctx)
							}
						case "criticality":
							if !config.DefaultCriticality.IsNull() {
								value, _ = config.DefaultCriticality.GetChar(ctx)
							}
						case "initiative":
							if !config.DefaultInitiative.IsNull() {
								value, _ = config.DefaultInitiative.GetChar(ctx)
							}
						case "solution":
							if !config.DefaultSolution.IsNull() {
								value, _ = config.DefaultSolution.GetChar(ctx)
							}
						}

						// Fallback to first character of fullname if char is empty
						if value == "" && len(defaultValue) > 0 {
							logDebug(ctx, "Char is empty, using first character of fullname")
							value = string(defaultValue[0])
						}
					}
				}

				// Add to placeholders map
				if value != "" {
					placeholders[placeholder] = value
					logDebugWithFields(ctx, "Added placeholder value", map[string]interface{}{
						"placeholder": placeholder,
						"value":       value,
					})
				} else {
					logDebugWithFields(ctx, "Empty value for placeholder", map[string]interface{}{
						"component":   compType.Name,
						"placeholder": placeholder,
					})
				}
			}
		}
	}

	// Process additional components from the additional_components map
	if attrs, ok := params.Attributes()["additional_components"]; ok && !attrs.IsNull() && !attrs.IsUnknown() {
		logDebugWithFields(ctx, "Processing additional_components for placeholders", map[string]interface{}{
			"additional_components": attrs.String(),
		})

		// If it's a map, extract the components
		if additionalMap, ok := attrs.(types.Map); ok {
			// Create a map to group attributes by component name
			componentGroups := make(map[string]map[string]string)

			// Process each flattened entry in the additional_components map
			for key, val := range additionalMap.Elements() {
				// Parse the dotted key to extract component name and attribute
				parts := strings.Split(key, ".")
				if len(parts) == 2 {
					componentName := parts[0]
					attrName := parts[1]

					// Initialize the component group if it doesn't exist
					if _, exists := componentGroups[componentName]; !exists {
						componentGroups[componentName] = make(map[string]string)
					}

					// Add the value to the component group
					if strVal, ok := val.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
						componentGroups[componentName][attrName] = strVal.ValueString()
					}
				}
			}

			// For each component, add placeholders for it
			for componentName, attrs := range componentGroups {
				logDebugWithFields(ctx, "Processing additional component", map[string]interface{}{
					"component":  componentName,
					"attributes": attrs,
				})

				// Add placeholders for this component with different formats
				if fullname, ok := attrs["fullname"]; ok && fullname != "" {
					placeholders["{"+componentName+"}"] = fullname
					placeholders["{"+componentName+":full}"] = fullname
					logDebugWithFields(ctx, "Added additional component placeholder", map[string]interface{}{
						"placeholder": "{" + componentName + "}",
						"value":       fullname,
					})
				}

				if shortcode, ok := attrs["shortcode"]; ok && shortcode != "" {
					placeholders["{"+componentName+":short}"] = shortcode
					logDebugWithFields(ctx, "Added additional component placeholder", map[string]interface{}{
						"placeholder": "{" + componentName + ":short}",
						"value":       shortcode,
					})
				}

				if char, ok := attrs["char"]; ok && char != "" {
					placeholders["{"+componentName+":char}"] = char
					logDebugWithFields(ctx, "Added additional component placeholder", map[string]interface{}{
						"placeholder": "{" + componentName + ":char}",
						"value":       char,
					})
				}
			}
		}
	}

	// Log all naming patterns before generating the result
	logDebugWithFields(ctx, "Final naming patterns for resource name generation", map[string]interface{}{
		"resource_type":         resourceTypeFull,
		"naming_patterns_count": len(patternElements),
		"naming_patterns":       fmt.Sprintf("%v", patternElements),
	})

	// Log all placeholders before generating the result
	logDebugWithFields(ctx, "Final placeholders for pattern substitution", map[string]interface{}{
		"resource_type":      resourceTypeFull,
		"placeholders_count": len(placeholders),
		"placeholders":       fmt.Sprintf("%v", placeholders),
	})

	// Generate the resource name by replacing all placeholders
	result := pattern
	for placeholder, value := range placeholders {
		result = strings.ReplaceAll(result, placeholder, value)
		logDebugWithFields(ctx, "Replaced placeholder", map[string]interface{}{
			"placeholder":   placeholder,
			"value":         value,
			"result_so_far": result,
		})
	}

	logInfoWithFields(ctx, "Generated resource name", map[string]interface{}{
		"resource_type":     resourceTypeFull,
		"pattern_used":      pattern,
		"placeholders_used": fmt.Sprintf("%v", placeholders),
		"result":            result,
	})

	// Verify result and return errors
	if strings.Contains(result, "{") {
		unresolvedComponents := make([]string, 0)
		for _, part := range strings.Split(result, "{") {
			if strings.Contains(part, "}") {
				unresolvedComponents = append(unresolvedComponents, "{"+strings.Split(part, "}")[0]+"}")
			}
		}

		logErrorWithFields(ctx, "Resource name contains unresolved components", map[string]interface{}{
			"result":                result,
			"unresolved_components": unresolvedComponents,
		})
		diags.AddError("Unresolved Components", fmt.Sprintf("Resource name contains unrecognized components: %s", result))
		return "", diags
	}
	if result == "" {
		logError(ctx, "Generated resource name is empty")
		diags.AddError("Empty Name", "Resource name cannot be empty")
		return "", diags
	}
	if len(result) > 90 {
		logErrorWithFields(ctx, "Generated resource name is too long", map[string]interface{}{
			"length": len(result),
			"result": result,
		})
		diags.AddError("Name Too Long", fmt.Sprintf("Resource name exceeds 90 characters: %s", result))
		return "", diags
	}
	if len(result) < 3 {
		logErrorWithFields(ctx, "Generated resource name is too short", map[string]interface{}{
			"length": len(result),
			"result": result,
		})
		diags.AddError("Name Too Short", fmt.Sprintf("Resource name must be at least 3 characters long: %s", result))
		return "", diags
	}

	logDebugWithFields(ctx, "Successfully generated resource name", map[string]interface{}{
		"resource_type": resourceTypeFull,
		"result":        result,
	})
	return result, diags
}
