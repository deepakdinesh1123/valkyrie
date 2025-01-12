# Shared Nix Store
The Shared nix store is a persistent compute disk which contains all the packages in a nix channel.

## Prerequisites
- Disk attached as read-only to all the Odin workers.
- Disk should be in the same region as worker and must be comaptabile with gcp compute families (e-series, n-series ...etc)

## Disk population  
- Download the [.sqlite](https://valnix-stage-bucket.s3.amazonaws.com/rippkgs-24.11.sqlite) for nix 24.11 channel.
- Download the [odin binary](https://valnix-stage-bucket.s3.amazonaws.com/odinb).
- Run generate cmd to populate sqlite db
  ```
  ./odinb store generate
  ```
- Run realise cmd to populate nix packages
  ```
  ./odinb store realise
  ```

!!! info "**Disk Requirements**"
    Current Odin setup has a gcp persistent SSD disk, which allows read-only attachment to multiple compute instances.