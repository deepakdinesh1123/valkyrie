ARG UBUNTU_IMAGE=ubuntu:24.04

FROM ${UBUNTU_IMAGE}

# Install required packages and configure system
RUN apt update && \
    apt install -y adduser curl xz-utils && \
    groupadd -o -g 1024 -r valnix && \
    adduser --uid 1024 --gid 1024 --disabled-password --gecos "" valnix && \
    mkdir /nix && \
    mkdir -p /home/valnix/.config/process-compose && \
    chown -R valnix:valnix /nix /home/valnix && \
    # Configure Nix
    mkdir /etc/nix && \
    echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf

# Switch to the valnix user
USER valnix
WORKDIR /home/valnix/

ARG NIX_CHANNELS_ENVIRONMENT
ARG NIX_USER_ENVIRONMENT

ENV NIX_CHANNELS_ENVIRONMENT=${NIX_CHANNELS_ENVIRONMENT}
ENV NIX_USER_ENVIRONMENT=${NIX_USER_ENVIRONMENT}

COPY --chown=1024:1024 hack/store/start_store.bash /tmp/start_store.bash

ENTRYPOINT [ "/bin/bash", "/tmp/start_store.bash" ]
