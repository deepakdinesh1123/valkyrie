FROM nixos/nix:latest

RUN echo 'experimental-features = nix-command flakes' >> /etc/nix/nix.conf
RUN nix-channel --add https://nixos.org/channels/nixpkgs-unstable unstable
RUN nix-channel --update
RUN nix-env -iA nixpkgs.go nixpkgs.unixtools.netstat nixpkgs.air nixpkgs.sqlc nixpkgs.go-migrate
WORKDIR /valkyrie/odin
COPY go.mod /valkyrie/odin
COPY go.sum /valkyrie/odin
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY cmd /valkyrie/odin/cmd
COPY internal /valkyrie/odin/internal
COPY database /valkyrie/odin/database
COPY configs/air/.air.odin.toml /valkyrie/odin/
COPY .env /valkyrie/odin/
CMD ["bin/bash", "-c", "sleep infinity"]