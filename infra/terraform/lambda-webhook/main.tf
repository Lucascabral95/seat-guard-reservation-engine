terraform {
  required_version = ">= 1.5.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_dir  = "${path.module}/../../../lambdas/payment-processor"
  output_path = "${path.module}/lambda.zip"
}

# --- 1. ROL DE LAMBDA (Actualizado con permisos SQS) ---
resource "aws_iam_role" "lambda_role" {
  name = "stripe-payment-processor-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect    = "Allow",
      Principal = { Service = "lambda.amazonaws.com" },
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_sqs" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole"
}

# --- 2. LAMBDA FUNCTION ---
resource "aws_lambda_function" "stripe_processor" {
  function_name = "stripe-payment-processor"
  role          = aws_iam_role.lambda_role.arn

  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256

  runtime = "nodejs20.x"
  handler = "index.handler"
  timeout = 30 

  environment {
    variables = {
      STRIPE_SECRET_KEY = var.stripe_secret_key
    }
  }
}

# --- 3. TRIGGER SQS ---
resource "aws_lambda_event_source_mapping" "sqs_trigger" {
  event_source_arn = "arn:aws:sqs:us-east-1:560765037562:payment-queue-reservation"
  function_name    = aws_lambda_function.stripe_processor.arn
  
  batch_size       = 1   
  enabled          = true
}