// Copyright Thomas Geens 2025, 2026

package provider

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestGenerateResourceNameFunction_Basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
output "test" {
  # depends_on forces evaluation after the provider is configured.
  # Without this, Terraform 1.8+ may call the function before ValidateProviderConfig
  # runs (empty-literal arguments have no graph dependency on the provider), so
  # p.config would be nil and default_resource_type would not be available.
  depends_on = [data.resourcenamingtool_status.init]
  value      = provider::resourcenamingtool::generate_resource_name([])
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the output has the expected format for resource group
					resource.TestMatchOutput("test", regexp.MustCompile(`^rg-example-prd-we$`)),
				),
			},
		},
	})
}

func TestGenerateResourceNameFunction_WithAdditionalComponents(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
output "test" {
  value = provider::resourcenamingtool::generate_resource_name([{
	resource_type = {
		"fullname"  = "azurerm_resource_group"
		"shortcode" = "rg"
		"char"      = "r"
	},
	additional_components = {
		"instance.fullname"  = "00002"
		"instance.shortcode" = "002"
		"instance.char"      = "2"
	},
	additional_naming_patterns = {
		"azurerm_resource_group" = "{basename}-{environment:short}-{region:short}-{instance:char}"
	}
	}])
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the output has the expected format with custom pattern and instance
					resource.TestMatchOutput("test", regexp.MustCompile(`^example-prd-we-2$`)),
				),
			},
		},
	})
}

func TestGenerateResourceNameFunction_WithRegionOverride(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
output "test" {
  value = provider::resourcenamingtool::generate_resource_name([{
		additional_components = {
			"region.fullname"  = "Germany West Central"
			"region.shortcode" = "gwc"
			"region.char"      = "g"
		}
	}])
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the output has the expected format with custom region
					resource.TestMatchOutput("test", regexp.MustCompile(`^rg-example-prd-gwc$`)),
				),
			},
		},
	})
}

func TestGenerateResourceNameFunction_StorageAccount(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
output "test" {
  value = provider::resourcenamingtool::generate_resource_name([{
      resource_type = {
        "fullname"  = "azurerm_storage_account"
      }
    }])
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the output has the expected format for storage account (no hyphens)
					resource.TestMatchOutput("test", regexp.MustCompile(`^examplepw00001$`)),
				),
			},
		},
	})
}

func TestGenerateResourceNameFunction_LocationComponent(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// {location:short} resolves via the location componentType (separate from region).
				// Pass location as additional_components to isolate the placeholder logic.
				Config: providerConfig + `
output "test" {
  value = provider::resourcenamingtool::generate_resource_name([{
    resource_type = { "fullname" = "myresource" },
    additional_naming_patterns = {
      "myresource" = "res-{location:short}"
    },
    additional_components = {
      "location.fullname"  = "Germany West Central"
      "location.shortcode" = "gwc"
      "location.char"      = "g"
    }
  }])
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchOutput("test", regexp.MustCompile(`^res-gwc$`)),
				),
			},
		},
	})
}

func TestGenerateResourceNameFunction_ProviderLevelAdditionalComponents(t *testing.T) {
	// Direct Go unit test — bypasses the Terraform test framework to avoid
	// file-mechanism cross-test pollution on the default config file.
	ctx := context.Background()

	deptComp, diags := CreateComponentValueObjectFromParts(ctx, "finance", "fin", "f")
	if diags.HasError() {
		t.Fatalf("CreateComponentValueObjectFromParts: %s", diags.Errors()[0].Summary())
	}

	additionalComponents, mapDiags := types.MapValueFrom(ctx, NewComponentValueType(), map[string]attr.Value{
		"{dept}": deptComp,
	})
	if mapDiags.HasError() {
		t.Fatalf("MapValueFrom(AdditionalComponents): %s", mapDiags.Errors()[0].Summary())
	}

	additionalPatterns, mapDiags := types.MapValueFrom(ctx, types.StringType, map[string]attr.Value{
		"myresource": types.StringValue("res-{dept}"),
	})
	if mapDiags.HasError() {
		t.Fatalf("MapValueFrom(AdditionalNamingPatterns): %s", mapDiags.Errors()[0].Summary())
	}

	config := resourcenamingtoolProviderModel{
		AdditionalComponents:     additionalComponents,
		AdditionalNamingPatterns: additionalPatterns,
	}

	resourceTypeComp, compDiags := CreateComponentValueObjectFromParts(ctx, "myresource", "mr", "m")
	if compDiags.HasError() {
		t.Fatalf("CreateComponentValueObjectFromParts(resource_type): %s", compDiags.Errors()[0].Summary())
	}

	params, objDiags := types.ObjectValue(
		map[string]attr.Type{"resource_type": NewComponentValueType()},
		map[string]attr.Value{"resource_type": resourceTypeComp},
	)
	if objDiags.HasError() {
		t.Fatalf("ObjectValue(params): %s", objDiags.Errors()[0].Summary())
	}

	result, resultDiags := generateResourceName(ctx, ResourceNamingParametersValue{ObjectValue: params}, config)
	if resultDiags.HasError() {
		t.Fatalf("generateResourceName: %s", resultDiags.Errors()[0].Summary())
	}
	if result != "res-finance" {
		t.Errorf("expected res-finance, got %s", result)
	}
}

func TestGenerateResourceNameFunction_FunctionLevelOverridesProviderAdditionalComponents(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Function-level additional_components override provider-level
				Config: `
provider "resourcenamingtool" {
  additional_components = {
    "{dept}" = {
      fullname  = "finance"
      shortcode = "fin"
      char      = "f"
    }
  }
  additional_naming_patterns = {
    "myresource" = "res-{dept}"
  }
}
data "resourcenamingtool_status" "init" {}
output "test" {
  depends_on = [data.resourcenamingtool_status.init]
  value      = provider::resourcenamingtool::generate_resource_name([{
    resource_type = { "fullname" = "myresource" },
    additional_components = {
      "dept.fullname"  = "engineering"
      "dept.shortcode" = "eng"
      "dept.char"      = "e"
    }
  }])
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchOutput("test", regexp.MustCompile(`^res-engineering$`)),
				),
			},
		},
	})
}

func TestGenerateResourceNameFunction_Null(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
output "test" {
  value = provider::resourcenamingtool::generate_resource_name()
}
`,
				ExpectError: regexp.MustCompile(`Error: Not enough function arguments`),
			},
		},
	})
}
