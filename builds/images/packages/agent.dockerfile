ARG NIX_CHANNEL=24.11

FROM nixos/nix:2.24.0

RUN mkdir -p /etc/nix && echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf

ARG NIX_CHANNEL

RUN nix-channel --remove nixpkgs
RUN nix-channel --add https://nixos.org/channels/nixos-${NIX_CHANNEL} nixpkgs && nix-channel --update

COPY builds/base.go.nix /tmp/go/shell.nix
WORKDIR /tmp/go
RUN nix-shell --run 'exit'

WORKDIR /valkyrie
COPY flake.nix /valkyrie
COPY flake.lock /valkyrie
RUN nix develop .#agent --command 'ls'

COPY agent /valkyrie/agent
COPY builds/nix/agent.nix /valkyrie/builds/nix/agent.nix

RUN nix build .#agent

RUN mkdir /tmp/nix-store-closure
RUN cp -R $(nix-store -qR result/) /tmp/nix-store-closure
