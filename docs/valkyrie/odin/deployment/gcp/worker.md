- The Odin worker is deployed in the google cloud compute engine without a public IP.
- GCP instance groups with a minimum worker count of 1 is required.
- The worker setup: 
    - Instance groups are created with a preemptable VMs [(spot vm)](../../../../blog/posts/spot_instances.md)
    - [Shared-nix-store](./shared_nix_store.md)
    - Container Engine (Podman/Docker)
    - Golang
    - Environment variables
        * NIX_USER_ENVIRONMENT
        * NIX_CHANNELS_ENVIRONMENT

!!! info "Info"
    - The Odin server, worker and shared nix store should be 

    **Caddyfile**
    ```
    backend.evnix.cloud {  
        reverse_proxy localhost:8080  
    }
    ```