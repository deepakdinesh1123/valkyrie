ARG AGENT_BUILDER=odin_agent:0.0.1
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

COPY --chown=1024:1024 hack/sandbox/install_nix.bash /tmp/install_nix.bash

# Install Nix
RUN chmod +x /tmp/install_nix.bash && \
    bash /tmp/install_nix.bash && \
    # Cleanup the install script
    rm -f /tmp/install_nix.bash

USER root
ARG NIXPKGS_REV
RUN curl -L "https://api.github.com/repos/NixOS/nixpkgs/zipball/${NIXPKGS_REV}" -o commit.zip && unzip commit.zip -d /var/cache/nixpkgs && rm commit.zip
RUN curl -L "https://api.github.com/repos/numtide/flake-utils/zipball/11707dc2f618dd54ca8739b309ec4fc024de578b" -o utils.zip && unzip utils.zip -d /var/cache/flake-utils && rm utils.zip
USER valnix

# Ensure the Nix binaries are available in PATH
ENV PATH="/home/valnix/.nix-profile/bin:${PATH}"

# RUN nix profile install nixpkgs#nix-direnv --extra-experimental-features 'nix-command flakes'
# RUN nix profile install nixpkgs#direnv --extra-experimental-features 'nix-command flakes'
# RUN echo 'source $HOME/.nix-profile/share/nix-direnv/direnvrc' >> /home/valnix/.bashrc && \
#     echo 'eval "$(direnv hook bash)"' >> /home/valnix/.bashrc

COPY --from=agent /tmp/nix-store-closure /tmp/agent/closure
COPY --from=agent /valkyrie/result /home/valnix

COPY --chown=1024:1024 configs/nix/flake.nix /home/valnix/flake.nix
RUN nix profile install . --extra-experimental-features 'nix-command flakes'

USER root
COPY configs/nix/nix.conf /etc/nix/nix.conf
ARG NIX_CACHE_PUBLIC_KEY
RUN echo "trusted-public-keys = ${NIX_CACHE_PUBLIC_KEY} cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY=" >> /etc/nix/nix.conf
RUN chown -R valnix:valnix /tmp/agent/closure
RUN cp -a /tmp/agent/closure/* /nix/store

USER valnix
WORKDIR /home/valnix/

ENTRYPOINT [ "/home/valnix/bin/agent"]
