- The Odin server is deployed in the google cloud compute engine with a public IP. 
- Caddy will be running as a reverse proxy, routing traffic from client requests to the Odin server running on a compute engine.
- The environment variables will include 
    * DB Secrets 
    * Odin worker execution mode
    * Odin worker container Engine (Podman, Docker)
    * Odin worker base image
    * Odin worker hot containers
    * Log level

!!! note "Prerequisites"
    - Caddy

    **Caddyfile**
    ```
    backend.evnix.cloud {  
        reverse_proxy localhost:8080  
    }
    ```