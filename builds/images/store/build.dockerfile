FROM ubuntu:24.04

# Install required packages and configure system
RUN apt update && \
    apt install -y adduser curl xz-utils sudo && \
    groupadd -o -g 1024 -r valnix && \
    adduser --uid 1024 --gid 1024 --disabled-password --gecos "" valnix && \
    usermod -aG sudo valnix && \
    echo "valnix ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/valnix && \
    # Configure Nix
    mkdir /etc/nix && \
    echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf

# Switch to the valnix user
USER valnix
WORKDIR /home/valnix/

COPY --chown=1024:1024 hack/store/install_nix.bash /tmp/install_nix.bash

# Install Nix
RUN chmod +x /tmp/install_nix.bash
CMD [ "/bin/bash", "/tmp/install_nix.bash" ]
