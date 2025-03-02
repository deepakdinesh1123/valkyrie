ARG UBUNTU_IMAGE=ubuntu:24.04
ARG NIXPKGS_REV=b27ba4eb322d9d2bf2dc9ada9fd59442f50c8d7c

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

COPY --chown=1024:1024 hack/sandbox/install_nix.bash /tmp/install_nix.bash

# Install Nix
RUN chmod +x /tmp/install_nix.bash && \
    bash /tmp/install_nix.bash && \
    # Cleanup the install script
    rm -f /tmp/install_nix.bash

# Ensure the Nix binaries are available in PATH
ENV PATH="/home/valnix/.nix-profile/bin:${PATH}"

RUN nix-env -iA nixpkgs.nix-serve-ng
RUN nix-store --generate-binary-cache-key odin-store cache-priv-key.pem cache-pub-key.pem && \
    chown valnix cache-priv-key.pem && \
    chmod 600 cache-priv-key.pem

ENTRYPOINT [ "/bin/sh", "-c", "NIX_SECRET_KEY_FILE=/home/valnix/cache-priv-key.pem nix-serve --listen 0.0.0.0:5000" ]
