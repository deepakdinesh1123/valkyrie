# Build Alpine-based Docker image
docker-odin-worker:
    docker build \
        -t odin_nix_env:alpine \
        -f build/platforms/alpine/Dockerfile .

# Build Ubuntu-based Docker image
docker-odin-worker-ubuntu:
    docker build \
        -t odin_nix_env:ubuntu \
        -f build/platforms/ubuntu/Dockerfile .

# Build Ubuntu-based Podman image
podman-odin-worker-ubuntu:
    podman build \
        -t odin_nix_env:ubuntu \
        -f build/platforms/ubuntu/Containerfile .

# Build Alpine-based Podman image
podman-odin-worker:
    podman build \
        -t odin_nix_env:alpine \
        -f build/platforms/alpine/Containerfile .