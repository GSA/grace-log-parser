output "grace-log-parser-name" {
  value       = aws_lambda_function.self.function_name
  description = "Function name of grace-log-parser Lambda function"
}

output "grace-log-parser-arn" {
  value       = aws_lambda_function.self.arn
  description = "ARN of grace-log-parser Lambda function"
}
