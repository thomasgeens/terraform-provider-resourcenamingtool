// Copyright (c) Thomas Geens

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestProvider(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Provider test - verify provider instance can be instantiated
			{
				Config: providerConfig + `
output "provider_test" {
  value = "Provider successfully configured"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("provider_test", "Provider successfully configured"),
				),
			},
		},
	})
}
