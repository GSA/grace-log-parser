# GRACE Log Parser

Lambda function to parse CloudWatch Alarms and send legible email to designated recipients

## Repository contents

[handler](handler/) - Python function to parse logs and send eMail
[terraform](https://github.com/GSA/grace-log-parser) - Terraform module to install the Lambda function and set IAM role, policy, triggers and environment variables

## Inputs ##

|     Name     | Description |
| ------------ | ----------- |
| env | (optional) The environment in which the script is running (development | test | production) |
| recipients | (required) comma delimited list of AWS SES eMail recipients |
| sender | (required) eMail address of sender for AWS SES |
| subject | (optional) Subject Header of  Email sent notifications |
| region | (optional) AWS region to deploy lambda function. |
| source_arn | (required) Source ARN of Cloudtrail Log Group |
| log_group_name | (required) Cloudtrail Log Group Name |

## Outputs ##

|     Name     | Description |
| ------------ | ----------- |
| grace-log-parser-name | Function name of grace-log-parser Lambda function |
| grace-log-parser-arn | ARN of grace-log-parser Lambda function |

## Public domain

This project is in the worldwide [public domain](LICENSE.md). As stated in [CONTRIBUTING](CONTRIBUTING.md):

> This project is in the public domain within the United States, and copyright and related rights in the work worldwide are waived through the [CC0 1.0 Universal public domain dedication](https://creativecommons.org/publicdomain/zero/1.0/).
>
> All contributions to this project will be released under the CC0 dedication. By submitting a pull request, you are agreeing to comply with this waiver of copyright interest.
