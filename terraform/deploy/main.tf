resource "null_resource" "deploy" {
  triggers = {
    instance_id = var.instance_id
  }

  provisioner "local-exec" {
    command = data.template_file.deploy.rendered
  }
}
