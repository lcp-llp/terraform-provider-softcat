terraform {
  required_providers {
    softcat = {
      source = "softcat/softcat"
    }
  }
}

provider "softcat" {
  hostname = var.softcat_hostname
  username = var.softcat_username
  password = var.softcat_password
}

variable "softcat_hostname" {
  type        = string
  description = "GraphQL endpoint for the Softcat API."
}

variable "softcat_username" {
  type        = string
  description = "Username used for Softcat API authentication."
}

variable "softcat_password" {
  type        = string
  sensitive   = true
  description = "Password used for Softcat API authentication."
}

variable "msid" {
  type        = string
  description = "Microsoft tenant ID required by the Softcat API."
}

variable "azure_contact" {
  type        = string
  description = "Primary contact email for the Azure subscription order."
}

resource "softcat_azure_subscription" "example" {
  basket_name    = "test10-10-2025"
  msid           = var.msid
  azure_budget   = "1"
  azure_nickname = "TestApiAzureSub"
  azure_contact  = var.azure_contact
  quantity       = 1

  checkout_data {
    purchase_order_number = "test10-10-2025"
    additional_information = var.azure_contact
    csp_terms              = true
  }
}

# Import format:
# terraform import softcat_azure_subscription.example <msid>/<order_id>
# Example:
# terraform import softcat_azure_subscription.example 00000000-0000-0000-0000-000000000000/SC-ORDER-12345

output "softcat_order_id" {
  value = softcat_azure_subscription.example.order_id
}

output "softcat_order_status" {
  value = softcat_azure_subscription.example.status
}