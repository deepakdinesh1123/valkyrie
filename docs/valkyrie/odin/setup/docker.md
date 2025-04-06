# Docker Mode

Docker mode uses Docker containers as isolated execution environments for running code, providing better security and resource isolation compared to native mode.

## Prerequisites

- Docker installed and running on your system
- Sufficient permissions to manage Docker containers

## Image Setup

Build the required Docker images using the provided Makefile targets:

```bash
# Build Alpine-based image
make build-docker-image

# Build Ubuntu-based image
make build-docker-image-ubuntu
```

## Configuration

Set the following environment variables:

```bash
ODIN_MODE=docker
ODIN_DOCKER_IMAGE=odin:alpine  # or odin:ubuntu
ODIN_DOCKER_NETWORK=<network_name>  # Optional
ODIN_WORKER_COUNT=<number_of_workers>
```

## Security Features

- Process isolation through containerization
- Resource limiting via Docker constraints
- Network isolation options
- Filesystem isolation

## Resource Management

Configure container resources:
```bash
ODIN_DOCKER_MEMORY_LIMIT=<memory_limit>
ODIN_DOCKER_CPU_LIMIT=<cpu_limit>
ODIN_DOCKER_TIMEOUT=<timeout_seconds>
```

## Networking

1. Create a dedicated Docker network (recommended):
   ```bash
   docker network create odin-network
   ```
2. Configure network settings in your environment variables

## Troubleshooting

1. Check container logs:
   ```bash
   docker logs <container_id>
   ```
2. Verify container status:
   ```bash
   docker ps -a | grep odin
   ```
3. Inspect network connectivity:
   ```bash
   docker network inspect odin-network
   ```
