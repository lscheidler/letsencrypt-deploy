data "template_file" "deploy" {
  template = "ansible-playbook $${ansible_args} ${path.module}/deploy.yml -e '$${variables}'"

  vars = {
    ansible_args = local.ansible_args
    dependencies = jsonencode(var.dependencies)
    variables    = jsonencode(local.template_variables)
  }
}
