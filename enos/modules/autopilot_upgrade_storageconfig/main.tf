variable "vault_product_version" {}

output "storage_addl_config" {
  value = {
    autopilot_upgrade_version = var.vault_product_version
  }
}
