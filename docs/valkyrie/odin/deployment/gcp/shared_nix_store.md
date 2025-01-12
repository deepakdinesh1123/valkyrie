The Shared nix store is a persistent compute disk which contains all the packages in a nix channel.

## Prerequisites
- Disk attached as read-only to all the Odin workers.
- Disk should be in the same region as worker and must be comaptabile with gcp compute families (e-series, n-series ...etc)

## Disk population  
- 
- 

!!! info "**Disk Requirements**"
    Current Odin setup has a gcp persistent SSD disk, which allows only read-only attachment to multiple compute instances.