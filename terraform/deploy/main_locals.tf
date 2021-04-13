locals {
  ansible_args_vault_password_file = (var.ansible_vault_password_file != "") ? "--vault-password-file ${var.ansible_vault_password_file}" : ""
  ansible_args_vault_files         = join(" ", formatlist("-e @%s", var.ansible_vault_files))
  ansible_args = join(" ", [
    local.ansible_args_vault_password_file,
    local.ansible_args_vault_files,
  ])

  letsencrypt_deploy_arguments_delay           = (var.delay != "") ? "-delay ${var.delay}" : ""
  letsencrypt_deploy_arguments_hook            = join(" ", formatlist("-hook \"%s\"", var.hooks))
  letsencrypt_deploy_arguments_output_location = (var.output_location != "") ? "-o ${var.output_location}" : ""

  template_variables = {
    config                           = local.letsencrypt_config
    create_systemd_timer             = var.create_systemd_timer
    instance_ip                      = var.instance_ip
    letsencrypt_deploy_version       = var.letsencrypt_deploy_version
    letsencrypt_deploy_checksum_type = var.letsencrypt_deploy_checksum_type
    letsencrypt_deploy_checksum      = var.letsencrypt_deploy_checksum
    start_systemd_service            = var.start_systemd_service
    additional_letsencrypt_deploy_arguments = join(" ",
      compact([
        local.letsencrypt_deploy_arguments_hook,
      ])
    ),
  }

  passphrase = var.client_passphrase != "" ? var.client_passphrase : "{{ ${var.client_passphrase_key} }}"

  fortios_access_token = var.fortios_access_token != "" ? var.fortios_access_token : (var.fortios_access_token_key != "" ? "{{ ${var.fortios_access_token_key} }}" : "")

  # create map with config values, remove empty config settings
  letsencrypt_config = merge([for key, val in
    {
      delay : var.delay,
      domains : var.domains,
      dynamodbTableName : var.dynamodb_table_name,
      email : var.email,
      outputLocation : var.output_location,
      passphrase : local.passphrase,
      prefix : var.prefix,
      local : var.local,
      fortios : var.fortios,
      fortiosAccessToken : local.fortios_access_token,
      fortiosAdminServerCert : var.fortios_admin_server_cert,
      fortiosBaseUrl : var.fortios_base_url,
      fortiosSslSshProfiles : var.fortios_ssl_ssh_profiles,
    } : ((jsonencode(val) != "\"\"") ? { format("%s", key) : val } : {})]...
  )
}
