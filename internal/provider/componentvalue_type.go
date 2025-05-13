// Copyright (c) Thomas Geens

// Package provider implements the terraform provider resource naming tool functionality.
package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ComponentValueType is the custom type for a component value
type ComponentValueType struct {
	basetypes.ObjectType
}

// NewComponentValueType creates a new ComponentValueType
func NewComponentValueType() ComponentValueType {
	return ComponentValueType{
		ObjectType: basetypes.ObjectType{
			AttrTypes: map[string]attr.Type{
				"fullname":  types.StringType,
				"shortcode": types.StringType,
				"char":      types.StringType,
			},
		},
	}
}

// Equal returns true if the given type is equivalent
func (t ComponentValueType) Equal(o attr.Type) bool {
	other, ok := o.(ComponentValueType)
	if !ok {
		return false
	}
	return t.ObjectType.Equal(other.ObjectType)
}

// String returns a human-readable representation of the type
func (t ComponentValueType) String() string {
	return "ComponentValueType"
}

// TerraformType returns the tftypes.Type that should be used to represent this type
func (t ComponentValueType) TerraformType(ctx context.Context) tftypes.Type {
	return t.ObjectType.TerraformType(ctx)
}

// ValueFromTerraform transforms a tftypes.Value into the appropriate Go type
func (t ComponentValueType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	// Use the base ObjectType's ValueFromTerraform
	attrValue, err := t.ObjectType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	// Convert the basetypes.ObjectValue to our ComponentValueObject
	objValue, ok := attrValue.(basetypes.ObjectValue)
	if !ok {
		return nil, fmt.Errorf("expected basetypes.ObjectValue, got: %T", attrValue)
	}

	// Return our ComponentValueObject wrapping the base ObjectValue
	return ComponentValueObject{
		ObjectValue: objValue,
	}, nil
}

// ValueType returns the value type for this type
func (t ComponentValueType) ValueType(ctx context.Context) attr.Value {
	// Return a ComponentValueObject wrapping a null ObjectValue
	return ComponentValueObject{
		ObjectValue: basetypes.NewObjectNull(
			map[string]attr.Type{
				"fullname":  types.StringType,
				"shortcode": types.StringType,
				"char":      types.StringType,
			},
		),
	}
}

// ComponentValueObject is the value type for ComponentValueType
type ComponentValueObject struct {
	basetypes.ObjectValue
}

// Type returns the type of the value
func (v ComponentValueObject) Type(ctx context.Context) attr.Type {
	return ComponentValueType{
		ObjectType: basetypes.ObjectType{
			AttrTypes: map[string]attr.Type{
				"fullname":  types.StringType,
				"shortcode": types.StringType,
				"char":      types.StringType,
			},
		},
	}
}

// Equal returns true if the values are equivalent
func (v ComponentValueObject) Equal(o attr.Value) bool {
	other, ok := o.(ComponentValueObject)
	if !ok {
		return false
	}
	return v.ObjectValue.Equal(other.ObjectValue)
}

// IsNull returns true if the value is null
func (v ComponentValueObject) IsNull() bool {
	return v.ObjectValue.IsNull()
}

// IsUnknown returns true if the value is unknown
func (v ComponentValueObject) IsUnknown() bool {
	return v.ObjectValue.IsUnknown()
}

// String returns a human-readable representation of the value
func (v ComponentValueObject) String() string {
	return v.ObjectValue.String()
}

// GetFullname returns the fullname value from component
func (v ComponentValueObject) GetFullname(ctx context.Context) (string, diag.Diagnostics) {
	return GetComponentValue(ctx, v.ObjectValue, "fullname")
}

// GetShortcode returns the shortcode value from component
func (v ComponentValueObject) GetShortcode(ctx context.Context) (string, diag.Diagnostics) {
	return GetComponentValue(ctx, v.ObjectValue, "shortcode")
}

// GetChar returns the char value from component
func (v ComponentValueObject) GetChar(ctx context.Context) (string, diag.Diagnostics) {
	return GetComponentValue(ctx, v.ObjectValue, "char")
}

// GetComponentValue extracts values from a basetypes.ObjectValue
func GetComponentValue(ctx context.Context, obj basetypes.ObjectValue, key string) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if obj.IsNull() || obj.IsUnknown() {
		return "", diags
	}

	attrs := obj.Attributes()

	attrVal, ok := attrs[key]
	if !ok {
		logDebug(ctx, "%s attribute not found in component value", key)
		return "", diags
	}

	strVal, ok := attrVal.(types.String)
	if !ok {
		diags.AddError(
			fmt.Sprintf("Invalid %s type", key),
			fmt.Sprintf("The %s attribute is not a string type.", key),
		)
		return "", diags
	}

	if strVal.IsNull() || strVal.IsUnknown() {
		return "", diags
	}

	return strVal.ValueString(), diags
}

