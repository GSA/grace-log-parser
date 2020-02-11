resource "aws_cloudwatch_log_subscription_filter" "self" {
  depends_on      = [aws_lambda_permission.self]
  name            = "grace-log-parser-cloudwatch-subscription"
  log_group_name  = var.log_group_name
  filter_pattern  = ""
  destination_arn = aws_lambda_function.self.arn
}
