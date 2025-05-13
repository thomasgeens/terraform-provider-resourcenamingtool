# Copyright (c) Thomas Geens

# Generate a resource name using solely the provider's configuration:
#   - default_resource_type
#   - default naming patterns
#   - default components
output "azurerm_resource_group_example_1" {
  value = provider::resourcenamingtool::generate_resource_name([])
}

# Generate a resource name using a non-default resource type, in combination with the provider's configuration:
#   - default naming patterns
#   - default components
output "azurerm_storage_account_example_1" {
  value = provider::resourcenamingtool::generate_resource_name([{
    resource_type = {
      "fullname" = "azurerm_storage_account"
    }
  }])
}

# Generate a resource name using a non-default resource type, an additional component, an additional naming pattern, in combination with the provider's configuration:
#   - default components
output "azurerm_resource_group_example_2" {
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

# Generate a resource name using additional components, in combination with the provider's configuration:
#   - default resource type
#   - default naming patterns
#   - default components
output "azurerm_resource_group_example_3" {
  value = provider::resourcenamingtool::generate_resource_name([{
    additional_components = {
      "region.fullname"    = "Germany West Central"
      "region.shortcode"   = "gwc"
      "region.char"        = "g"
      "basename.fullname"  = "custom"
      "basename.shortcode" = "cus"
      "basename.char"      = "c"
    }
  }])
}
