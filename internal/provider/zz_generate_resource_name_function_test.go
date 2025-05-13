// Copyright (c) Thomas Geens

package provider

import (
	"regexp"
	"testing"

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
  value = provider::resourcenamingtool::generate_resource_name([])
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
