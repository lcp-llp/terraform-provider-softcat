---
page_title: "softcat_azure_subscription Resource - Softcat"
subcategory: ""
description: |-
  Manages a Softcat Azure subscription order and subscription metadata.
---

# softcat_azure_subscription

Manages a Softcat Azure subscription order. After creation, the provider also looks up the Azure subscription details for the created order and exposes the resulting subscription metadata in state.

## Example Usage

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
  azure_nickname = "TestApiAzureSub"
  azure_contact  = var.azure_contact
  quantity       = 1

  checkout_data {
    purchase_order_number  = "test10-10-2025"
    additional_information = var.azure_contact
    csp_terms              = true
  }
}
```

## Schema

### Required

- `msid` (String) Microsoft tenant ID required by the Softcat API.
- `azure_budget` (String) Budget value sent to the Softcat Azure subscription mutation.
- `azure_contact` (String) Primary contact email sent to the Softcat Azure subscription mutation. This is also used as the budget contact for updates.
- `azure_nickname` (String) Friendly name for the Azure subscription. This is the same underlying value returned as the subscription display name.
- `checkout_data` (Block List, Min: 1, Max: 1) Checkout metadata required by the create mutation.

### Optional

- `basket_name` (String) Optional basket name used when placing the Azure subscription order.
- `quantity` (Number) Quantity passed to the order mutation. Defaults to `1`.

### Read-Only

- `id` (String) Terraform resource ID. This is the Softcat `order_id`.
- `order_id` (String) Softcat order identifier returned by the mutation.
- `subscription_id` (String) Azure subscription identifier.
- `display_name` (String) Azure subscription display name.
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

Import is supported using the Microsoft tenant ID and Softcat order ID.

```terraform
terraform import softcat_azure_subscription.example <msid>/<order_id>
```

Example:

```terraform
terraform import softcat_azure_subscription.example 00000000-0000-0000-0000-000000000000/SC-ORDER-12345
```
