// Copyright (c) Thomas Geens
// SPDX-License-Identifier: MIT

//go:build generate

package tools

import (
	_ "github.com/hashicorp/copywrite"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)

// Generate copyright headers
//go:generate go run github.com/hashicorp/copywrite headers -d .. --config ../.copywrite.hcl

// Format Terraform code for use in documentation.
// If you do not have Terraform installed, you can remove the formatting command, but it is suggested
// to ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ../examples/

// Generate documentation.
// Clean up the old documentation first.
//go:generate rm -rf ../docs/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir .. -provider-name resourcenamingtool
