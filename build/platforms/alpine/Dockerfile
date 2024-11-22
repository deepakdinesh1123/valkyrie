FROM alpine:3.20

RUN apk update && \
    apk add --no-cache shadow xz curl git vim && \
    addgroup -g 1024 -S valnix && \
    adduser -u 1024 -G valnix -S -D valnix

RUN curl -L https://nixos.org/nix/install -o install_nix.sh
RUN chmod +x install_nix.sh
RUN mkdir /nix
RUN chown -R valnix /nix
COPY configs/nix/nix.conf /etc/nix/nix.conf

USER valnix
RUN sh install_nix.sh --no-daemon
ENV PATH="$PATH:/home/valnix/.nix-profile/bin"

RUN nix-channel --remove nixpkgs
RUN nix-channel --add https://nixos.org/channels/nixos-24.05 nixpkgs
RUN nix-channel --update

RUN mkdir ~/odin

COPY configs/nix/packages.csv /tmp/packages.csv

COPY hack/execution/nix_run_docker.sh /home/valnix/nix_run.sh
COPY hack/execution/nix_stop.sh /home/valnix/nix_stop.sh
COPY hack/packages/store_packages_docker.sh /home/valnix/store_packages.sh
WORKDIR /home/valnix/

RUN sh store_packages.sh

RUN nix-env -if https://github.com/valkyrie-sh/cached-nix-shell/tarball/master

CMD [ "/bin/sh", "-c", "sleep infinity" ]