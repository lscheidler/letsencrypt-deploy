# letsencrypt-deploy

Terraform module for deploying letsencrypt-deploy to an instance

# Requirements

- [ansible](https://www.ansible.com/) (>= 2.9.13)
- [terraform](https://www.terraform.io) (>= 0.13.5)
- working ssh access to instance with super user permissions (sudo)

## Usage

```
module "letsencrypt-deploy" {
  source = "github.com/lscheidler/letsencrypt-deploy//terraform/deploy?ref=main"

  instance_ip = aws_instance.instance.private_ip
  instance_id = aws_instance.instance.id

  dependencies = [
    aws_instance.instance.private_ip,
    null_resource.prerequirement.id,
  ]

  #delay           = "10"
  domains         = "example.com,*.example.com"
  email           = "me@example.com"
  hooks = [
    "exec;systemctl restart nginx.service",
    "sns;arn:aws:sns:<region>:<account-id>:<topic>",
  ]
  
  #client_passphrase = "<client-passphrase>"

  client_passphrase_key       = "vault_client_passphrase"
  ansible_vault_file          = "vault_letsencrypt.yml"
  ansible_vault_password_file = "vault-passfile"
}
```

It is going to deploy and configure
- letsencrypt-deploy

## Argument Reference

| Name                                    | Required  | Default                                   | Description                                             |
|-----------------------------------------|-----------|-------------------------------------------|---------------------------------------------------------|
| `domains`                               | 🗹         |                                           | Domains to get certificates for                         |
| `email`                                 | 🗹         |                                           | Registration email for letsencrypt                      |
| `instance_ip`                           | 🗹         |                                           |                                                         |
| `ansible_vault_file`                    | 🗷         | `""`                                      | should be set, if `client_passphrase_key` is set        |
| `ansible_vault_password_file`           | 🗷         | `""`                                      | should be set, if `client_passphrase_key` is set        |
| `client_passphrase`                     | 🗷         | `""`                                      | this argument or `client_passphrase_key` should be set  |
| `client_passphrase_key`                 | 🗷         | `""`                                      |                                                         |
| `create_systemd_timer`                  | 🗷         | `true`                                    |                                                         |
| `delay`                                 | 🗷         | `""`                                      | set delay argument for letsencrypt-deploy               |
| `dependencies`                          | 🗷         | `[]`                                      | add additional dependencies to wait for                 |
| `hooks`                                 | 🗷         | `[]`                                      | add hooks                                               |
| `output_location`                       | 🗷         | `""`                                      | set `output_location` for letsencrypt-deploy            |
