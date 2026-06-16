---
page_title: "softcat_azure_subscription Resource - Softcat"
subcategory: ""
description: |-
  Manages a Softcat Azure subscription order and subscription metadata.
---

# softcat_azure_subscription

Manages a Softcat Azure subscription. The resource supports two workflows:

- **New subscriptions** — placed via the Softcat API using `checkout_data`. The provider looks up the resulting Azure subscription by order ID after creation.
- **CSP-moved subscriptions** — existing Azure subscriptions that have been moved to Softcat CSP and already have a known subscription ID. These can be imported directly and managed without `checkout_data`.

## Example Usage

### New subscription

```terraform
terraform {
  required_providers {
    softcat = {
      source = "lcp-llp/softcat"
    }
  }
}

resource "softcat_azure_subscription" "example" {
  basket_name    = "test10-10-2025"
  msid           = var.msid
  azure_budget   = "1"
  friendly_name  = "TestApiAzureSub"
  azure_contact  = var.azure_contact
  quantity       = 1

  checkout_data {
    purchase_order_number  = "test10-10-2025"
    additional_information = var.azure_contact
    csp_terms              = true
  }
}
```

### CSP-moved subscription (import only)

```terraform
resource "softcat_azure_subscription" "csp_moved" {
  msid          = var.msid
  friendly_name = "MyExistingSubscription"
  azure_contact = var.azure_contact
}
```

## Schema

### Required

- `msid` (String) Microsoft tenant ID required by the Softcat API.
- `azure_contact` (String) Primary contact email sent to the Softcat Azure subscription mutation. This is also used as the budget contact for updates.
- `friendly_name` (String) Friendly name for the Azure subscription.


### Optional

- `basket_name` (String) Optional basket name used when placing the Azure subscription order.
- `quantity` (Number) Quantity passed to the order mutation. Defaults to `1`.
- `checkout_data` (Block List, Max: 1) Checkout metadata required when placing a new subscription order. Not required for CSP-moved subscriptions that are imported.
- `azure_budget` (String) Budget value sent to the Softcat Azure subscription mutation.
### Read-Only

- `id` (String) Terraform resource ID. This is the Azure `subscription_id`.
- `order_id` (String) Softcat order identifier returned by the mutation.
- `subscription_id` (String) Azure subscription identifier.
- `status` (String) Order or subscription status returned by the API.
- `order_name` (String) Order name returned by the Softcat API.
- `currency` (String) Currency associated with the created order.
- `po_number` (String) Purchase order number echoed by the API.
- `date_created` (String) Timestamp when the order was created.
- `date_stored` (String) Timestamp when the order was stored.
- `creator` (Block List, Max: 1) Creator metadata returned by the Softcat API.

### Nested Schema for `checkout_data`

Required:

- `purchase_order_number` (String) Purchase order number to include in the checkout payload.
- `csp_terms` (Boolean) Confirms acceptance of CSP terms for the order.

Optional:

- `additional_information` (String) Optional free-form information sent with the checkout payload.

### Nested Schema for `creator`

Read-Only:

- `user_id` (String) Creator user ID.
- `name` (String) Creator name.
- `email` (String) Creator email address.
- `account` (String) Creator account identifier.

## Import

Import is supported using the Microsoft tenant ID and Azure subscription ID. This is the primary workflow for **CSP-moved subscriptions** — existing Azure subscriptions that have been transferred to Softcat CSP and already have a known subscription ID. After import the resource is managed in-place without needing `checkout_data`.

```terraform
terraform import softcat_azure_subscription.example <msid>/<subscription_id>
```

Example:

```terraform
terraform import softcat_azure_subscription.example 00000000-0000-0000-0000-000000000000/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```
