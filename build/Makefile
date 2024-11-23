# Phony Targets
.PHONY: all build-docker-image build-docker-image-alpine build-podman-image

# Default target
all: build-docker-image build-docker-image-ubuntu build-podman-image build-podman-image-ubuntu

# Build Docker images
build-docker-image:
	docker build \
		-t odin:alpine \
		-f build/platforms/alpine/Dockerfile .

build-docker-image-ubuntu:
	docker build \
		-t odin:ubuntu \
		-f build/platforms/ubuntu/Dockerfile .

build-podman-image-ubuntu:
	podman build \
		-t odin:ubuntu \
		-f build/platforms/ubuntu/Containerfile .

build-podman-image:
	podman build \
		-t odin:alpine \
		-f build/platforms/alpine/Containerfile .

build-sandbox-image:
	docker build \
	-t odin:sandbox \
	-f build/platforms/sandbox/base.dockerfile
