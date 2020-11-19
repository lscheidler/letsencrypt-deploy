################################################################################
# aws_iam_role.letsencrypt
################################################################################
variable "aws_iam_role_name" {
  default = "letsencrypt"
}

################################################################################
# aws_iam_instance_profile.letsencrypt
################################################################################
variable "aws_iam_instance_profile_name" {
  default = "letsencrypt"
}

################################################################################
# aws_iam_policy.letsencrypt
################################################################################
variable "aws_iam_policy_name" {
  default = "letsencrypt"
}

variable "aws_iam_policy_description" {
  default = "letsencrypt client permissions"
}

variable "aws_iam_policy_path" {
  default = "/"
}

################################################################################
# aws_iam_policy_document.letsencrypt
################################################################################
variable "aws_dynamodb_table_arn" {
  type = string
}

variable "sns_topic" {
  type    = list
  default = []
}
