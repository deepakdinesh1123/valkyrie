# Running Valkyrie in Local

This guide explains how to run the Valkyrie application stack using Docker Compose local with the `--profile dev`. This setup involves using [nixery.dev](https://Nixery.dev) to pull container images with packages. For detailed instructions on setting up private Nixery instance, refer to [Nixery Setup](../deployment/nixery.md).

---

## Prerequisites

Ensure your machine meets the following prerequisites:

* **Docker and Docker Compose** installed:

```bash
sudo apt-get update
sudo apt-get install docker.io docker-compose -y
```

* **Git** installed to clone repositories:

```bash
sudo apt-get install git -y
```

---

## Clone the Valkyrie Repository

Clone the Valkyrie repository from GitHub:

```bash
git clone -b feat/nixery https://github.com/deepakdinesh1123/valkyrie.git
cd valkyrie
```

---

# Valkyrie Environment Configuration

This document explains the configuration options (environment variables) required to configure your self-hosted Valkyrie environment.

## Environment Variables

Below is the list of environment variables along with explanations of their purposes:

| Variable Name       | Default Value                             | Description                                                              |
| ------------------- | ----------------------------------------- | ------------------------------------------------------------------------ |
| `POSTGRES_HOST`     | `localhost`                               | Hostname or IP address of your PostgreSQL database.                      |
| `POSTGRES_DB`       | `valkyrie`                                | Name of the PostgreSQL database to be used by Valkyrie.                  |
| `POSTGRES_USER`     | `thors`                                   | Username for PostgreSQL database authentication.                         |
| `POSTGRES_PASSWORD` | `thorkell`                                | Password for PostgreSQL database authentication.                         |
| `POSTGRES_PORT`     | `5432`                                    | Port number on which your PostgreSQL database is listening.              |
| `POSTGRES_SSL_MODE` | `disable`                                 | SSL mode for PostgreSQL connection (`disable`, `require`, etc.).         |
| `CONTAINER_RUNTIME` | `runsc`                                   | Container runtime to use (`runsc` for gVisor or `runc` default).         |
| `DB_MIGRATE`        | `true`                                    | Enables or disables automatic database migration on application startup. |
| `ENABLE_EXECUTION`  | `true`                                    | Enables or disables the execution functionality within Valkyrie.         |
| `NIXERY_URL`        | `nixery.dev` | URL of your Nixery instance for fetching container images dynamically. Use custom URL if using private nixery instance    |

## Example `.env` File

Here’s an example of a `.env` file for Valkyrie:

```env
POSTGRES_HOST=localhost
POSTGRES_DB=valkyrie
POSTGRES_USER=thors
POSTGRES_PASSWORD=thorkell
POSTGRES_PORT=5432
POSTGRES_SSL_MODE=disable
CONTAINER_RUNTIME=runsc
DB_MIGRATE=true
ENABLE_EXECUTION=true
NIXERY_URL=nixery.dev
```

---

## Configure Docker Runtime to use gVisor (Optional but Recommended for Security)

To enhance container isolation, you can change the container runtime from `runc` to `runsc` using gVisor. Follow these steps:

### Install gVisor

```bash
curl -fsSL https://gvisor.dev/archive.key | sudo gpg --dearmor -o /usr/share/keyrings/gvisor-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/gvisor-archive-keyring.gpg] https://storage.googleapis.com/gvisor/releases release main" | sudo tee /etc/apt/sources.list.d/gvisor.list
sudo apt-get update
sudo apt-get install -y runsc
```

### Configure Docker to use gVisor

Edit Docker's daemon configuration file:

```bash
sudo nano /etc/docker/daemon.json
```

Add the following content:

```json
{
  "runtimes": {
    "runsc": {
      "runtimeType": "io.containerd.runsc.v1",
      "options": {
        "TypeUrl": "io.containerd.runsc.v1.options",
        "ConfigPath": "/etc/containerd/runsc.toml"
      }
    }
  }
}
```

Restart Docker to apply the changes:

```bash
sudo systemctl restart docker
```

In your `.env`, specify the runtime for the services you wish to run with gVisor:

```
CONTAINER_RUNTIME=runsc
```

---

## Build and Run Valkyrie with Docker Compose

### Valkyrie Docker Networks

This document explains how to create Docker networks required for your Valkyrie setup.

#### Docker Network Setup

Docker networks help in isolating containers and managing their internal communication effectively. Valkyrie requires the following Docker networks:

##### Networks

* `valkyrie-network`: Main network for Valkyrie components to communicate internally.
* `devpi-network`: Network specifically used for Devpi components.

##### Creating Docker Networks

Execute the following commands to create these networks:

```bash
docker network create valkyrie-network
docker network create devpi-network
```

##### Verify Network Creation

You can verify that the networks were successfully created by running:

```bash
docker network ls
```

You should see output similar to:

```
NETWORK ID     NAME               DRIVER    SCOPE
abc123def456   valkyrie-network   bridge    local
def789abc123   devpi-network      bridge    local
```

### Build Docker Images

Build all the necessary Docker images with the `--profile dev`:

```bash
docker compose --profile dev build
```

### Run Containers

Place the `.env` file in the root directory of your Valkyrie Docker Compose setup. Docker Compose will automatically read and apply the environment variables when starting the services.
Start the containers with the dev profile:

```bash
docker compose --profile dev up -d
```

Check the status of running containers:

```bash
docker compose ps
```

---

## Accessing Valkyrie

Your application should now be accessible via HTTPS at:

```bash
https://localhost:8080
```

---

## Monitoring and Tracing (Optional)

If you want to enable detailed tracing with Jaeger, run the following commands:

### Build and Run with Instrumentation

```bash
docker compose --profile instrument build
docker compose --profile instrument up -d
```

Jaeger UI will be available at:

```bash
http://localhost:16686
```

Use Jaeger to view detailed application traces for debugging and monitoring.

---

## Troubleshooting

### Check logs for any container issues:

```bash
docker compose logs -f <container-name>
```

Replace `<container-name>` with the respective container.

---

## Conclusion

You now have a fully operational Valkyrie environment deployed using Docker Compose in your machine and detailed tracing capabilities provided by Jaeger.
