ARG NIX_CHANNEL=24.11

FROM nixos/nix:2.24.0 AS builder

WORKDIR /valkyrie
COPY flake.nix /valkyrie
COPY flake.lock /valkyrie

COPY cmd /valkyrie/cmd
COPY builds/nix/valkyrie.nix /valkyrie/builds/nix/valkyrie.nix
COPY internal /valkyrie/internal
COPY pkg /valkyrie/pkg
COPY go.mod /valkyrie
COPY go.sum /valkyrie

RUN \
    --mount=type=cache,target=/nix,from=nixos/nix:2.24.0,source=/nix \
    mkdir -p /tmp/nix-store-closure  && \
    nix \
    --extra-experimental-features "nix-command flakes" \
    --option filter-syscalls false \
    --show-trace \
    --log-format raw \
    build --out-link /tmp/valkyrie/result && \
    cp -R $(nix-store -qR /tmp/valkyrie/result) /tmp/nix-store-closure

FROM ubuntu:24.04

RUN apt update && \
    apt install -y adduser curl xz-utils && \
    groupadd -o -g 1024 -r valnix && \
    adduser --uid 1024 --gid 1024 --disabled-password --gecos "" valnix

COPY --from=builder /tmp/nix-store-closure /nix/store
COPY --from=builder /tmp/valkyrie/ /home/valnix

USER valnix
WORKDIR /home/valnix
ENTRYPOINT ["/home/valnix/result/bin/valkyrie"]
