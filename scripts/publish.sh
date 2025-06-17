#!/bin/bash

# Configuration
IMAGE_NAME="rodrwan/secretly"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Build the image
echo "Building Docker image..."
docker build -t ${IMAGE_NAME}:${VERSION} -t ${IMAGE_NAME}:latest .

# Login to Docker Hub
echo "Logging in to Docker Hub..."
docker login

# Push the images
echo "Pushing images to Docker Hub..."
docker push ${IMAGE_NAME}:${VERSION}
docker push ${IMAGE_NAME}:latest

echo "Done! Image published as ${IMAGE_NAME}:${VERSION} and ${IMAGE_NAME}:latest"