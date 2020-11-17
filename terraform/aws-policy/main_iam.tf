resource "aws_iam_role" "letsencrypt" {
  name = var.aws_iam_role_name

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "letsencrypt" {
  name = var.aws_iam_instance_profile_name
  role = aws_iam_role.letsencrypt.name
}

resource "aws_iam_role_policy_attachment" "letsencrypt" {
  role       = aws_iam_role.letsencrypt.name
  policy_arn = aws_iam_policy.letsencrypt.arn
}

resource "aws_iam_policy" "letsencrypt" {
  name        = var.aws_iam_policy_name
  description = var.aws_iam_policy_description
  path        = var.aws_iam_policy_path
  policy      = data.aws_iam_policy_document.letsencrypt.json
}

data "aws_iam_policy_document" "letsencrypt" {
  statement {
    effect = "Allow"

    actions = [
      "dynamodb:GetItem",
    ]

    resources = [
      var.aws_dynamodb_table_arn,
    ]
  }
}
