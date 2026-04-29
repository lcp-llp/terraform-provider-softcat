---
page_title: "Softcat Provider"
subcategory: ""
description: |-
  The Softcat provider configures access to the Softcat GraphQL API.
---

# Softcat Provider

The Softcat provider is used to manage Softcat resources through the Softcat GraphQL API.

## Example Usage

```terraform
terraform {
  required_providers {
    softcat = {
      source = "lcp-llp/softcat"
    }
  }
}

provider "softcat" {
  hostname = var.softcat_hostname
  username = var.softcat_username
  password = var.softcat_password
}
```

## Schema

### Required

- `hostname` (String) GraphQL endpoint for the Softcat API. Can also be set with the `SOFTCAT_HOSTNAME` environment variable.
- `username` (String) Username used for authentication to the Softcat API. Can also be set with the `SOFTCAT_USERNAME` environment variable.
- `password` (String, Sensitive) Password used for authentication to the Softcat API. Can also be set with the `SOFTCAT_PASSWORD` environment variable.
