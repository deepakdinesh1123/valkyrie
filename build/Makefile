# Variables
HOST_UID := $(shell id -u)
HOST_GID := $(shell id -g)
HOST_USER := $(shell whoami)
HOST_GROUP := $(shell whoami)
DOCKER_IMAGE_UBUNTU := odin:ubuntu
DOCKER_IMAGE_ALPINE := odin:alpine
PODMAN_IMAGE := odin:latest
DOCKERFILE_UBUNTU := build/platforms/ubuntu.dockerfile
DOCKERFILE_ALPINE := build/platforms/alpine.dockerfile

# Phony Targets
.PHONY: all build-docker-image build-docker-image-alpine build-podman-image

# Default target
all: build-docker-image build-docker-image-alpine build-podman-image

# Build Docker image for Ubuntu
build-docker-image:
	docker build \
		-t $(DOCKER_IMAGE_UBUNTU) \
		--build-arg HOST_UID=$(HOST_UID) \
		--build-arg HOST_GID=$(HOST_GID) \
		--build-arg HOST_USER=$(HOST_USER) \
		--build-arg HOST_GROUP=$(HOST_GROUP) \
		-f $(DOCKERFILE_UBUNTU) .

# Build Docker image for Alpine
build-docker-image-alpine:
	docker build \
		-t $(DOCKER_IMAGE_ALPINE) \
		--build-arg HOST_UID=$(HOST_UID) \
		--build-arg HOST_GID=$(HOST_GID) \
		--build-arg HOST_USER=$(HOST_USER) \
		--build-arg HOST_GROUP=$(HOST_GROUP) \
		-f $(DOCKERFILE_ALPINE) .

# Build Podman image
build-podman-image:
	podman build \
		-t $(PODMAN_IMAGE) \
		--build-arg HOST_UID=$(HOST_UID) \
		--build-arg HOST_GID=$(HOST_GID) \
		--build-arg HOST_USER=$(HOST_USER) \
		--build-arg HOST_GROUP=$(HOST_GROUP) \
		-f $(DOCKERFILE_UBUNTU) .
