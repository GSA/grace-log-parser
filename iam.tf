resource "aws_iam_role" "self" {
  name = "grace-${var.env}-log-parser-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}


resource "aws_iam_policy" "self" {
  name        = "grace-${var.env}-log-parser"
  path        = "/"
  description = "IAM policy for logging from a lambda"

  #tfsec:ignore:AWS099
  policy = <<END_OF_POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "graceLogParser",
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents",
                "ses:SendEmail"
            ],
            "Resource": "*"
        }
    ]
}
END_OF_POLICY
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.self.name
  policy_arn = aws_iam_policy.self.arn
}
