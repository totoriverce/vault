# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

scenario "agent" {
  matrix {
    arch            = ["amd64", "arm64"]
    artifact_source = ["local", "crt", "artifactory"]
    distro          = ["ubuntu", "rhel"]
    edition         = ["oss", "ent", "ent.fips1402", "ent.hsm", "ent.hsm.fips1402"]
  }

  terraform_cli = terraform_cli.default
  terraform     = terraform.default
  providers = [
    provider.aws.default,
    provider.enos.ubuntu,
    provider.enos.rhel
  ]

  locals {
    build_tags = {
      "oss"              = ["ui"]
      "ent"              = ["ui", "enterprise", "ent"]
      "ent.fips1402"     = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.fips1402"]
      "ent.hsm"          = ["ui", "enterprise", "cgo", "hsm", "venthsm"]
      "ent.hsm.fips1402" = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.hsm.fips1402"]
    }
    bundle_path = matrix.artifact_source != "artifactory" ? abspath(var.vault_bundle_path) : null
    packages    = ["jq"]
    enos_provider = {
      rhel   = provider.enos.rhel
      ubuntu = provider.enos.ubuntu
    }
    install_artifactory_artifact = local.bundle_path == null
    spot_price_max = {
      // These prices are based on on-demand cost for t3.medium in us-east
      "rhel"   = "0.1016"
      "ubuntu" = "0.0416"
    }
    tags = merge({
      "Project Name" : var.project_name
      "Project" : "Enos",
      "Environment" : "ci"
    }, var.tags)
    vault_instance_types = {
      amd64 = "t3a.small"
      arm64 = "t4g.small"
    }
    vault_instance_type = coalesce(var.vault_instance_type, local.vault_instance_types[matrix.arch])
    vault_license_path  = abspath(var.vault_license_path != null ? var.vault_license_path : joinpath(path.root, "./support/vault.hclic"))
  }

  step "get_local_metadata" {
    skip_step = matrix.artifact_source != "local"
    module    = module.get_local_metadata
  }

  step "build_vault" {
    module = "build_${matrix.artifact_source}"

    variables {
      build_tags           = var.vault_local_build_tags != null ? var.vault_local_build_tags : local.build_tags[matrix.edition]
      bundle_path          = local.bundle_path
      goarch               = matrix.arch
      goos                 = "linux"
      artifactory_host     = matrix.artifact_source == "artifactory" ? var.artifactory_host : null
      artifactory_repo     = matrix.artifact_source == "artifactory" ? var.artifactory_repo : null
      artifactory_username = matrix.artifact_source == "artifactory" ? var.artifactory_username : null
      artifactory_token    = matrix.artifact_source == "artifactory" ? var.artifactory_token : null
      arch                 = matrix.artifact_source == "artifactory" ? matrix.arch : null
      product_version      = var.vault_product_version
      artifact_type        = matrix.artifact_source == "artifactory" ? var.vault_artifact_type : null
      distro               = matrix.artifact_source == "artifactory" ? matrix.distro : null
      edition              = matrix.artifact_source == "artifactory" ? matrix.edition : null
      instance_type        = matrix.artifact_source == "artifactory" ? local.vault_instance_type : null
      revision             = var.vault_revision
    }
  }

  step "find_azs" {
    module = module.az_finder

    variables {
      instance_type = [
        var.backend_instance_type,
        local.vault_instance_type
      ]
    }
  }

  step "create_vpc" {
    module = module.create_vpc

    variables {
      ami_architectures  = distinct([matrix.arch, "amd64"])
      availability_zones = step.find_azs.availability_zones
      common_tags        = local.tags
    }
  }

  step "read_license" {
    skip_step = matrix.edition == "oss"
    module    = module.read_license

    variables {
      file_name = local.vault_license_path
    }
  }

  step "create_vault_cluster_targets" {
    module     = module.target_ec2_spot_fleet // "target_ec2_instances" can be used for on-demand instances
    depends_on = [step.create_vpc]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id                = step.create_vpc.ami_ids[matrix.distro][matrix.arch]
      awskms_unseal_key_arn = step.create_vpc.kms_key_arn
      common_tags           = local.tags
      instance_type         = local.vault_instance_type // only used for on-demand instances
      spot_price_max        = local.spot_price_max[matrix.distro]
      vpc_id                = step.create_vpc.vpc_id
    }
  }

  step "create_vault_cluster" {
    module = module.vault_cluster
    depends_on = [
      step.build_vault,
      step.create_vault_cluster_targets
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      artifactory_release   = matrix.artifact_source == "artifactory" ? step.build_vault.vault_artifactory_release : null
      awskms_unseal_key_arn = step.create_vpc.kms_key_arn
      cluster_name          = step.create_vault_cluster_targets.cluster_name
      config_env_vars = {
        VAULT_LOG_LEVEL = var.vault_log_level
      }
      install_dir         = var.vault_install_dir
      license             = matrix.edition != "oss" ? step.read_license.license : null
      local_artifact_path = local.bundle_path
      packages            = local.packages
      storage_backend     = "raft"
      target_hosts        = step.create_vault_cluster_targets.hosts
      unseal_method       = "shamir"
    }
  }

  step "start_vault_agent" {
    module = "vault_agent"
    depends_on = [
      step.build_vault,
      step.create_vault_cluster,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances                  = step.create_vault_cluster_targets.hosts
      vault_root_token                 = step.create_vault_cluster.root_token
      vault_agent_template_destination = "/tmp/agent_output.txt"
      vault_agent_template_contents    = "{{ with secret \\\"auth/token/lookup-self\\\" }}orphan={{ .Data.orphan }} display_name={{ .Data.display_name }}{{ end }}"
    }
  }

  step "verify_vault_agent_output" {
    module = module.vault_verify_agent_output
    depends_on = [
      step.create_vault_cluster,
      step.start_vault_agent,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances                  = step.create_vault_cluster_targets.hosts
      vault_agent_template_destination = "/tmp/agent_output.txt"
      vault_agent_expected_output      = "orphan=true display_name=approle"
    }
  }

  output "awkms_unseal_key_arn" {
    description = "The Vault cluster KMS key arn"
    value       = step.create_vpc.kms_key_arn
  }

  output "cluster_name" {
    description = "The Vault cluster name"
    value       = step.create_vault_cluster.cluster_name
  }

  output "hosts" {
    description = "The Vault cluster target hosts"
    value       = step.create_vault_cluster.target_hosts
  }

  output "private_ips" {
    description = "The Vault cluster private IPs"
    value       = step.create_vault_cluster.private_ips
  }

  output "public_ips" {
    description = "The Vault cluster public IPs"
    value       = step.create_vault_cluster.public_ips
  }

  output "root_token" {
    description = "The Vault cluster root token"
    value       = step.create_vault_cluster.root_token
  }

  output "recovery_key_shares" {
    description = "The Vault cluster recovery key shares"
    value       = step.create_vault_cluster.recovery_key_shares
  }

  output "recovery_keys_b64" {
    description = "The Vault cluster recovery keys b64"
    value       = step.create_vault_cluster.recovery_keys_b64
  }

  output "recovery_keys_hex" {
    description = "The Vault cluster recovery keys hex"
    value       = step.create_vault_cluster.recovery_keys_hex
  }

  output "unseal_keys_b64" {
    description = "The Vault cluster unseal keys"
    value       = step.create_vault_cluster.unseal_keys_b64
  }

  output "unseal_keys_hex" {
    description = "The Vault cluster unseal keys hex"
    value       = step.create_vault_cluster.unseal_keys_hex
  }
}
