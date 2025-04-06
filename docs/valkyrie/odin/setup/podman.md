# Podman Mode

Podman mode utilizes Podman containers for code execution, offering a daemonless container solution with similar isolation benefits to Docker mode.

## Platform Support

> **Important Note**: Podman is not available for macOS systems natively. This is because Podman relies heavily on Linux kernel features and namespace isolation that are not present in macOS. While Docker for Mac uses a Linux VM under the hood to provide container support, Podman's architecture is more directly tied to Linux system calls and cannot be easily virtualized on macOS.

If you're using macOS, please use Docker mode instead.

## Prerequisites

- Podman installed on your system (Linux only)
- Proper system permissions configured

## Image Setup

Build the required Podman images using the provided Makefile targets:

```bash
# Build Alpine-based image
make build-podman-image

# Build Ubuntu-based image
make build-podman-image-ubuntu
```

## Configuration

Set the following environment variables:

```bash
ODIN_MODE=podman
ODIN_PODMAN_IMAGE=odin:alpine  # or odin:ubuntu
ODIN_PODMAN_NETWORK=<network_name>  # Optional
ODIN_WORKER_COUNT=<number_of_workers>
```

## Key Features

- Rootless container execution
- Daemonless architecture
- Compatible with Docker images
- Integrated with systemd

## Security Configuration

1. Configure SELinux labels (if applicable):
   ```bash
   ODIN_PODMAN_SELINUX_ENABLED=true
   ```

2. Set up rootless mode:
   ```bash
   # Configure subuid/subgid mappings
   sudo usermod --add-subuids 100000-165535 --add-subgids 100000-165535 $USER
   ```

## Resource Controls

```bash
ODIN_PODMAN_MEMORY_LIMIT=<memory_limit>
ODIN_PODMAN_CPU_LIMIT=<cpu_limit>
ODIN_PODMAN_PIDS_LIMIT=<process_limit>
```

## Troubleshooting

1. Check container status:
   ```bash
   podman ps -a --format "table {{.ID}} {{.Image}} {{.Status}}"
   ```

2. View container logs:
   ```bash
   podman logs <container_id>
   ```

3. Inspect container details:
   ```bash
   podman inspect <container_id>
   ```
