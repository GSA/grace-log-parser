variable "env" {
  type        = string
  description = "(optional) The environment in which the script is running (development | test | production)"
  default     = "development"
}

variable "recipients" {
  type        = string
  description = "(required) comma delimited list of AWS SES eMail recipients"
}

variable "sender" {
  type        = string
  description = "(required) eMail address of sender for AWS SES"
}

variable "region" {
  description = "(optional) AWS region to deploy lambda function."
  default     = "us-east-1"
}

variable "source_arn" {
  description = "(required) Source ARN of Cloudtrail Log Group"
}

variable "source_file" {
  type        = string
  description = "(optional) full or relative path to zipped binary of lambda handler"
  default     = "../release/grace-log-parser.zip"
}

variable "log_group_name" {
  description = "(required) Cloudtrail Log Group Name"
}
