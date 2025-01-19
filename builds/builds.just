# Build Alpine-based Docker image
docker-odin-worker:
    docker build \
        -t odin_nix_env:alpine \
        -f builds/platforms/alpine/Dockerfile .

# Build Ubuntu-based Docker image
docker-odin-worker-ubuntu:
    docker build \
        -t odin_nix_env:ubuntu \
        -f builds/platforms/ubuntu/Dockerfile .

# Build Ubuntu-based Podman image
podman-odin-worker-ubuntu:
    podman build \
        -t odin_nix_env:ubuntu \
        -f builds/platforms/ubuntu/Containerfile .

# Build Alpine-based Podman image
podman-odin-worker:
    podman build \
        -t odin_nix_env:alpine \
        -f builds/platforms/alpine/Containerfile .

docker-odin-sandbox:
    docker build \
        -t odin_sandbox:ubuntu \
        -f builds/platforms/sandbox/ubuntu/Dockerfile .

podman-odin-sandbox:
    podman build \
        -t odin_sandbox:alpine \
        -f builds/platforms/sandbox/alpine/Containerfile .