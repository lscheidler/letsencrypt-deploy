resource "null_resource" "deploy" {
  triggers = {
    instance_id = var.instance_id
    variables   = sha512(jsonencode(local.template_variables))
  }

  provisioner "local-exec" {
    command = join(" ", [
      "ansible-playbook",
      local.ansible_args,
      "${path.module}/deploy.yml",
      "-e",
      join("", [
        "'",
        jsonencode(local.template_variables),
        "'",
      ]),
    ])
  }
}
