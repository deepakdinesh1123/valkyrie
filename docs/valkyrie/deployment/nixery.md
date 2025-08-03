# Self-Hosting Nixery

**Nixery** is an ad-hoc container image registry that dynamically provides container images with packages from the **Nix** package manager. You can pull container images directly from this registry at [nixery.dev](https://nixery.dev) by specifying the required packages in the URL.

You can also self-host Nixery as a Docker container or on a serverless compute platform like Google Cloud Run. **HTTPS** is mandatory for Nixery to function correctly.

## Self-Hosting with Docker

### Prerequisites

* A public domain name.
* A virtual machine with a public IP address.
* Caddy (for HTTPS reverse proxy).

### Installation Steps

1. **Clone the Nixery Repository:**

```bash
git clone https://github.com/tazjin/nixery.git
cd nixery
```

2. **Build Nixery Image:**

Run the following command to build and load the Nixery Docker image:

```bash
IMG=$(docker load -q -i "$(nix-build -A nixery-image)" | awk '{ print $3 }')
echo "Loaded Nixery image as ${IMG}"
```

The image name will be displayed once the process completes.

### Configuration

Run the Nixery container with environment variables to configure:

* `PORT`: Port where Nixery listens for HTTP requests.
* `NIXERY_CHANNEL`: Nix/NixOS channel for building images.
* `NIXERY_STORAGE_BACKEND`: Storage backend (`gcs` or `filesystem`).
* `GCS_BUCKET`: Name of your Google Cloud Storage bucket (required if using `gcs`).
* `GOOGLE_APPLICATION_CREDENTIALS`: Path to your GCP service account JSON key file (optional for `gcs`).

```bash
docker run --privileged --rm -p 127.0.0.1:8080:8080 --name nixery \
    -e PORT=8080 \
    -e NIXERY_CHANNEL=<nix-channel> \
    -e NIXERY_STORAGE_BACKEND=gcs \
    -e GCS_BUCKET=<gcs-bucket> \
    -e GOOGLE_APPLICATION_CREDENTIALS=/gcpkeys.json \
    -v /gcpkeys.json:/gcpkeys.json \
    ${IMG}
```

### Serving Nixery over HTTPS with Caddy

Configure Caddy to reverse-proxy your Nixery container securely:

```caddyfile
yourdomain.com {
    reverse_proxy 127.0.0.1:8080
}
```

### Usage

Once set up, use your private Nixery instance to pull container images:

```bash
docker run -it yourdomain.com/shell bash
```

## Self-Hosting on Google Cloud Platform (GCP)

You can also host your Nixery instance on **Google Cloud Run**, leveraging its built-in HTTPS and auto-scaling capabilities.

### Cloud Run Setup

1. **Build and Load Nixery Image:**

Follow the [Docker installation steps](#installation-steps) above to build the Nixery image.

2. **Deploy to Cloud Run:**

Deploy the built Nixery image with the following configurations:

* Set environment variables as detailed in the [Docker Configuration](#configuration) section.
* Attach your GCS bucket as a volume to the Cloud Run instance.
* Set up a Startup Probe with an initial delay of `120s` and a frequency of `30s`.
* Mount your service account key JSON as another volume, either from Secret Manager or directly from the GCS bucket.

Once deployed, your private Nixery instance will be available at:

`https://<cloudrun-service-name>.<region>.run.app`

### Usage

Pull container images directly from your Cloud Run Nixery instance:

```bash
docker run -it <cloudrun-service-name>.<region>.run.app/shell bash
```

---
