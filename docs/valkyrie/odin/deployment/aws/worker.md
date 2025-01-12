# Odin Worker
The Odin worker is deployed in the AWS ec2 without a public endpoint. Workers are run in aws spot fleet, this fleet help us maintain minimum number of workers and are scalable as utilization.
## Prerequisites
- Container Engine (Podman/Docker) installed and running on your system
- Golang installed

## Worker setup
- AWS spot fleet with a minimum worker count of 1 is required.
  Fleet instances can be of any supported ec2 familiy

- Instance groups are created with a preemptable VMs (spot vm)
- [Shared-nix-store](./shared_nix_store.md)
- Environment variables
    * NIX_USER_ENVIRONMENT
    * NIX_CHANNELS_ENVIRONMENT
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
- For convenience, [worker setup script](https://valnix-stage-bucket.s3.amazonaws.com/deploy.sh) has been added.  
  Although the  script will be executed in worker while setting up using opentofu it will serve as a reference.  
  The script has aws secrets in it.
- Run the [odin binary](https://valnix-stage-bucket.s3.amazonaws.com/odinb)
  ```
  ./odinb worker start
  ```


!!! note "Note"
    - The shared nix store should be in same aws region as Odin server and worker for reduced Latency, network speed and reduced
      the data transfer costs.
