terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }
}

provider "aws" {
  region                   = "eu-central-1"
  shared_credentials_files = ["~/.aws/credentials"]
  profile                  = "MohamedASalah"
}

// Creating policy
data "aws_iam_policy_document" "AWSLambdaTrustPolicy" {
  statement {
    actions = ["sts:AssumeRole"]
    effect  = "Allow"
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

// Creating role
resource "aws_iam_role" "terraform_function_role" {
  name               = "terraform_function_role"
  assume_role_policy = data.aws_iam_policy_document.AWSLambdaTrustPolicy.json
}

resource "aws_iam_role_policy_attachment" "terraform_lambda_policy" {
  role       = aws_iam_role.terraform_function_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

// Build main.go first
resource "null_resource" "build-canvas-lambda" {
  provisioner "local-exec" {
    command = "cd .. && make build-lambda"
  }
  triggers = {
    always_run = "${timestamp()}"
  }
}

// archive the build
data "archive_file" "zip-canvas-lambda" {
  depends_on  = [null_resource.build-canvas-lambda]
  type        = "zip"
  source_file = "../bin/main"
  output_path = "../bin/main_archived.zip"
}

// Creating lambda function containing the archived file
resource "aws_lambda_function" "canvas-lambda" {
  depends_on       = [null_resource.build-canvas-lambda]
  filename         = "./../bin_archived/main_archived.zip"
  source_code_hash = data.archive_file.zip-canvas-lambda.output_base64sha256
  function_name    = "canvas-lambda"
  handler          = "bin/main"
  role             = aws_iam_role.terraform_function_role.arn
  runtime          = "go1.x"
}

resource "aws_api_gateway_rest_api" "rest-api-gateway" {
  name = "rest-api-gateway"
}

resource "aws_api_gateway_resource" "rest-api-resource-canvas" {
  path_part   = "canvas"
  parent_id   = aws_api_gateway_rest_api.rest-api-gateway.root_resource_id
  rest_api_id = aws_api_gateway_rest_api.rest-api-gateway.id
}

resource "aws_api_gateway_resource" "rest-api-resource-canvas-proxy" {
  path_part   = "{proxy+}"
  parent_id   = aws_api_gateway_resource.rest-api-resource-canvas.id
  rest_api_id = aws_api_gateway_rest_api.rest-api-gateway.id
}

resource "aws_api_gateway_method" "rest-api-canvas-method" {
  rest_api_id   = aws_api_gateway_rest_api.rest-api-gateway.id
  resource_id   = aws_api_gateway_resource.rest-api-resource-canvas-proxy.id
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "lambda_integration" {
  rest_api_id             = aws_api_gateway_rest_api.rest-api-gateway.id
  resource_id             = aws_api_gateway_resource.rest-api-resource-canvas-proxy.id
  http_method             = aws_api_gateway_method.rest-api-canvas-method.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.canvas-lambda.invoke_arn
}

resource "aws_api_gateway_deployment" "api_deployment" {
  rest_api_id = aws_api_gateway_rest_api.rest-api-gateway.id
  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_integration.lambda_integration,
    ]))
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.canvas-lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.rest-api-gateway.execution_arn}/*/*${aws_api_gateway_resource.rest-api-resource-canvas.path}/*"
}

resource "aws_api_gateway_stage" "rest-api-services-stage" {
  deployment_id = aws_api_gateway_deployment.api_deployment.id
  rest_api_id   = aws_api_gateway_rest_api.rest-api-gateway.id
  stage_name    = "services"
}
