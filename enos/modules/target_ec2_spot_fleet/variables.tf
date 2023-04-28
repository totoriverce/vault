variable "ami_id" {
  description = "The machine image identifier"
  type        = string
}

variable "awskms_unseal_key_arn" {
  type        = string
  description = "The AWSKMS key ARN if using the awskms unseal method. If specified the instances will be granted kms permissions to the key"
  default     = null
}

variable "cluster_name" {
  type        = string
  description = "A unique cluster identifier"
  default     = null
}

variable "common_tags" {
  description = "Common tags for cloud resources"
  type        = map(string)
  default = {
    Project = "Vault"
  }
}

variable "instance_mem_min" {
  description = "The minimum amount of memory in mebibytes for each instance in the fleet. (1 MiB = 1024 bytes)"
  type        = number
  default     = 4096 // ~4 GB
}

variable "instance_mem_max" {
  description = "The maximum amount of memory in mebibytes for each instance in the fleet. (1 MiB = 1024 bytes)"
  type        = number
  default     = 16385 // ~16 GB
}

variable "instance_cpu_min" {
  description = "The minimum number of vCPU's for each instance in the fleet"
  type        = number
  default     = 2
}

variable "instance_cpu_max" {
  description = "The maximum number of vCPU's for each instance in the fleet"
  type        = number
  default     = 8 // Unlikely we'll ever get that high due to spot price bid protection
}

variable "instance_count" {
  description = "The number of target instances to create"
  type        = number
  default     = 3
}

variable "instance_type" {
  description = "Shim variable for target module variable compatibility that is not used. The spot fleet determines instance sizes"
  type        = string
  default     = null
}

variable "project_name" {
  description = "A unique project name"
  type        = string
}

variable "spot_price_max" {
  description = "The maximum hourly price to pay for each target instance"
  type        = string
  // Current on-demand cost of linux t3.medium in us-east.
  default = "0.0416"
}

variable "ssh_allow_ips" {
  description = "Allowlisted IP addresses for SSH access to target nodes. The IP address of the machine running Enos will automatically allowlisted"
  type        = list(string)
  default     = []
}

variable "ssh_keypair" {
  description = "SSH keypair used to connect to EC2 instances"
  type        = string
}

variable "vpc_id" {
  description = "The identifier of the VPC where the target instances will be created"
  type        = string
}
