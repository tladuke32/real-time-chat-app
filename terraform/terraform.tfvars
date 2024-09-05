aws_region       = "us-east-2"
ami_id           = "ami-0490fddec0cbeb88b" # Example AMI ID for Amazon Linux 2
instance_type    = "t2.micro"
key_name         = "EC2KP"
private_key_path = "./kp/EC2KP.pem"
vpc_cidr         = "10.0.0.0/16"
subnet_cidr      = "10.0.1.0/24"