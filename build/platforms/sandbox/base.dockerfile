FROM ubuntu:24.04

# Copy the required files
COPY hack/sandbox/start.sh /home/valnix/start.sh
# COPY configs/sandbox/code-server.yaml /home/valnix/.config/code-server/config.yaml

# Install required packages and configure system
RUN apt update && \
    apt install -y adduser curl xz-utils && \
    groupadd -o -g 1024 -r valnix && \
    adduser --uid 1024 --gid 1024 --disabled-password --gecos "" valnix && \
    echo 'valnix ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers && \
    mkdir /nix && \
    mkdir -p /home/valnix/.config/process-compose && \
    chown -R valnix:valnix /nix /home/valnix && \
    chmod +x /home/valnix/start.sh && \
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

RUN nix-env -iA nixpkgs.code-server
RUN mkdir -p /home/valnix/.config/code-server

# COPY configs/sandbox/code-server.yaml /home/valnix/.config/code-server/config.yaml

# Entry point to start the sandbox environment
# ENTRYPOINT ["/home/valnix/start.sh"]
ENTRYPOINT [ "/bin/sh", "-c", "sleep infinity"]
