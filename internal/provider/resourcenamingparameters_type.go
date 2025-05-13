// Copyright (c) Thomas Geens

package provider

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Simple global cache - store the last valid parameters received
// This is a very direct approach but should work for this naming tool use case
var (
	// lastValidParams stores the last set of valid parameters received
	lastValidParams map[string]tftypes.Value
	// globalCacheMutex ensures thread-safety
	globalCacheMutex sync.RWMutex
)

// ResourceNamingParametersType is a custom type for resource naming parameters
type ResourceNamingParametersType struct {
	basetypes.ObjectType

	// Store our custom tftypes.Object with optional attributes
	tfObject tftypes.Object
}

// init initializes the global cache
func init() {
	lastValidParams = make(map[string]tftypes.Value)
}

// NewResourceNamingParametersType creates a new instance of ResourceNamingParametersType
func NewResourceNamingParametersType() *ResourceNamingParametersType {
	componentType := NewComponentValueType()

	// Create attribute types map for our resource naming parameters
	attrTypes := map[string]attr.Type{
		// Core components
		"resource_type":   componentType,
		"resource_prefix": componentType,
		"basename":        componentType,
		"environment":     componentType,
		"region":          componentType,
		"instance":        componentType,

		// Organization related components
		"organization":  componentType,
		"business_unit": componentType,
		"cost_center":   componentType,
		"project":       componentType,
		"application":   componentType,
		"workload":      componentType,

		// Provider specific components
		"subscription": componentType,
		"location":     componentType,
		"domain":       componentType,
		"criticality":  componentType,

		// Initiative/solution related
		"initiative": componentType,
		"solution":   componentType,

		// Extension points - for custom components and patterns
		"additional_components":      types.MapType{ElemType: componentType},
		"additional_naming_patterns": types.MapType{ElemType: types.StringType},
	}

	// Create corresponding tftypes map for our Terraform type representation
	tfAttrTypes := map[string]tftypes.Type{
		// Core components
		"resource_type":   componentType.TerraformType(context.TODO()),
		"resource_prefix": componentType.TerraformType(context.TODO()),
		"basename":        componentType.TerraformType(context.TODO()),
		"environment":     componentType.TerraformType(context.TODO()),
		"region":          componentType.TerraformType(context.TODO()),
		"instance":        componentType.TerraformType(context.TODO()),

		// Organization related components
		"organization":  componentType.TerraformType(context.TODO()),
		"business_unit": componentType.TerraformType(context.TODO()),
		"cost_center":   componentType.TerraformType(context.TODO()),
		"project":       componentType.TerraformType(context.TODO()),
		"application":   componentType.TerraformType(context.TODO()),
		"workload":      componentType.TerraformType(context.TODO()),

		// Provider specific components
		"subscription": componentType.TerraformType(context.TODO()),
		"location":     componentType.TerraformType(context.TODO()),
		"domain":       componentType.TerraformType(context.TODO()),
		"criticality":  componentType.TerraformType(context.TODO()),

		// Initiative/solution related
		"initiative": componentType.TerraformType(context.TODO()),
		"solution":   componentType.TerraformType(context.TODO()),

		// Extension points
		"additional_components":      tftypes.Map{ElementType: componentType.TerraformType(context.TODO())},
		"additional_naming_patterns": tftypes.Map{ElementType: tftypes.String},
	}

	optionalAttrs := map[string]struct{}{}
	for name := range tfAttrTypes {
		optionalAttrs[name] = struct{}{}
	}

	return &ResourceNamingParametersType{
		ObjectType: basetypes.ObjectType{
			AttrTypes: attrTypes,
		},
		tfObject: tftypes.Object{
			AttributeTypes:     tfAttrTypes,
			OptionalAttributes: optionalAttrs,
		},
	}
}

// Equal returns true if the given type is equivalent
func (t ResourceNamingParametersType) Equal(o attr.Type) bool {
	other, ok := o.(*ResourceNamingParametersType)
	if !ok {
		return false
	}
	return t.ObjectType.Equal(other.ObjectType)
}

// String returns a human-readable representation of the type
func (t ResourceNamingParametersType) String() string {
	return "ResourceNamingParametersType"
}

// TerraformType returns the tftypes.Type that should be used to represent this type.
// This is used by the framework for various purposes including to validate if the user-supplied value can be correctly applied.
func (t ResourceNamingParametersType) TerraformType(ctx context.Context) tftypes.Type {
	return tftypes.Object{
		AttributeTypes: t.tfObject.AttributeTypes,
		// Don't use OptionalAttributes here - it's causing the provider to panic
	}
}

