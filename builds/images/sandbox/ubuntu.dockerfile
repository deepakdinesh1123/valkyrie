ARG NIX_CHANNEL=24.11
ARG AGENT_BUILDER=odin_agent:0.0.1
ARG UBUNTU_IMAGE=ubuntu:24.04
ARG ODIN_STORE_BUILDER=odin_store:0.0.1

FROM ${AGENT_BUILDER} AS agent

FROM ${ODIN_STORE_BUILDER} AS odin_store

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
    mkdir /etc/nix

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

COPY --from=agent /tmp/nix-store-closure /tmp/agent/closure
COPY --from=agent /valkyrie/result /home/valnix

USER root
COPY configs/nix/nix.conf /etc/nix/nix.conf
COPY --from=odin_store /home/valnix/cache-pub-key.pem /tmp/cache-pub-key.pem
RUN echo "trusted-public-keys = $(cat /tmp/cache-pub-key.pem) cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY=" >> /etc/nix/nix.conf
RUN chown -R valnix:valnix /tmp/agent/closure
RUN cp -a /tmp/agent/closure/* /nix/store

USER valnix
WORKDIR /home/valnix

ENTRYPOINT [ "/home/valnix/bin/agent"]
# CMD [ "/bin/sh", "-c", "sleep infinity" ]
