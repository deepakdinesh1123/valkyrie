FROM ubuntu:22.04

RUN apt update && \
    apt install -y adduser xz-utils curl && \
    groupadd -o -g 1024 -r valnix && \
    adduser --uid 1024 --gid 1024 --disabled-password --gecos "" valnix && \
    echo 'valnix ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers
RUN curl -L https://nixos.org/nix/install -o install_nix.sh
RUN chmod +x install_nix.sh
RUN mkdir /nix
RUN chown -R valnix /nix
RUN mkdir /etc/nix && echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf

USER valnix
RUN sh install_nix.sh --no-daemon
ENV PATH="$PATH:/home/valnix/.nix-profile/bin"

CMD [ "/bin/sh", "-c", "sleep infinity" ]