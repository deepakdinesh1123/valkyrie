FROM nixos/nix:2.23.1 as BUILDER

COPY cmd /valkyrie/cmd
COPY build/package/nix/odin.nix /valkyrie/build/package/nix/odin.nix
COPY internal /valkyrie/internal
COPY pkg /valkyrie/pkg
COPY go.mod /valkyrie
COPY go.sum /valkyrie
COPY flake.nix /valkyrie
COPY flake.lock /valkyrie

WORKDIR /valkyrie
RUN nix \
    --extra-experimental-features "nix-command flakes" \
    --option filter-syscalls false \
    build .#odin

RUN mkdir /tmp/nix-store-closure
RUN cp -R $(nix-store -qR result/) /tmp/nix-store-closure

FROM alpine:3.20

RUN apk update && \
    apk add --no-cache sudo shadow && \
    addgroup -g 1024 -S valnix && \
    adduser -u 1024 -G valnix -S -D valnix && \
    echo 'valnix ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers

COPY --from=BUILDER /tmp/nix-store-closure /nix/store
COPY --from=BUILDER /valkyrie/result /app

USER valnix

COPY .env /valkyrie/odin/.env

WORKDIR /valkyrie/odin
ENTRYPOINT ["/app/bin/odin"]