FROM nixos/nix:latest

RUN echo 'experimental-features = nix-command flakes' >> /etc/nix/nix.conf
RUN nix-channel --add https://nixos.org/channels/nixpkgs-unstable unstable
RUN nix-channel --update
RUN nix-env -iA nixpkgs.go nixpkgs.unixtools.netstat nixpkgs.air nixpkgs.crun nixpkgs.openssh
WORKDIR /valkyrie/skald
COPY go.mod /valkyrie/skald/
COPY go.sum /valkyrie/skald/
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY cmd /valkyrie/skald/cmd
COPY internal /valkyrie/skald/internal
COPY configs/air/.air.skald.toml /valkyrie/skald/
COPY .env /valkyrie/skald/
RUN ssh-keygen -t rsa -f /root/.ssh/id_rsa -q -N ''
RUN mkdir -p /valkyrie/config/keys
RUN cp /root/.ssh/id_rsa.pub /valkyrie/config/keys/
CMD ["bin/bash", "-c", "sleep infinity"]