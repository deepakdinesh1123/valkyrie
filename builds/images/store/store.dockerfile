ARG BASE_IMAGE=odin_sandbox:ubuntu

FROM ${BASE_IMAGE}

RUN nix-env -iA nixpkgs.nix-serve-ng
ENTRYPOINT [ "/bin/sh", "-c", "nix-serve --listen 0.0.0.0:5000" ]