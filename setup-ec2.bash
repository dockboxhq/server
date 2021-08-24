#!/bin/bash
sudo apt-get update
sudo apt-get upgrade

echo "Setting up docker..."
sudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

    
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

sudo echo \
  "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get -y install docker-ce docker-ce-cli containerd.io git binutils

cd $HOME && git clone https://github.com/aws/efs-utils

cd efs-utils && ./build-deb.sh

sudo apt-get -y install ./build/amazon-efs-utils*deb

sudo mkdir -p /efs/dockbox
echo "Mounting efs..."
sudo mount -t efs -o tls,rw fs-86a8f032:/ /efs/dockbox

sudo groupadd docker
sudo usermod -aG docker ubuntu

sudo chgrp -R docker /efs/dockbox
sudo chmod go+rw /efs/dockbox
sudo chgrp -R docker /usr/bin/git
sudo chmod g+x /usr/bin/git

sudo mkdir -p /etc/systemd/system/docker.service.d && touch /etc/systemd/system/docker.service.d/override.conf

sudo echo "[Service]
ExecStart=
ExecStart=/usr/bin/dockerd" > /etc/systemd/system/docker.service.d/override.conf

sudo touch /etc/docker/daemon.json

sudo echo '{"hosts": ["tcp://0.0.0.0:2375", "unix:///var/run/docker.sock"]}' > /etc/docker/daemon.json


echo "Starting services"

sudo systemctl daemon-reload

sudo systemctl restart docker.service

sudo systemctl status docker.service

docker pull ubuntu:latest

sudo touch /etc/systemd/system/api.service

sudo echo "[Unit]
Description=API Service
Wants=network.target
After=network.target

[Service]
Type=simple
DynamicUser=yes
Group=docker
ExecStart=/usr/local/bin/api
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/api.service

echo "Starting Reboot in 5 seconds"
sleep 5

sudo reboot