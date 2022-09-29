terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

variable "vault_instances" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances that were created"
}

variable "vault_product_version" {
  type        = string
  description = "The vault product version"
  default     = null
}

variable "vault_edition" {
  type        = string
  description = "The vault product edition"
  default     = null
}

variable "vault_revision" {
  type        = string
  description = "The vault product revision"
  default     = null
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
  default     = null
}

locals {
  instances = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.vault_instances)[idx].public_ip
      private_ip = values(var.vault_instances)[idx].private_ip
    }
  }
}

resource "enos_remote_exec" "verify_all_nodes_have_updated_version" {
  for_each = local.instances

  content = templatefile("${path.module}/templates/verify-cluster-version.sh", {
    vault_install_dir = var.vault_install_dir,
    vault_version     = var.vault_product_version,
    vault_edition     = var.vault_edition,
    vault_revision    = var.vault_revision,
    vault_token       = var.vault_root_token,
  })

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
