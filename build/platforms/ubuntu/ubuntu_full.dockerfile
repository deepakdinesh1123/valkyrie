FROM ubuntu:latest

RUN apt update && \
    apt install -y sudo adduser xz-utils curl && \
    groupadd -o -g 1024 -r valnix && \
    adduser --uid 1024 --gid 1024 --disabled-password --gecos "" valnix && \
    echo 'valnix ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers
RUN curl -L https://nixos.org/nix/install -o install_nix.sh
RUN chmod +x install_nix.sh
RUN mkdir /nix
RUN chown -R valnix /nix

USER valnix
RUN sh install_nix.sh --no-daemon
ENV PATH="{$PATH}:/home/valnix/.nix-profile/bin:/nix/var/nix/profiles/default/bin"

CMD [ "/bin/bash", "-c", "sleep infinity" ]
