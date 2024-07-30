FROM alpine:3.20

RUN apk update && \
    apk add --no-cache sudo shadow xz curl && \
    addgroup -g 1024 -S valnix && \
    adduser -u 1024 -G valnix -S -D valnix && \
    echo 'valnix ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers

RUN curl -L https://nixos.org/nix/install -o install_nix.sh
RUN chmod +x install_nix.sh
RUN mkdir /nix
RUN chown -R valnix /nix
RUN mkdir /etc/nix && echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf

USER valnix
RUN sh install_nix.sh --no-daemon
ENV PATH="{$PATH}:/home/valnix/.nix-profile/bin:/nix/var/nix/profiles/default/bin:/bin"

CMD [ "/bin/sh", "-c", "sleep infinity" ]