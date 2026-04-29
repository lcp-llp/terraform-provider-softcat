# Terraform Provider Softcat

Terraform provider for managing Softcat resources through the Softcat GraphQL API.

## Current Support

The provider currently supports:

- `softcat_azure_subscription`

## Provider Source

```hcl
terraform {
  required_providers {
    softcat = {
      source = "lcp-llp/softcat"
    }
  }
}
```

## Provider Configuration

```hcl
provider "softcat" {
  hostname = var.softcat_hostname
  username = var.softcat_username
  password = var.softcat_password
}
```

The provider accepts these configuration arguments:

- `hostname` - GraphQL endpoint for the Softcat API.
- `username` - Username used for API authentication.
- `password` - Password used for API authentication.

These can also be set with environment variables:

- `SOFTCAT_HOSTNAME`
- `SOFTCAT_USERNAME`
- `SOFTCAT_PASSWORD`

## Documentation

- Provider docs: [docs/index.md](docs/index.md)
- Resource docs: [docs/resources/azure_subscription.md](docs/resources/azure_subscription.md)
- Example configuration: [examples/resources/softcat_azure_subscription/resource.tf](examples/resources/softcat_azure_subscription/resource.tf)
