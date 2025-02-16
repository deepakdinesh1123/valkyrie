# Common variables
build_tool := "docker"
version := "0.0.1"

# Image names
nix_env_image := "odin_nix_env"
sandbox_image := "odin_sandbox"
store_image := "odin_store"
agent_image := "odin_agent"
sandbox_cs_image := "odin_sandbox_cs"

# Build worker images with specified base
build-worker base_os="alpine":
    {{build_tool}} build \
        -t {{nix_env_image}}:{{version}}-{{base_os}} \
        -f builds/images/execution/{{base_os}}.dockerfile .

# Build sandbox image
build-sandbox:
    {{build_tool}} build \
        -t {{sandbox_image}}:{{version}}-ubuntu \
        -f builds/images/sandbox/ubuntu.dockerfile .

# Build store image
build-store:
    {{build_tool}} build \
        -t {{store_image}}:{{version}} \
        -f builds/images/store/store.dockerfile .

# Build code-server image
build-code-server:
    {{build_tool}} build \
        -t {{sandbox_cs_image}}:{{version}}ubuntu \
        -f builds/images/sandbox/code_server.dockerfile .

# Build agent image
build-agent:
    {{build_tool}} build \
        -t {{agent_image}}:{{version}} \
        -f builds/images/packages/dockerfiles/agent.dockerfile .

# Docker-specific aliases
docker-odin-worker: (build-worker "alpine")
docker-odin-worker-ubuntu: (build-worker "ubuntu")
docker-odin-sandbox: build-sandbox
docker-odin-store: build-store
docker-code-server: build-code-server
docker-odin-agent: build-agent

# Podman-specific aliases
podman-odin-worker: (build-worker "alpine")
    just --set build_tool "podman" build-worker "alpine"
podman-odin-worker-ubuntu: (build-worker "ubuntu")
    just --set build_tool "podman" build-worker "ubuntu"
podman-odin-sandbox: build-sandbox
    just --set build_tool "podman" build-sandbox
podman-odin-store: build-store
    just --set build_tool "podman" build-store
podman-code-server: build-code-server
    just --set build_tool "podman" build-code-server
podman-odin-agent: build-agent
    just --set build_tool "podman" build-agent