// ValueFromTerraform transforms a tftypes.Value into the appropriate Go type
func (t *ResourceNamingParametersType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	// Use basic logging to help with debugging
	logDebugWithFields(ctx, "ResourceNamingParametersType.ValueFromTerraform called", map[string]interface{}{
		"value_type":   in.Type().String(),
		"is_null":      in.IsNull(),
		"is_known":     in.IsKnown(),
		"debug_string": fmt.Sprintf("%#v", in),
	})

	// If the value is null or unknown, create a null or unknown value
	if !in.IsKnown() {
		logDebug(ctx, "Value is unknown, returning unknown object")
		return basetypes.NewObjectUnknown(t.AttrTypes), nil
	}

	if in.IsNull() {
		logDebug(ctx, "Value is null, returning null object")
		return basetypes.NewObjectNull(t.AttrTypes), nil
	}

	// First, try to convert using the base ObjectType (following the HashiCorp pattern)
	attrValue, err := t.ObjectType.ValueFromTerraform(ctx, in)
	if err != nil {
		logError(ctx, "Error converting value from Terraform using ObjectType: %s", err.Error())
		return nil, err
	}

	// Convert to ObjectValue as per HashiCorp's documentation pattern
	objectValue, ok := attrValue.(basetypes.ObjectValue)
	if !ok {
		logError(ctx, "Unexpected value type: expected basetypes.ObjectValue, got %T", attrValue)
		return nil, fmt.Errorf("unexpected value type %T", attrValue)
	}

	logDebug(ctx, "Successfully converted to ObjectValue with %d attributes", len(objectValue.Attributes()))

	// Now extract the raw values for our caching mechanism
	var tfValues map[string]tftypes.Value
	err = in.As(&tfValues)
	if err != nil {
		logError(ctx, "Failed to extract values from tftypes.Value: %s", err.Error())
		return nil, err
	}

	// Check if the object is empty and if we have cached values
	globalCacheMutex.RLock()
	cacheEmpty := len(lastValidParams) == 0
	globalCacheMutex.RUnlock()

	logDebugWithFields(ctx, "Extracted object value", map[string]interface{}{
		"keys":           fmt.Sprintf("%v", getMapKeys(tfValues)),
		"is_empty":       len(tfValues) == 0,
		"cache_is_empty": cacheEmpty,
	})

	// If we received an empty object but have cached parameters, use the cached version
	if len(tfValues) == 0 && !cacheEmpty {
		globalCacheMutex.RLock()
		// Create a deep copy to avoid potential race conditions
		tfValues = make(map[string]tftypes.Value, len(lastValidParams))
		for k, v := range lastValidParams {
			tfValues[k] = v
		}
		globalCacheMutex.RUnlock()

		logDebugWithFields(ctx, "Using cached parameters instead of empty object", map[string]interface{}{
			"cached_params_count": len(tfValues),
			"cached_keys":         fmt.Sprintf("%v", getMapKeys(tfValues)),
		})

		// We need to reconstruct the ObjectValue with our cached parameters
		in = tftypes.NewValue(in.Type(), tfValues)
		attrValue, err = t.ObjectType.ValueFromTerraform(ctx, in)
		if err != nil {
			logError(ctx, "Error converting cached values from Terraform: %s", err.Error())
			return nil, err
		}

		objectValue, ok = attrValue.(basetypes.ObjectValue)
		if !ok {
			logError(ctx, "Unexpected value type after using cache: expected basetypes.ObjectValue, got %T", attrValue)
			return nil, fmt.Errorf("unexpected value type %T after using cache", attrValue)
		}
	} else if len(tfValues) > 0 {
		// Store valid parameters for future use - make a deep copy
		paramsCopy := make(map[string]tftypes.Value, len(tfValues))
		for k, v := range tfValues {
			paramsCopy[k] = v
		}

		globalCacheMutex.Lock()
		lastValidParams = paramsCopy
		globalCacheMutex.Unlock()

		logDebugWithFields(ctx, "Cached new parameters", map[string]interface{}{
			"params_count": len(tfValues),
			"param_keys":   fmt.Sprintf("%v", getMapKeys(tfValues)),
		})
	}

	// Now we can return our custom value type
	logDebug(ctx, "Returning ResourceNamingParametersValue from ObjectValue with %d attributes", len(objectValue.Attributes()))

	return ResourceNamingParametersValue{
		ObjectValue: objectValue,
	}, nil
}

// Helper function to get map keys
func getMapKeys(m map[string]tftypes.Value) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// ValueType returns the value type for this type
func (t ResourceNamingParametersType) ValueType(ctx context.Context) attr.Value {
	return ResourceNamingParametersValue{
		ObjectValue: basetypes.NewObjectNull(map[string]attr.Type{}),
	}
}

