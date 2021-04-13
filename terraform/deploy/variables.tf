################################################################################
# data.template_file.deploy
################################################################################
variable "ansible_vault_files" {
  type    = list(string)
  default = []
}
variable "ansible_vault_password_file" {
  type    = string
  default = ""
}
variable "client_passphrase" {
  type    = string
  default = ""
}
variable "client_passphrase_key" {
  type    = string
  default = ""
}
variable "create_systemd_timer" {
  type    = bool
  default = true
}

variable "delay" {
  type    = number
  default = 0
}

variable "domains" {
  type = list(any)
}

variable "dynamodb_table_name" {
  type    = string
  default = ""
}

variable "email" {}

variable "hooks" {
  type    = list(any)
  default = []
}

variable "instance_ip" {}

variable "output_location" {
  type    = string
  default = ""
}

variable "prefix" {
  type    = string
  default = ""
}

variable "start_systemd_service" {
  type    = bool
  default = true
}

variable "local" {
  type    = bool
  default = true
}

variable "fortios" {
  type    = bool
  default = false
}

variable "fortios_admin_server_cert" {
  type    = bool
  default = true
}

variable "fortios_base_url" {
  type    = string
  default = ""
}

variable "fortios_ssl_ssh_profiles" {
  type    = list(any)
  default = []
}

variable "fortios_access_token" {
  type    = string
  default = ""
}

variable "fortios_access_token_key" {
  type    = string
  default = ""
}

################################################################################
# null_resource.deploy
################################################################################
variable "instance_id" {
  default = ""
}
