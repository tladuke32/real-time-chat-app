#!/bin/bash
# Update package index
sudo yum update -y

# Install Docker and Git
sudo yum install docker git -y
sudo service docker start

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

#Set executable permissions for docker-compose and the ec2-user
sudo chmod +x /usr/local/bin/docker-compose
sudo usermod -a -G docker ec2-user

# Set up SSH for GitHub
mkdir -p /home/ec2-user/.ssh
echo "${SSH_PRIVATE_KEY}" > /home/ec2-user/.ssh/id_rsa
chmod 600 /home/ec2-user/.ssh/id_rsa
# Add GitHub to known hosts to prevent SSH verification prompt
ssh-keyscan github.com >> /home/ec2-user/.ssh/known_hosts

chown -R ec2-user:ec2-user /home/ec2-user/.ssh

# Switch to ec2-user and execute commands as this user
sudo -i -u ec2-user bash << EOF

# Clone the GitHub repository
cd /home/ec2-user
git clone -b rfranco/deployment-updates git@github.com:tladuke32/real-time-chat-app.git

# Create the .env file with necessary environment variables
cat <<EOL > /home/ec2-user/real-time-chat-app/.env
MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
MYSQL_DATABASE=${MYSQL_DATABASE}
MYSQL_USER=${MYSQL_USER}
MYSQL_PASSWORD=${MYSQL_PASSWORD}
JWT_SECRET=${JWT_SECRET}
REACT_APP_API_URL=${REACT_APP_API_URL}
EOL

EOF

# Change to the cloned directory and run Docker Compose
cd /home/ec2-user/real-time-chat-app
docker-compose up -d
