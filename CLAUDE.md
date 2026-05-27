# CLAUDE.md

Guidance for Claude Code (claude.ai/code) when working with this repo.

## What this project is

Terraform provider (`terraform-provider-resourcenamingtool`) exposing single provider function — `generate_resource_name` — for uniformly named cloud resources. Uses [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) (not older SDK).

## Commands

```bash
make test          # Unit tests only (TF_ACC=0)
make testacc       # Acceptance tests (TF_ACC=1, requires Terraform)
make build         # Runs test + testacc, then builds the binary
make install       # build + copy binary to GOPATH/bin + update ~/.terraformrc
make verify        # install + run terraform plan in examples/provider
make generate      # Re-generate docs (runs tfplugindocs, copywrite, terraform fmt)
make fmt           # go fmt ./...
make lint          # golangci-lint run
make clean         # remove binary, GOPATH copy, and internal/provider/.terraform

# After changes to examples/, internal/provider/descriptions/, or provider schema (*.go):
make generate      # regenerate docs — must be committed or CI generate job fails

# Run a specific test by name
TF_ACC=0 go test ./internal/provider -run=TestGenerateResourceNameFunction -v
TF_ACC=1 go test ./internal/provider -run=TestAccProviderStatusDataSource -v
```

## Architecture

### Provider lifecycle

`ValidateConfig` → `Configure` → `Functions()` — Terraform call order. Key design: `ValidateConfig` and function `Run` execute in the **same process**; `f.config` (set during `Functions()`) is the primary config source — no file I/O on the normal path. File-based fallback (`GetSharedProviderConfig`) is used only when `f.config` is nil (e.g. function evaluated before `ValidateConfig` runs). Config files stored in `$TF_DATA_DIR/resourcenamingtool_cache/` → `{cwd}/.terraform/resourcenamingtool_cache/` → `os.TempDir()/resourcenamingtool_cache/` (fallback chain).

- `provider.go` — provider registration, schema, `ValidateConfig` (saves config to file), `Configure` (loads config from file into `p.config`), `Functions()` (wires `p.config` into function instance as `f.config`).
- `config_manager.go` — file I/O: `saveProviderConfigToFile` / `loadProviderConfigFromFile`. Uses `sync.Mutex` for in-process safety and `github.com/gofrs/flock` for cross-process file locking. Config files named `resourcenamingtool_config_{instanceID}.json` (or `_default.json` when `provider_instance_id` unset).
- `zz_generate_resource_name_function.go` — `generate_resource_name` function. `Run` uses `f.config` (in-memory) as primary; falls back to file only if `f.config` nil. Merges provider-level defaults with call-time parameters, applies naming pattern, returns generated string.
- `componentvalue_type.go` — `ComponentValueType` / `ComponentValueObject`: custom Terraform attribute type for naming component with three representations: `fullname`, `shortcode`, `char`. Used for every `default_*` provider attribute and function parameters.
- `resourcenamingparameters_type.go` — `ResourceNamingParametersType`: custom object type for function's `parameters` argument; mirrors all provider-level component fields plus `additional_components` and `additional_naming_patterns`.
- `provider_status_datasource.go` — read-only data source exposing provider config state for debugging.
- `logging_helpers.go` — thin wrappers around `terraform-plugin-log` (`logDebug`, `logInfo`, `logWarn`, `logError`, `logDebugWithFields`, `logErrorWithFields`).

### Naming pattern resolution

Patterns use `{component_name}` placeholders. `:short` and `:char` suffixes (e.g. `{environment:short}`, `{region:char}`) select `shortcode` or `char` representation of component. Resolved in `zz_generate_resource_name_function.go`. Built-in cloud patterns defined (currently commented out) in `config_manager.go` under `builtin_NamingPatterns`. Custom patterns injectable via `additional_naming_patterns`.

### Multi-instance support

Multiple provider blocks coexist if each sets unique `provider_instance_id`. Only needed when two provider blocks with different defaults exist in the same workspace — value becomes config filename suffix, keeping fallback files isolated on disk.

### Documentation generation

`tools/tools.go` uses `//go:generate` directives. Run `make generate` to regenerate docs under `docs/` via `tfplugindocs`. Description strings embedded from `internal/provider/descriptions/*.txt` and `*.md` files.
