FROM nixos/nix:2.23.1 AS builder

RUN mkdir -p /etc/nix && echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf

WORKDIR /valkyrie
COPY flake.nix /valkyrie
COPY flake.lock /valkyrie
RUN nix develop --command 'ls'

COPY cmd /valkyrie/cmd
COPY builds/nix/odin.nix /valkyrie/builds/nix/odin.nix
COPY internal /valkyrie/internal
COPY pkg /valkyrie/pkg
COPY go.mod /valkyrie
COPY go.sum /valkyrie

RUN nix build

RUN mkdir /tmp/nix-store-closure
RUN cp -R $(nix-store -qR result/) /tmp/nix-store-closure

FROM alpine:3.20

RUN apk update && \
    apk add --no-cache shadow && \
    addgroup -g 1024 -S valnix && \
    adduser -u 1024 -G valnix -S -D valnix

COPY --from=builder /tmp/nix-store-closure /nix/store
COPY --from=builder /valkyrie/result /home/valnix

USER valnix
WORKDIR /home/valnix
ENTRYPOINT ["/home/valnix/bin/odin"]
