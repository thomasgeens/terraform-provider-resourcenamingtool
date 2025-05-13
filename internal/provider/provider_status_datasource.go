// Copyright (c) Thomas Geens

// provider_status.go - Provider status datasource implementation
package provider

import (
	"context"
	"runtime"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &ProviderStatusDataSource{}

// NewProviderStatusDataSource is a helper function to simplify the provider implementation.
func NewProviderStatusDataSource() datasource.DataSource {
	return &ProviderStatusDataSource{}
}

// ProviderStatusDataSource is the data source implementation.
type ProviderStatusDataSource struct{}

// Metadata returns the data source type name.
func (d *ProviderStatusDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status"
}

// Schema defines the schema for the data source.
func (d *ProviderStatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Status information about the provider",
		Attributes: map[string]schema.Attribute{
			"provider_version": schema.StringAttribute{
				Description: "Version of the provider",
				Computed:    true,
			},
			"go_version": schema.StringAttribute{
				Description: "Version of Go used to build the provider",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *ProviderStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerStatusModel

	// Set computed attributes
	state.ProviderVersion = types.StringValue("dev") // This would normally be set via build flags
	state.GoVersion = types.StringValue(runtime.Version())

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// providerStatusModel is the data source implementation model.
type providerStatusModel struct {
	ProviderVersion types.String `tfsdk:"provider_version"`
	GoVersion       types.String `tfsdk:"go_version"`
}
