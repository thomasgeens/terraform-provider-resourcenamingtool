// Copyright (c) Thomas Geens

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProviderStatusDataSource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
data "resourcenamingtool_status" "test" {}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the datasource has the expected attributes
					resource.TestCheckResourceAttrSet(
						"data.resourcenamingtool_status.test",
						"provider_version",
					),
					resource.TestCheckResourceAttrSet(
						"data.resourcenamingtool_status.test",
						"go_version",
					),
				),
			},
		},
	})
}
