package main

import (
	"github.com/GSA/grace-log-parser/lambda/app"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	a := app.New()
	lambda.Start(a.Run)
}