// Constructor functions

// NewComponentValueObject creates a new null ComponentValueObject
func NewComponentValueObject() ComponentValueObject {
	return ComponentValueObject{
		ObjectValue: basetypes.NewObjectNull(
			map[string]attr.Type{
				"fullname":  types.StringType,
				"shortcode": types.StringType,
				"char":      types.StringType,
			},
		),
	}
}

// NewComponentValueObjectUnknown creates a new unknown ComponentValueObject
func NewComponentValueObjectUnknown() ComponentValueObject {
	return ComponentValueObject{
		ObjectValue: basetypes.NewObjectUnknown(
			map[string]attr.Type{
				"fullname":  types.StringType,
				"shortcode": types.StringType,
				"char":      types.StringType,
			},
		),
	}
}

// NewComponentValueObjectNull creates a new null ComponentValueObject
func NewComponentValueObjectNull() ComponentValueObject {
	return NewComponentValueObject()
}

// CreateComponentObject creates a base component object
func CreateComponentObject(ctx context.Context, fullname, shortcode, char string) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := map[string]attr.Value{
		"fullname":  types.StringValue(fullname),
		"shortcode": types.StringValue(shortcode),
		"char":      types.StringValue(char),
	}

	objValue, objDiags := types.ObjectValue(
		map[string]attr.Type{
			"fullname":  types.StringType,
			"shortcode": types.StringType,
			"char":      types.StringType,
		},
		attributes,
	)
	diags.Append(objDiags...)

	return objValue, diags
}

// CreateComponentValueObject creates a new ComponentValueObject from a string
func CreateComponentValueObject(ctx context.Context, value string) (ComponentValueObject, diag.Diagnostics) {
	var diags diag.Diagnostics

	baseObject, objDiags := CreateComponentObject(ctx, value, value, value)
	diags.Append(objDiags...)

	return ComponentValueObject{
		ObjectValue: baseObject,
	}, diags
}

// CreateComponentValueObjectFromParts creates a ComponentValueObject from individual parts
func CreateComponentValueObjectFromParts(ctx context.Context, fullname, shortcode, char string) (ComponentValueObject, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Use default values if not provided
	if fullname == "" {
		fullname = ""
	}
	if shortcode == "" {
		shortcode = ""
	}
	if char == "" {
		char = ""
	}

	// Log the values we're using
	logDebug(ctx, "Creating ComponentValueObject with: fullname=%s, shortcode=%s, char=%s",
		fullname, shortcode, char)

	baseObject, objDiags := CreateComponentObject(ctx, fullname, shortcode, char)
	diags.Append(objDiags...)

	if objDiags.HasError() {
		logError(ctx, "Error creating component object: %v", objDiags)
	}

	return ComponentValueObject{
		ObjectValue: baseObject,
	}, diags
}

// MarshalJSON implements custom JSON marshaling for ComponentValueObject
func (v ComponentValueObject) MarshalJSON() ([]byte, error) {
	if v.IsNull() || v.IsUnknown() {
		return []byte("null"), nil
	}

	// Use a context without cancellation for marshaling
	ctx := context.Background()

	// Create a map to store component attributes
	result := make(map[string]interface{})

	// Get component values
	fullname, fullnameDiags := v.GetFullname(ctx)
	if !fullnameDiags.HasError() && fullname != "" {
		result["fullname"] = fullname
	}

	shortcode, shortcodeDiags := v.GetShortcode(ctx)
	if !shortcodeDiags.HasError() && shortcode != "" {
		result["shortcode"] = shortcode
	}

	char, charDiags := v.GetChar(ctx)
	if !charDiags.HasError() && char != "" {
		result["char"] = char
	}

	// If no attributes were added, return an empty object
	if len(result) == 0 {
		return []byte("{}"), nil
	}

	return json.Marshal(result)
}

// UnmarshalJSON implements custom JSON unmarshaling for ComponentValueObject
func (v *ComponentValueObject) UnmarshalJSON(data []byte) error {
	// Handle null case
	if string(data) == "null" {
		*v = NewComponentValueObjectNull()
		return nil
	}

	// Use a context without cancellation for unmarshaling
	ctx := context.Background()

	// Parse JSON data into a map
	var componentData map[string]string
	if err := json.Unmarshal(data, &componentData); err != nil {
		return err
	}

	// Extract values
	fullname := componentData["fullname"]
	shortcode := componentData["shortcode"]
	char := componentData["char"]

	// Create new ComponentValueObject
	compObj, diags := CreateComponentValueObjectFromParts(ctx, fullname, shortcode, char)
	if diags.HasError() {
		return fmt.Errorf("failed to create component value: %v", diags)
	}

	*v = compObj
	return nil
}
