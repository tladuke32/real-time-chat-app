#!/bin/bash
# Update package index
sudo yum update -y

# Install Docker
sudo amazon-linux-extras install docker -y
sudo service docker start
sudo usermod -a -G docker ec2-user

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Install Git
sudo yum install git -y

# Clone the GitHub repository
cd /home/ec2-user
git clone git@github.com:tladuke32/real-time-chat-app.git

# Change to the cloned directory and run Docker Compose
cd real-time-chat-app
sudo docker-compose up -d