resource "aws_lambda_function" "self" {
  filename         = data.archive_file.lambda_zip_inline.output_path
  source_code_hash = data.archive_file.lambda_zip_inline.output_base64sha256
  function_name    = "grace-${var.env}-log-parser"
  role             = aws_iam_role.self.arn
  handler          = "lambda_function.lambda_handler"
  runtime          = "python3.8"
  timeout          = 500
  environment {
    variables = {
      TO_EMAIL   = var.recipients
      FROM_EMAIL = var.sender
      SUBJECT    = var.subject
      REGION     = var.region
    }
  }
}

data "archive_file" "lambda_zip_inline" {
  type        = "zip"
  output_path = "/tmp/lambda_zip_inline.zip"
  source_file = "${path.module}/handler/lambda_function.py"
}

# Lambda subscription to CloudWatch log group
resource "aws_lambda_permission" "self" {
  statement_id  = "grace-${var.env}-log-parser-allow-cloudwatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.self.function_name
  principal     = "logs.${var.region}.amazonaws.com"
  source_arn    = var.source_arn
}
