################################################################################
# data.template_file.deploy
################################################################################
variable "ansible_vault_file" {
  type    = string
  default = ""
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

variable "dependencies" {
  type    = list
  default = []
}

variable "delay" {
  type    = string
  default = ""
}

variable "domains" {}

variable "email" {}

variable "hooks" {
  type    = list
  default = []
}

variable "instance_ip" {}

variable "output_location" {
  type    = string
  default = ""
}

################################################################################
# null_resource.deploy
################################################################################
variable "instance_id" {
  default = ""
}
