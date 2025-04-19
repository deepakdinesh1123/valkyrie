ARG UBUNTU_IMAGE=ubuntu:24.04
ARG NIXPKGS_REV=b27ba4eb322d9d2bf2dc9ada9fd59442f50c8d7c
ARG NIX_CACHE_PUBLIC_KEY

FROM ${UBUNTU_IMAGE}

# Install required packages and configure system
RUN apt update && \
    apt install -y adduser curl xz-utils unzip && \
    groupadd -o -g 1024 -r valnix && \
    adduser --uid 1024 --gid 1024 --disabled-password --gecos "" valnix && \
    mkdir /nix && \
    mkdir -p /home/valnix/.config/process-compose && \
    chown -R valnix:valnix /nix /home/valnix && \
    # Configure Nix
    mkdir /etc/nix

# Switch to the valnix user
USER valnix
WORKDIR /home/valnix/

COPY --chown=valnix:valnix hack/install_nix.bash /tmp/install_nix.bash

# Install Nix
RUN chmod +x /tmp/install_nix.bash && \
    bash /tmp/install_nix.bash && \
    # Cleanup the install script
    rm -f /tmp/install_nix.bash

# Ensure the Nix binaries are available in PATH
ENV PATH="/home/valnix/.nix-profile/bin:${PATH}"

USER root
COPY configs/nix/nix.conf /etc/nix/nix.conf
ARG NIX_CACHE_PUBLIC_KEY
RUN echo "trusted-public-keys = ${NIX_CACHE_PUBLIC_KEY} cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY=" >> /etc/nix/nix.conf

USER valnix
RUN mkdir -p /home/valnix/valkyrie
COPY --chown=valnix:valnix hack/execution/* /home/valnix/

VOLUME [ "/home/valnix" ]
WORKDIR /home/valnix/

USER root
RUN apt update && apt install -y ca-certificates
USER valnix

ENTRYPOINT [ "/bin/sh", "-c", "sleep infinity"]
