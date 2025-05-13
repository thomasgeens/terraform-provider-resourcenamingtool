# Resource Naming Tool Provider

The **Resource Naming Tool** provider offers a flexible and standardized way to generate resource names across various cloud environments, including Azure, AWS, and GCP. It aligns with established best practices such as Microsoft's Cloud Adoption Framework (CAF), AWS Well-Architected Framework (WAF), and Google Cloud naming conventions.

## Key Features

*   **Consistent Naming**: Enforces uniform naming conventions across your infrastructure, reducing ambiguity and improving resource discoverability.
*   **Multi-Cloud Support**: Comes with built-in, sensible default naming patterns tailored for popular services on Azure, AWS, and GCP.
*   **Customizable Defaults**: Allows you to set default values for common naming components (e.g., `default_environment`, `default_region`, `default_basename`) at the provider level. This simplifies individual `generate_resource_name` function calls by pre-filling common values.
*   **Extensibility**:
    *   `additional_components`: Define your own custom naming components (e.g., `department`, `cost_center_short`) to be used in naming patterns.
    *   `additional_naming_patterns`: Override or add new naming patterns for specific resource types to perfectly match your organization's standards.
*   **Simplified Configuration**: Configure shared settings once at the provider level, and these settings will be available to all `generate_resource_name` function calls, promoting consistency and reducing boilerplate.

This provider helps improve resource organization, simplifies management, and enhances clarity in complex cloud deployments by ensuring that all resources are named predictably and meaningfully. It is particularly useful in environments where maintaining a strict and understandable naming strategy is crucial for operational efficiency and governance.
