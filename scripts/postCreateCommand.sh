#!/bin/bash
# This script is executed after the creation of a new project.

go install github.com/goreleaser/goreleaser/v2@latest

curl -sS https://pkg.authzed.com/apt/gpg.key | sudo gpg --dearmor --yes -o /etc/apt/keyrings/authzed.gpg
echo "deb [signed-by=/etc/apt/keyrings/authzed.gpg] https://pkg.authzed.com/apt/ * *"  | sudo tee /etc/apt/sources.list.d/authzed.list
sudo chmod 644 /etc/apt/sources.list.d/authzed.list  # helps tools such as command-not-found to work correctly

sudo apt update
sudo apt install -y zed