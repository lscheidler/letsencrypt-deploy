locals {
  ansible_args_vault_password_file = (var.ansible_vault_password_file != "") ? "--vault-password-file ${var.ansible_vault_password_file}" : ""
  ansible_args_vault_file          = (var.ansible_vault_file != "") ? "-e @${var.ansible_vault_file}" : ""
  ansible_args = join(" ", [
    local.ansible_args_vault_password_file,
    local.ansible_args_vault_file,
  ])

  letsencrypt_deploy_arguments_delay           = (var.delay != "") ? "-delay ${var.delay}" : ""
  letsencrypt_deploy_arguments_hook            = join(" ", formatlist("-hook \"%s\"", var.hooks))
  letsencrypt_deploy_arguments_output_location = (var.output_location != "") ? "-o ${var.output_location}" : ""

  template_variables = {
    client_passphrase                = var.client_passphrase,
    client_passphrase_key            = var.client_passphrase_key,
    create_systemd_timer             = var.create_systemd_timer,
    domains                          = var.domains,
    email                            = var.email,
    instance_ip                      = var.instance_ip,
    letsencrypt_deploy_version       = var.letsencrypt_deploy_version,
    letsencrypt_deploy_checksum_type = var.letsencrypt_deploy_checksum_type
    letsencrypt_deploy_checksum      = var.letsencrypt_deploy_checksum
    additional_letsencrypt_deploy_arguments = join(" ",
      compact([
        local.letsencrypt_deploy_arguments_output_location,
        local.letsencrypt_deploy_arguments_hook,
        local.letsencrypt_deploy_arguments_delay,
      ])
    ),
  }
}
