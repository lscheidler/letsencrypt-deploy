resource "null_resource" "deploy" {
  triggers = {
    instance_id = var.instance_id
    variables   = jsonencode(local.template_variables)
  }

  provisioner "local-exec" {
    command = data.template_file.deploy.rendered
  }
}
