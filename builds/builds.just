# Build Alpine-based Docker image
docker-odin-worker:
    docker build \
        -t odin_nix_env:alpine \
        -f builds/images/execution/alpine.dockerfile .

# Build Ubuntu-based Docker image
docker-odin-worker-ubuntu:
    docker build \
        -t odin_nix_env:ubuntu \
        -f builds/images/execution/ubuntu.dockerfile .

# Build Ubuntu-based Podman image
podman-odin-worker-ubuntu:
    podman build \
        -t odin_nix_env:ubuntu \
        -f builds/images/execution/ubuntu.dockerfile .

# Build Alpine-based Podman image
podman-odin-worker:
    podman build \
        -t odin_nix_env:alpine \
        -f builds/images/execution/alpine.dockerfile .

docker-odin-sandbox:
    docker build \
        -t odin_sandbox:ubuntu \
        -f builds/images/sandbox/ubuntu.dockerfile .

podman-odin-sandbox:
    podman build \
        -t odin_sandbox:ubuntu \
        -f builds/images/sandbox/ubuntu.dockerfile .

docker-odin-store:
    docker build \
        -t odin_store:ubuntu \
        -f builds/images/store/store.dockerfile .

docker-code-server:
    docker build \
        -t odin_store:ubuntu \
        -f builds/images/sandbox/code_server.dockerfile .