variable "azure_subscription_id" {
  description = "Azure Subscription ID"
  type        = string
  default     = "00000000-0000-0000-0000-000000000000"
}

variable "azure_client_id" {
  description = "Azure Client ID"
  type        = string
  default     = "11111111-1111-1111-1111-111111111111"
}

variable "azure_client_secret" {
  description = "Azure Client Secret"
  type        = string
  default     = "dummy-secret"
  sensitive   = true
}

variable "azure_tenant_id" {
  description = "Azure Tenant ID"
  type        = string
  default     = "22222222-2222-2222-2222-222222222222"
}
