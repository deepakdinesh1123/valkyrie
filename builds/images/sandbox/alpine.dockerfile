ARG NIX_CHANNEL=24.11
ARG ALPINE_IMAGE=alpine:3.20

FROM ${ALPINE_IMAGE}

RUN apk update && \
    apk add --no-cache curl xz sudo && \
    addgroup -g 1024 -S valnix && \
    adduser -D -u 1024 -G valnix -h /home/valnix valnix && \
    mkdir /nix && \
    mkdir -p /home/valnix/.config/process-compose && \
    chown -R valnix:valnix /nix /home/valnix && \
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

ENTRYPOINT [ "/bin/sh", "-c", "sleep infinity"]