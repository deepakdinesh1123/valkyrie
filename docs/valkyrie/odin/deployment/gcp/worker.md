The Odin worker is deployed in the google cloud compute engine without a public endpoint. Workers are run in gcp instance groups, these groups help us maintain minimum number of workers and are scalable as utilization.
## Prerequisites
- Container Engine (Podman/Docker) installed and running on your system
- Golang installed

## Worker setup: 
- GCP instance groups with a minimum worker count of 1 is required.

- Instance groups are created with a preemptable VMs [(spot vm)](../../../../blog/posts/spot_instances.md)
- [Shared-nix-store](./shared_nix_store.md)
- Environment variables
    * NIX_USER_ENVIRONMENT
    * NIX_CHANNELS_ENVIRONMENT
    * DB Secrets 
    * Odin worker execution mode
    * Odin worker container Engine (Podman, Docker)
    * Odin worker base image
    * Odin worker hot containers
    * Log level

!!! success "Tip"
    - The Odin server, worker and shared nix store should be in same gcp region for compatibility.