// ResourceNamingParametersValue is the value type for ResourceNamingParametersType
type ResourceNamingParametersValue struct {
	basetypes.ObjectValue
}

// Type returns the type of the value
func (v ResourceNamingParametersValue) Type(ctx context.Context) attr.Type {
	return &ResourceNamingParametersType{} // Return a pointer instead of a value
}

// GetAdditionalNamingPatterns retrieves the additional_naming_patterns map as a types.Map
func (v ResourceNamingParametersValue) GetAdditionalNamingPatterns(ctx context.Context) (types.Map, diag.Diagnostics) {
	logDebugWithFields(ctx, "GetAdditionalNamingPatterns called", map[string]interface{}{
		"resourcenamingparametersvalue": fmt.Sprintf("%#v", v),
	})

	var diags diag.Diagnostics

	if v.IsNull() || v.IsUnknown() {
		return types.MapNull(types.StringType), diags
	}

	// Get attributes from the object
	attrs := v.Attributes()

	// Look for additional_naming_patterns in the attributes
	patternsAttr, ok := attrs["additional_naming_patterns"]
	if !ok {
		// No additional patterns defined, return a null map
		logDebug(ctx, "No additional_naming_patterns found, returning null map")
		return types.MapNull(types.StringType), diags
	} else {
		logDebugWithFields(ctx, "Found additional_naming_patterns", map[string]interface{}{
			"patterns": fmt.Sprintf("%#v", patternsAttr),
		})
	}

	// Check if it's already a Map
	if patternsMap, ok := patternsAttr.(types.Map); ok {
		logDebugWithFields(ctx, "Returning existing Map for additional_naming_patterns", map[string]interface{}{
			"patterns": fmt.Sprintf("%#v", patternsMap),
		})
		return patternsMap, diags
	}

	// If it's an object, try to convert it to a map
	if objVal, ok := patternsAttr.(basetypes.ObjectValue); ok {
		// Extract string values from object attributes
		elements := make(map[string]attr.Value)
		for k, v := range objVal.Attributes() {
			if strVal, ok := v.(types.String); ok {
				elements[k] = strVal
			} else {
				diags.AddWarning(
					"Invalid pattern value type",
					fmt.Sprintf("Expected String but got %T for key %s", v, k),
				)
			}
		}

		mapVal, mapDiags := types.MapValue(types.StringType, elements)
		diags.Append(mapDiags...)
		return mapVal, diags
	}

	diags.AddError(
		"Invalid additional_naming_patterns type",
		fmt.Sprintf("Expected Map or Object but got %T", patternsAttr),
	)
	return types.MapNull(types.StringType), diags
}

// GetComponentValue retrieves a component value directly from the object attributes
func (v ResourceNamingParametersValue) GetComponentValue(ctx context.Context, name string) (ComponentValueObject, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v.IsNull() || v.IsUnknown() {
		return ComponentValueObject{}, diags
	}

	// Get attributes from the object
	attrs := v.Attributes()

	// Look for the component directly in the attributes
	componentAttr, ok := attrs[name]
	if !ok {
		// Component not found, check in additional_components if it exists
		additionalComponentsAttr, hasAdditional := attrs["additional_components"]
		if hasAdditional {
			if additionalMap, ok := additionalComponentsAttr.(types.Map); ok && !additionalMap.IsNull() && !additionalMap.IsUnknown() {
				// Get the elements of the additional_components map
				additionalElements := additionalMap.Elements()

				// Check if the component exists in the additional_components map
				if componentInMap, exists := additionalElements[name]; exists {
					if compObj, ok := componentInMap.(ComponentValueObject); ok {
						return compObj, diags
					} else {
						diags.AddError(
							"Invalid component value type in additional_components",
							fmt.Sprintf("Expected ComponentValueObject but got %T", componentInMap),
						)
					}
				}
			}
		}

		// Component not found - return a null component value with no error
		// Since components are optional, this is not an error condition
		nullObj := basetypes.NewObjectNull(map[string]attr.Type{
			"fullname":  types.StringType,
			"shortcode": types.StringType,
			"char":      types.StringType,
		})
		return ComponentValueObject{ObjectValue: nullObj}, diags
	}

	// Convert to ComponentValueObject
	componentValue, ok := componentAttr.(ComponentValueObject)
	if !ok {
		// If it's a regular object, try to convert it
		if objVal, isObj := componentAttr.(basetypes.ObjectValue); isObj {
			return ComponentValueObject{ObjectValue: objVal}, diags
		}

		diags.AddError(
			"Invalid component value type",
			fmt.Sprintf("Expected ComponentValueObject but got %T", componentAttr),
		)
		return ComponentValueObject{}, diags
	}

	return componentValue, diags
}
