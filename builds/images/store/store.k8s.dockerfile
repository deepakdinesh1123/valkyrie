FROM ubuntu:24.04

ARG DEBIAN_FRONTEND=noninteractive

RUN apt update && \
    apt install -y adduser curl xz-utils sudo && \
    groupadd -o -g 1024 -r valnix && \
    adduser --uid 1024 --gid 1024 --disabled-password --gecos "" valnix && \
    usermod -aG sudo valnix && \
    echo "valnix ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/valnix && \
    # Configure Nix
    mkdir /etc/nix && \
    echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf

# Switch to the valnix user and set the working directory
USER valnix
WORKDIR /home/valnix

# Copy scripts and secrets
COPY --chown=valnix:valnix hack/install_nix.bash /home/valnix/install_nix.bash
COPY --chown=valnix:valnix hack/store/install_store.bash /home/valnix/install_store.bash
COPY --chown=valnix:valnix hack/store/start_store.bash /home/valnix/start_store.bash
COPY --chown=valnix:valnix configs/secrets/cache-priv-key.pem /home/valnix/.config/cache-priv-key.pem
COPY --chown=valnix:valnix configs/secrets/cache-pub-key.pem /home/valnix/.config/cache-pub-key.pem

# Make scripts executable
RUN chmod +x /home/valnix/install_nix.bash /home/valnix/install_store.bash /home/valnix/start_store.bash

# Install Nix and then start the store
CMD ["/bin/bash",  "-c" ,"./install_store.bash k8s && ./start_store.bash"]
