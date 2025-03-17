# Common variables
build_tool := "docker"
version := "0.0.1"

# Image names
odin_exec_image := "odin_execution"
sandbox_image := "odin_sandbox"
store_image := "odin_store"
agent_image := "odin_agent"
sandbox_cs_image := "odin_sandbox_cs"

# Build execution images with specified base
build-execution base_os="ubuntu":
    {{build_tool}} build \
        --build-arg NIX_CACHE_PUBLIC_KEY=$(cat "$HOME/.odin/store/setup/cache-pub-key.pem") \
        -t {{odin_exec_image}}:{{version}}-{{base_os}} \
        -f builds/images/execution/{{base_os}}.dockerfile .

# Build sandbox image
build-sandbox:
    {{build_tool}} build \
        --build-arg NIX_CACHE_PUBLIC_KEY=$(cat "$HOME/.odin/store/setup/cache-pub-key.pem") \
        -t {{sandbox_image}}:{{version}}-ubuntu \
        -f builds/images/sandbox/ubuntu.dockerfile .

# Build code-server image
build-code-server:
    {{build_tool}} build \
        -t {{sandbox_cs_image}}:{{version}}-ubuntu \
        -f builds/images/sandbox/code_server.dockerfile .

# Build agent image
build-agent:
    {{build_tool}} build \
        -t {{agent_image}}:{{version}} \
        -f builds/images/packages/agent.dockerfile .

# Docker-specific aliases
docker-odin-worker: (build-execution "alpine")
docker-odin-worker-ubuntu: (build-execution "ubuntu")
docker-odin-sandbox: build-sandbox
docker-odin-store: build-store
docker-code-server: build-code-server
docker-odin-agent: build-agent

# Podman-specific aliases
podman-odin-worker: (build-execution "alpine")
    just --set build_tool "podman" build-execution "alpine"
podman-odin-worker-ubuntu: (build-execution "ubuntu")
    just --set build_tool "podman" build-execution "ubuntu"
podman-odin-sandbox: build-sandbox
    just --set build_tool "podman" build-sandbox
podman-odin-store: build-store
    just --set build_tool "podman" build-store
podman-code-server: build-code-server
    just --set build_tool "podman" build-code-server
podman-odin-agent: build-agent
    just --set build_tool "podman" build-agent


build-setup-store:
    {{build_tool}} build -f builds/images/store/build.dockerfile -t odin_setup_store . 
setup-store:
	{{build_tool}} run --rm -v $HOME/.odin/store/nix:/nix -v $HOME/.odin/store/setup:/tmp/setup odin_setup_store

# Build store image
build-store:
	export NIX_CHANNELS_ENVIRONMENT=$(grep -oP '(?<=^NIX_CHANNELS_ENVIRONMENT=).*' "$HOME/.odin/store/setup/.env"); \
	export NIX_USER_ENVIRONMENT=$(grep -oP '(?<=^NIX_USER_ENVIRONMENT=).*' "$HOME/.odin/store/setup/.env"); \
	{{build_tool}} build \
		--build-arg NIX_CHANNELS_ENVIRONMENT="$NIX_CHANNELS_ENVIRONMENT" \
		--build-arg NIX_USER_ENVIRONMENT="$NIX_USER_ENVIRONMENT" \
		-t {{store_image}}:{{version}} \
		-f builds/images/store/store.dockerfile .
run-store:
    #!/bin/bash
    # Extract environment variables from the .env file
    NIX_CHANNELS_ENVIRONMENT=$(grep -oP '(?<=^NIX_CHANNELS_ENVIRONMENT=).*' "$HOME/.odin/store/setup/.env")
    NIX_USER_ENVIRONMENT=$(grep -oP '(?<=^NIX_USER_ENVIRONMENT=).*' "$HOME/.odin/store/setup/.env")

    # Execute the build tool with the appropriate parameters
    {{build_tool}} run --rm -it \
        -p 5000:5000 \
        -v "$HOME/.odin/store/nix:/nix" \
        -v "$HOME/.odin/store/setup:/tmp/setup" \
        -e NIX_CHANNELS_ENVIRONMENT="$NIX_CHANNELS_ENVIRONMENT" \
        -e NIX_USER_ENVIRONMENT="$NIX_USER_ENVIRONMENT" \
        {{store_image}}:{{version}}
