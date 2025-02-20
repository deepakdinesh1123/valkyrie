FROM nixos/nix:2.23.1

RUN mkdir -p /etc/nix && echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf

WORKDIR /valkyrie
COPY flake.nix /valkyrie
COPY flake.lock /valkyrie
RUN nix develop .#agent --command 'ls'

COPY agent /valkyrie/agent
COPY builds/nix/agent.nix /valkyrie/builds/nix/agent.nix

RUN nix build .#agent

RUN mkdir /tmp/nix-store-closure
RUN cp -R $(nix-store -qR result/) /tmp/nix-store-closure
