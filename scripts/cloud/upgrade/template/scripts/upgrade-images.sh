#!/bin/bash

# This script is used to upgrade the images for sealos cloud components
kubectl set image deployment/applaunchpad-frontend -n applaunchpad-frontend applaunchpad-frontend=ghcr.io/nebstudio/sealos-applaunchpad-frontend:v5.0.0
