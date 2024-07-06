FROM nixos/nix:2.23.1 as BUILDER

COPY cmd /valkyrie/cmd
COPY build/package/nix/geri.nix /valkyrie/build/package/nix/geri.nix
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
    build .#geri

RUN mkdir /tmp/nix-store-closure
RUN cp -R $(nix-store -qR result/) /tmp/nix-store-closure

FROM nixos/nix:latest

COPY --from=BUILDER /tmp/nix-store-closure /nix/store
COPY --from=BUILDER /valkyrie/result/bin/geri /bin

RUN echo 'experimental-features = nix-command flakes' >> /etc/nix/nix.conf
RUN nix-channel --add https://nixos.org/channels/nixpkgs-unstable unstable
RUN nix-channel --update
RUN nix-env -iA devenv -f https://github.com/NixOS/nixpkgs/tarball/nixpkgs-unstable

COPY .env /valkyrie/geri/.env

WORKDIR /valkyrie/geri
ENTRYPOINT ["/bin/geri"]
CMD ["start"]
