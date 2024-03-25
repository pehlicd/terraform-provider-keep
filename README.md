# terraform-provider-keep

[![docs](https://img.shields.io/static/v1?label=docs&message=terraform&color=informational&style=for-the-badge&logo=terraform)](https://registry.terraform.io/providers/pehlicd/keep/latest/docs)
![downloads](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fregistry.terraform.io%2Fv2%2Fproviders%2Fpehlicd%2Fkeep%3Finclude%3Dcategories%2Cmoved-to%2Cpotential-fork-of%2Cprovider-versions%2Ctop-modules%26include%3Dcategories%252Cmoved-to%252Cpotential-fork-of%252Cprovider-versions%252Ctop-modules%26name%3Dkeep%26namespace%3Dpehlicd&query=data.attributes.downloads&style=for-the-badge&logo=terraform&label=downloads&color=brightgreen)
![latest version](https://img.shields.io/github/v/release/pehlicd/terraform-provider-keep?style=for-the-badge&label=latest%20version&color=orange)

This is a terraform provider for managing your [keep](https://github.com/keephq/keep).

> **Note:** This provider is not official terraform provider for keep.

### Installation

Add the following to your terraform configuration

```tf
terraform {
  required_providers {
    keep = {
      source  = "pehlicd/keep"
      version = "~> 0.0.1"
    }
  }
}
```

### Example

```hcl
provider "keep" {
  backend_url = "http://localhost:8080" # or use environment variable KEEP_BACKEND_URL
  api_key = "your apikey" # or use environment variable KEEP_API_KEY
}

resource "keep_workflow" "example_workflow" {
  workflow_file_path = "path/to/workflow.yml"
}

resource "keep_mapping" "example_mapping" {
  name = "example_mapping"
  mapping_file_path = "path/to/mapping.yml"
  matchers = [
    "your unique matcher",
  ]
}

resource "keep_provider" "example_provider" {
  name = "example_provider"
  type = "supported_provider_type"
  auth_config = {
    //...
    // Add your provider specific configuration
    //...
  }
}

data "keep_workflow" "example_workflow_data" {
  id = keep_workflow.example_workflow.id
}

data "keep_mapping" "example_mapping_data" {
  id = keep_mapping.example_mapping.id
}
```

For more information, please refer to the [documentation](https://registry.terraform.io/providers/pehlicd/keep/latest/docs).
