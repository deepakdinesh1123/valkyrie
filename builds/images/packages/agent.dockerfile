ARG NIX_CHANNEL=24.11

FROM nixos/nix:2.24.0

RUN mkdir -p /etc/nix && echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf

WORKDIR /valkyrie
COPY flake.nix /valkyrie
COPY flake.lock /valkyrie

COPY agent /valkyrie/agent
COPY builds/nix/agent.nix /valkyrie/builds/nix/agent.nix

RUN \
    --mount=type=cache,target=/nix,from=nixos/nix:2.24.0,source=/nix \
    mkdir -p /tmp/nix-store-closure  && \
    nix \
    --extra-experimental-features "nix-command flakes" \
    --option filter-syscalls false \
    --show-trace \
    --log-format raw \
    build .#agent --out-link /tmp/agent/result && \
    cp -R $(nix-store -qR /tmp/agent/result) /tmp/nix-store-closure

CMD ["/bin/sh", "-c", "sleep infinity"]
