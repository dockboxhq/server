#!/bin/sh
echo "Warning: The following commands require root access"
sudo rm /var/run/docker.sock
sudo ln -s $HOME/Library/Containers/com.docker.docker/Data/docker.raw.sock /var/run/docker.sock