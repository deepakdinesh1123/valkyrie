ARG NIX_CHANNEL=24.11

FROM nixos/nix:2.24.0 AS builder

RUN mkdir -p /etc/nix && echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf

ARG NIX_CHANNEL

RUN nix-channel --remove nixpkgs
RUN nix-channel --add https://nixos.org/channels/nixos-${NIX_CHANNEL} nixpkgs && nix-channel --update

WORKDIR /valkyrie
COPY flake.nix /valkyrie
COPY flake.lock /valkyrie
RUN nix develop .#agent --command 'ls'

COPY cmd /valkyrie/cmd
COPY builds/nix/valkyrie.nix /valkyrie/builds/nix/valkyrie.nix
COPY internal /valkyrie/internal
COPY pkg /valkyrie/pkg
COPY go.mod /valkyrie
COPY go.sum /valkyrie

RUN nix build

RUN mkdir /tmp/nix-store-closure
RUN cp -R $(nix-store -qR result/) /tmp/nix-store-closure

FROM ubuntu:24.04

RUN apt update && \
    apt install -y adduser curl xz-utils && \
    groupadd -o -g 1024 -r valnix && \
    adduser --uid 1024 --gid 1024 --disabled-password --gecos "" valnix

COPY --from=builder /tmp/nix-store-closure /nix/store
COPY --from=builder /valkyrie/result /home/valnix

USER valnix
WORKDIR /home/valnix
ENTRYPOINT ["/home/valnix/bin/valkyrie"]
