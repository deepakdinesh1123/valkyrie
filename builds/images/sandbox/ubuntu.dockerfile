ARG AGENT_BUILDER=valkyrie_agent:0.0.1
ARG UBUNTU_IMAGE=ubuntu:24.04
ARG NIXPKGS_REV=b27ba4eb322d9d2bf2dc9ada9fd59442f50c8d7c
ARG NIX_CACHE_PUBLIC_KEY

FROM ${AGENT_BUILDER} AS agent

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

COPY --chown=valnix:valnix hack/sandbox/install_nix.bash /tmp/install_nix.bash

# Install Nix
RUN chmod +x /tmp/install_nix.bash && \
    bash /tmp/install_nix.bash && \
    # Cleanup the install script
    rm -f /tmp/install_nix.bash

USER root
ARG NIXPKGS_REV
RUN curl -L "https://api.github.com/repos/NixOS/nixpkgs/zipball/${NIXPKGS_REV}" -o commit.zip && unzip commit.zip -d /var/cache/nixpkgs && rm commit.zip
RUN curl -L "https://github.com/hercules-ci/flake-parts/archive/refs/heads/main.zip" -o parts.zip && unzip parts.zip -d /var/cache && rm parts.zip
RUN curl -L "https://github.com/Platonic-Systems/process-compose-flake/archive/refs/heads/main.zip" -o pcompose.zip && unzip pcompose.zip -d /var/cache && rm pcompose.zip
RUN curl -L "https://github.com/juspay/services-flake/archive/refs/heads/main.zip" -o services.zip && unzip services.zip -d /var/cache && rm services.zip

USER valnix

# Ensure the Nix binaries are available in PATH
ENV PATH="/home/valnix/.nix-profile/bin:${PATH}"

RUN mkdir -p /home/valnix/work

COPY --chown=valnix:valnix configs/nix/flake.nix /home/valnix/work/flake.nix
RUN cd /home/valnix/work && nix profile install . --extra-experimental-features 'nix-command flakes'

RUN nix profile install nixpkgs#vim --extra-experimental-features 'nix-command flakes'
RUN nix profile install nixpkgs#gnupatch --extra-experimental-features 'nix-command flakes'

COPY --from=agent /tmp/nix-store-closure /tmp/agent/closure
COPY --from=agent /valkyrie/result /home/valnix

USER root
COPY configs/nix/nix.conf /etc/nix/nix.conf
ARG NIX_CACHE_PUBLIC_KEY
RUN echo "trusted-public-keys = ${NIX_CACHE_PUBLIC_KEY} cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY=" >> /etc/nix/nix.conf
RUN chown -R valnix:valnix /tmp/agent/closure
RUN cp -a /tmp/agent/closure/* /nix/store

USER valnix
WORKDIR /home/valnix/work/

ENTRYPOINT [ "/home/valnix/bin/agent"]
