ARG NIX_CHANNEL=24.11
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

# Install Nixsource
RUN curl -L https://nixos.org/nix/install -o /home/valnix/install_nix.sh && \
    chmod +x install_nix.sh && \
    sh install_nix.sh --no-daemon && \
    # Cleanup the install script
    rm -f install_nix.sh

# Ensure the Nix binaries are available in PATH
ENV PATH="/home/valnix/.nix-profile/bin:${PATH}"

ARG NIX_CHANNEL
RUN nix-channel --remove nixpkgs
RUN nix-channel --add https://nixos.org/channels/nixos-${NIX_CHANNEL} nixpkgs && nix-channel --update

RUN nix-env -iA nixpkgs.nix-serve-ng
RUN nix-store --generate-binary-cache-key odin-store cache-priv-key.pem cache-pub-key.pem && \
    chown valnix cache-priv-key.pem && \
    chmod 600 cache-priv-key.pem

ENTRYPOINT [ "/bin/sh", "-c", "NIX_SECRET_KEY_FILE=/home/valnix/cache-priv-key.pem nix-serve --listen 0.0.0.0:5000" ]