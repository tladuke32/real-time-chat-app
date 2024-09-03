provider "aws" {
  region = var.aws_region
}

resource "aws_instance" "demo_app" {
  ami           = var.ami_id
  instance_type = var.instance_type
  key_name      = var.key_name

  vpc_security_group_ids = [aws_security_group.allow_http.id]

  user_data = file("user_data.sh")  # Use the updated user_data.sh script

  tags = {
    Name = "RealTimeChatAppDemo"
  }
}

resource "aws_security_group" "allow_http" {
  name_prefix = "allow_http"

  ingress {
    description = "HTTP traffic"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
