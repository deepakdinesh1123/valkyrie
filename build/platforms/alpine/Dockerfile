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

RUN mkdir ~/odin
COPY hack/execution/* /home/valnix
WORKDIR /home/valnix/odin

CMD [ "/bin/sh", "-c", "sleep infinity" ]