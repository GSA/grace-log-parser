# GRACE Log Parser [![License](https://img.shields.io/badge/license-CC0-blue)](LICENSE.md) [![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg)](https://godoc.org/github.com/GSA/grace-log-parser/aws) [![CircleCI](https://circleci.com/gh/GSA/grace-log-parser.svg?style=shield)](https://circleci.com/gh/GSA/grace-log-parser) [![Go Report Card](https://goreportcard.com/badge/github.com/GSA/grace-log-parser)](https://goreportcard.com/report/github.com/GSA/grace-log-parser)

Lambda function that parses Cloudwatch events and performs different actions based on enabled modules and sub-modules

This function is executed on a configurable interval (defaults to 5m) stores the last read event ID in secretsmanager for use on the next execution. Then iterates through each event since the last known event ID executing all loaded modules and sub-modules.

Each module is responsible for carrying out its area of responsibility. It will be executed in a separate go-routine and given each full-parsed event for processing. Modules can decline to process a particular event by returning the `modules.NotApplicableErr` to indicate no action was taken.

## Repository contents

[lambda](lambda/) - Golang Lambda function handler to parse logs
[terraform](https://github.com/GSA/grace-log-parser) - Terraform module to install the Lambda function and set IAM role, policy, triggers and environment variables

## Usage

### Download (recommended)

Download the zip compressed executable (Note: Replace v0.0.2 with desired release):

```
mkdir -p release
curl -L https://github.com/GSA/grace-log-parser/releases/download/v0.0.2/grace-log-parser.zip -o release/grace-log-parser.zip
```

### Compile

Alternatively, you can compile the Lambda function handler yourself:

```
cd handler
GOOS=linux GOARCH=amd64 go build -o ../release/grace-log-parser -v
zip -j ../release/grace-log-parser.zip ../release/grace-log-parser
```

### Add Module

Add the module to your terraform project. Ensure path to `source_file` matches
where you downloaded the zip file. Replace v0.0.2 with desired release. Example below:

```
module "grace-log-parser" {
  source        = "github.com:GSA/grace-log-parser?ref=v0.0.2"
  source_file   = "../release/grace-log-parser.zip"
  env           = "development"
  sender        = "validated-sender@email.com"
  recipients    = "recipient@email.com,other-recipient@email.com"
  source_arn    = module.logging.cloudtrail_log_group_arn
  log_group_arn = module.logging.cloudtrail_log_group_name
}
```

## Inputs ##

|     Name     | Description |
| ------------ | ----------- |
| env | (optional) The environment in which the script is running (development | test | production) |
| recipients | (required) comma delimited list of AWS SES eMail recipients |
| sender | (required) eMail address of sender for AWS SES |
| region | (optional) AWS region to deploy lambda function. |
| source_arn | (required) Source ARN of Cloudtrail Log Group |
| source_file | "(optional) full or relative path to zipped binary of lambda handler" |
| log_group_name | (required) Cloudtrail Log Group Name |

## Outputs ##

|     Name     | Description |
| ------------ | ----------- |
| grace-log-parser-name | Function name of grace-log-parser Lambda function |
| grace-log-parser-arn | ARN of grace-log-parser Lambda function |


## Environment Variables ##

| Name | Description |
| ---- | ----------- |
| DISABLED_MODULES | A comma delimited list of modules and sub-modules to disable |


## Modules and Sub-modules ##

| Name | Type | Description |
| ---- | ---- | ----------- |
| email | module | Enables email notification and all sub-modules |
| email.login | sub-module | Sends email on console login |


## Public domain

This project is in the worldwide [public domain](LICENSE.md). As stated in [CONTRIBUTING](CONTRIBUTING.md):

> This project is in the public domain within the United States, and copyright and related rights in the work worldwide are waived through the [CC0 1.0 Universal public domain dedication](https://creativecommons.org/publicdomain/zero/1.0/).
>
> All contributions to this project will be released under the CC0 dedication. By submitting a pull request, you are agreeing to comply with this waiver of copyright interest.
