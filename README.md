# Valkyrie

Valkyrie is a powerful execution engine and sandbox provider that allows users to execute code or create sandboxes. Utilizing Nix to manage dependencies and containers for isolation, Valkyrie offers a secure, consistent, and isolated environment for running and testing code.

## Features

- **Code Execution**: Run your code in a secure and isolated environment powered by containers.
- **Sandbox Creation**: Create sandboxes for development, testing, or any other purpose with ease.
- **Powered by Nix**: Leverage the power of Nix for highly reproducible environments.
- **Containerization**: Utilize container technology for isolation and security.

## Getting Started

Follow these steps to get started with Valkyrie:

### Prerequisites

The only dependency required to run valkyrie is docker (podman will also be supported in the future) you can find docker installation instructions [here](https://www.docker.com/get-started)

### Setup

- Clone the repository and navigate to the project directory:

```sh
git clone https://github.com/yourusername/valkyrie.git
cd valkyrie
```

- Setup Valkyrie store

```sh
docker compose --profile setup build
docker compose --profile setup up
```

- Run Valkyrie

```sh
docker compose --profile dev build
docker compose --profile dev up
```

**Note:** Valkyrie can also be run with Tetragon for enhanced runtime security, and with Prometheus and Jaeger for comprehensive metrics collection. Use the following command:

```sh
docker compose --profile staging build
docker compose --profile staging up
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Nix](https://nixos.org/)
- [Docker](https://www.docker.com/)

## Contact

For any inquiries, please reach out to us at [d.deepakdinesh13@gmail.com](mailto:d.deepakdinesh13@gmail.com).
