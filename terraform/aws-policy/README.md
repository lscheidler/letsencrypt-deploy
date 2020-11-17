# letsencrypt-deploy aws policy

Terraform module for creating an aws iam policy to allow access for letsencrypt-deploy to dynamodb table

## Usage

You can use terraform (>= 0.13.5) to deploy the lambda function:

```
module "letsencrypt" {
  source = "github.com/lscheidler/letsencrypt-deploy/terraform/aws-policy?ref=main"

  aws_dynamodb_table_arn = "arn:aws:dynamodb:<region>:<account-id>:table/LetsencryptCA"
}
```

It is going to configure
- iam role and policy for required permissions

## Argument Reference

| Name                                    | Required  | Default                                   | Description                                     |
|-----------------------------------------|-----------|-------------------------------------------|-------------------------------------------------|
| `aws_dynamodb_table_arn`                | 🗹         |                                           | DynamoDB table arn                              |
| `aws_iam_role_name`                     | 🗷         | `letsencrypt`                             |                                                 |
| `aws_iam_instance_profile_name`         | 🗷         | `letsencrypt`                             |                                                 |
| `aws_iam_policy_name`                   | 🗷         | `letsencrypt`                             |                                                 |
| `aws_iam_policy_description`            | 🗷         | `letsencrypt client permissions`          |                                                 |
| `aws_iam_policy_path`                   | 🗷         | `/`                                       |                                                 |
