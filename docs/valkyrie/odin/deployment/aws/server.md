The Odin server is deployed in the AWS ec2 with a public endpoint. 

## Prerequisites
- Caddy installed and running on your system
- Golang installed

## Server Setup 
- Caddy will be running as a reverse proxy, routing traffic from client requests to the Odin server running on a compute engine.
- The environment variables will include 
    * DB Secrets 
    * Odin worker execution mode
    * Odin worker container Engine (Podman, Docker)
    * Odin worker base image
    * Odin worker hot containers
    * Log level

!!! info "**Caddyfile**"
    ```
    backend.evnix.cloud {  
        reverse_proxy localhost:8080  
    }
    ```