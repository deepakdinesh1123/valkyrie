# Shared Nix Store
The Shared nix store is AWS ebs which contains all the packages in a nix channel.

## Prerequisites
- EBS is created with multi-attach enabled

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

