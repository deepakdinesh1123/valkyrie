# Odin Server
The Odin server is deployed in the google cloud compute engine with a public endpoint. 

## Prerequisites
- Golang installed
- Caddyfile
    ```
    example.com {  
        reverse_proxy localhost:8080  
    }
    ```

## Server Setup 
- Caddy will be running as a reverse proxy, routing traffic from client requests to the Odin server running on a compute engine.  
  Run caddy as sudo 
  ```
  sudo $(which caddy) run --config ./Caddyfile
  ```
- The environment variables will include 
    * POSTGRES_HOST=host
    * POSTGRES_DB=dbname
    * POSTGRES_USER=user
    * POSTGRES_PASSWORD=password
    * POSTGRES_PORT=port
    * POSTGRES_SSL_MODE=mode
    * DB_MIGRATE=bool 
    * ODIN_LOG_LEVEL=debug
    * ODIN_WORKER_EXECUTOR=container
    * ODIN_CONTAINER_ENGINE=podman
    * ODIN_WORKER_SYSTEM_EXECUTOR=native
    * ODIN_NIX_STORE=/nix
    * ODIN_WORKER_PODMAN_IMAGE=odin:0.0.3
    * ODIN_HOT_CONTAINER=1
- Run the [odin binary](https://valnix-stage-bucket.s3.amazonaws.com/odinb) 
  ```
  ./odinb server start
  ```