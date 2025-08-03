# Running Valkyrie Docker Compose Setup on a VM

This guide explains how to run the Valkyrie application stack using Docker Compose on a virtual machine (VM) with the `--profile dev`. This setup includes a reverse proxy with Caddy and optional instrumentation to view application traces using Jaeger.

---

## Prerequisites

Ensure your VM meets the following prerequisites:

* **Docker and Docker Compose** installed:

```bash
sudo apt-get update
sudo apt-get install docker.io docker-compose -y
```

* **Git** installed to clone repositories:

```bash
sudo apt-get install git -y
```

* **Caddy** installed for HTTPS reverse proxy:

```bash
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/caddy-stable-archive-keyring.gpg] https://dl.cloudsmith.io/public/caddy/stable/deb/debian any-version main" | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install caddy -y
```

---

## Clone the Valkyrie Repository

Clone the Valkyrie repository from GitHub:

```bash
git clone -b feat/nixery https://github.com/deepakdinesh1123/valkyrie.git
cd valkyrie
```

---

## Configure Caddy

Use the following configuration in your Caddyfile. Replace `youdomain.io` with your domain:

```caddyfile
youdomain.io {
    reverse_proxy valkyrie:8080 {
        health_uri /api/health
        health_interval 30s
        health_timeout 5s

        # Header forwarding for proper client IP detection
        header_up Host {upstream_hostport}
        header_up X-Real-IP {remote_host}
        header_up X-Forwarded-For {remote_host}
        header_up X-Forwarded-Proto {scheme}
    }

    log {
        output file /var/log/caddy/access.log
        format json
    }

    encode gzip zstd
}
```

Ensure your DNS settings point your domain to your VM's public IP.

---

## Build and Run Valkyrie with Docker Compose

### Build Docker Images

Build all the necessary Docker images with the `--profile dev`:

```bash
docker compose --profile dev build
```

### Run Containers

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
https://yourdomain.io
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
http://<your-vm-ip-address>:16686
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

You now have a fully operational Valkyrie environment deployed using Docker Compose with a reverse proxy via Caddy and detailed tracing capabilities provided by Jaeger.
