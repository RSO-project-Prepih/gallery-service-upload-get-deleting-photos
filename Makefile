# Makefile

# Binary name
BINARY_NAME=server

# Docker image name
IMAGE_NAME=jabok123458/rso_prepih
IMAGE_TAG=latest

# Build the Go binary
build:
    CGO_ENABLED=0 GOOS=linux go build -o ${BINARY_NAME} main.go

# Build the Docker image
docker-build:
    docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .

# Push the Docker image
docker-push:
    echo "$$DOCKERHUB_ACCESS_TOKEN" | docker login --username $$DOCKERHUB_USERNAME --password-stdin
    docker push ${IMAGE_NAME}:${IMAGE_TAG}

# Clean up
clean:
    go clean
    rm -f ${BINARY_NAME}
    docker rmi ${IMAGE_NAME}:${IMAGE_TAG}

# Run the tests
test:
    go test ./...

# Default command
all: build

.PHONY: build docker-build docker-push clean test all