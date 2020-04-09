data "aws_iam_account_alias" "current" {}

resource "aws_lambda_function" "self" {
  filename         = var.source_file
  source_code_hash = filesha256(var.source_file)
  function_name    = "grace-${var.env}-log-parser"
  role             = aws_iam_role.self.arn
  handler          = "grace-log-parser"
  runtime          = "go1.x"
  timeout          = 500
  environment {
    variables = {
      TO_EMAIL      = var.recipients
      FROM_EMAIL    = var.sender
      ACCOUNT_ALIAS = data.aws_iam_account_alias.current.account_alias
      REGION        = var.region
    }
  }
}

# Lambda subscription to CloudWatch log group
resource "aws_lambda_permission" "self" {
  statement_id  = "grace-${var.env}-log-parser-allow-cloudwatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.self.function_name
  principal     = "logs.${var.region}.amazonaws.com"
  source_arn    = var.source_arn
}
