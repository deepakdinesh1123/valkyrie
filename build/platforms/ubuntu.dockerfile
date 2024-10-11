FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update && \
    apt install -y adduser xz-utils curl ca-certificates git vim sudo
RUN groupadd -o -g 1024 -r odnix && \
    adduser --uid 1024 --gid 1024 --disabled-password --gecos "" odnix && \
    echo 'odnix ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers

USER odnix

WORKDIR /home/odnix/
RUN curl -L https://github.com/DavHau/nix-portable/releases/latest/download/nix-portable-$(uname -m) > ./nix-portable &&  \
    chmod +x ./nix-portable
RUN mkdir odin
ENV NP_GIT=/usr/bin/git

CMD ["/bin/bash", "-c", "sleep infinity"]