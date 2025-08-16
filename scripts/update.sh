#!/bin/sh
git fetch
git pull

sudo docker compose down
sudo docker compose rm -s -f
sudo docker compose build
sudo docker compose up -d