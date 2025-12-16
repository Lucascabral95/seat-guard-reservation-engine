# Test para implementacion de Github Actions
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}
# terraform {
#   backend "s3" {
#     bucket = "seatguard-terraform-state-lucas" 
#     key    = "fargate/terraform.tfstate"       
#     region = "us-east-1"
#   }
# }
# Test para implementacion de Github Actions

provider "aws" {
  region = "us-east-1"
}

# locals {
#   image_auth    = "560765037562.dkr.ecr.us-east-1.amazonaws.com/auth-service:latest"
#   image_booking = "560765037562.dkr.ecr.us-east-1.amazonaws.com/booking-service:latest"
#   app_name      = "monorepo-prod"
# }

# Test para implementacion de Github Actions
locals {
  image_auth    = var.names_images_ecr["auth-service-image"]
  image_booking = var.names_images_ecr["booking-service-image"]
  app_name      = var.names_images_ecr["app_name"]
}
# Test para implementacion de Github Actions

# 1. LOGS
resource "aws_cloudwatch_log_group" "auth_logs" {
  name              = "/ecs/auth-service"
  retention_in_days = 7
}

resource "aws_cloudwatch_log_group" "booking_logs" {
  name              = "/ecs/booking-service"
  retention_in_days = 7
}

# 2. RED (VPC) - MODIFICADO PARA SEGURIDAD DE COSTOS
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  name   = "${local.app_name}-vpc"
  cidr   = "10.0.0.0/16"

  azs            = ["us-east-1a", "us-east-1b"]
  public_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
  
  # --- ESTAS LINEAS SON CLAVE PARA TU BOLSILLO ---
  enable_nat_gateway = false
  single_nat_gateway = false
  # -----------------------------------------------

  enable_dns_hostnames = true
  enable_dns_support   = true
}

# 3. BALANCEADOR (ALB)
resource "aws_lb" "main" {
  name               = "${local.app_name}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb_sg.id]
  subnets            = module.vpc.public_subnets
}

resource "aws_security_group" "alb_sg" {
  name   = "${local.app_name}-alb-sg"
  vpc_id = module.vpc.vpc_id

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 80
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    protocol    = "tcp"
    from_port   = 8080
    to_port     = 8080
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# 4. ECS CLUSTER
resource "aws_ecs_cluster" "main" {
  name = "${local.app_name}-cluster"
}

resource "aws_ecs_cluster_capacity_providers" "main" {
  cluster_name = aws_ecs_cluster.main.name
  capacity_providers = ["FARGATE_SPOT", "FARGATE"]
  
  default_capacity_provider_strategy {
    base              = 1
    weight            = 100
    capacity_provider = "FARGATE_SPOT"
  }
}

# 5. TAREAS Y SERVICIOS

# --- TAREA 1: AUTH SERVICE ---
resource "aws_ecs_task_definition" "auth" {
  family                   = "auth-service-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 512
  memory                   = 1024
  execution_role_arn       = aws_iam_role.ecs_execution_role.arn

  container_definitions = jsonencode([{
    name      = "auth-container"
    image     = local.image_auth
    essential = true
    portMappings = [{
      containerPort = 3000
      hostPort      = 3000
    }]
  
    environment = [
      for key, value in var.auth_service_envs : {
        name  = key 
        value = value
      }
    ]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        awslogs-group         = aws_cloudwatch_log_group.auth_logs.name
        awslogs-region        = "us-east-1"
        awslogs-stream-prefix = "ecs"
      }
    }
  }])
}

resource "aws_ecs_service" "auth" {
  name            = "auth-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.auth.arn
  desired_count   = 1

  network_configuration {
    subnets          = module.vpc.public_subnets
    security_groups  = [aws_security_group.ecs_tasks_sg.id]
    assign_public_ip = true # Esto permite salir a internet sin NAT
  }

  capacity_provider_strategy {
    capacity_provider = "FARGATE_SPOT"
    weight            = 100
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.auth.arn
    container_name   = "auth-container"
    container_port   = 3000
  }
}

# --- TAREA 2: BOOKING SERVICE ---
resource "aws_ecs_task_definition" "booking" {
  family                   = "booking-service-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_execution_role.arn

  container_definitions = jsonencode([{
    name      = "booking-container"
    image     = local.image_booking
    essential = true
    portMappings = [{
      containerPort = 4000
      hostPort      = 4000
    }]

    environment = [
      for key, value in var.booking_service_envs : {
        name  = key
        value = value
      }
    ]
    
    logConfiguration = {
      logDriver = "awslogs"
      options = {
        awslogs-group         = aws_cloudwatch_log_group.booking_logs.name
        awslogs-region        = "us-east-1"
        awslogs-stream-prefix = "ecs"
      }
    }
  }])
}

resource "aws_ecs_service" "booking" {
  name            = "booking-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.booking.arn
  desired_count   = 1

  network_configuration {
    subnets          = module.vpc.public_subnets
    security_groups  = [aws_security_group.ecs_tasks_sg.id]
    assign_public_ip = true # Esto permite salir a internet sin NAT
  }

  capacity_provider_strategy {
    capacity_provider = "FARGATE_SPOT"
    weight            = 100
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.booking.arn
    container_name   = "booking-container"
    container_port   = 4000
  }
}

# 6. ROUTING

# Target Group para Auth
resource "aws_lb_target_group" "auth" {
  name_prefix = "auth-"  
  port        = 3000
  protocol    = "HTTP"
  vpc_id      = module.vpc.vpc_id
  target_type = "ip"
  
  health_check {
    path    = "/health"
    matcher = "200-299"
  }

  lifecycle { create_before_destroy = true }
}

# Target Group para Booking
resource "aws_lb_target_group" "booking" {
  name_prefix = "book-" 
  port        = 4000
  protocol    = "HTTP"
  vpc_id      = module.vpc.vpc_id
  target_type = "ip"
  
  health_check {
    path    = "/api/v1/events" 
    matcher = "200-299"
  }

  lifecycle { create_before_destroy = true }
}

# Listeners
resource "aws_lb_listener" "http_auth" {
  load_balancer_arn = aws_lb.main.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.auth.arn
  }
}

resource "aws_lb_listener" "http_booking" {
  load_balancer_arn = aws_lb.main.arn
  port              = 8080
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.booking.arn
  }
}

# 7. SEGURIDAD Y IAM
resource "aws_security_group" "ecs_tasks_sg" {
  name   = "${local.app_name}-tasks-sg"
  vpc_id = module.vpc.vpc_id

  ingress {
    protocol        = "tcp"
    from_port       = 0
    to_port         = 65535
    security_groups = [aws_security_group.alb_sg.id]
  }

  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_iam_role" "ecs_execution_role" {
  name = "${local.app_name}-exec-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_execution_role_policy" {
  role       = aws_iam_role.ecs_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

output "url_servicio" {
  value = "http://${aws_lb.main.dns_name}"
}
