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

FROM scratch
COPY --from=BUILDER /tmp/nix-store-closure /nix/store
COPY --from=BUILDER /valkyrie/result /app


COPY .env /valkyrie/odin/.env

WORKDIR /valkyrie/odin
ENTRYPOINT ["/app/bin/odin"]