ARG NIXPKGS_COMMIT=70c74b02eac46f4e4aa071e45a6189ce0f6d9265

FROM alpine/curl:8.14.1 AS builder

ARG NIXPKGS_COMMIT
RUN curl -L https://github.com/NixOS/nixpkgs/archive/${NIXPKGS_COMMIT}.tar.gz -o nixpkgs.tar.gz \
    && mkdir /nixpkgs \
    && tar -xzf nixpkgs.tar.gz -C /nixpkgs --strip-components=1

FROM deepakdinesh/nixery:ubuntu-25.05

COPY --from=builder /nixpkgs /tmp/nixpkgs

ENV NIXERY_PKGS_PATH=/tmp/nixpkgs
