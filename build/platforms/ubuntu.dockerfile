FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

ARG HOST_UID
ARG HOST_GID
ARG HOST_USER
ARG HOST_GROUP

RUN apt update && \
    apt install -y adduser xz-utils curl ca-certificates git vim

RUN groupadd -o -g $HOST_GID -r $HOST_GROUP && \
    adduser --uid $HOST_UID --gid $HOST_GID --disabled-password --gecos "" $HOST_USER

RUN mkdir /etc/nix && echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf
USER $HOST_USER
RUN mkdir ~/odin && chown $HOST_USER:$HOST_GROUP ~/odin
RUN cd /tmp && git clone --depth 1 --branch 24.05 --single-branch https://github.com/NixOS/nixpkgs.git 

COPY hack/nix_setup.sh /home/$HOST_USER/nix_setup.sh
COPY hack/nix_run.sh /home/$HOST_USER/nix_run.sh
VOLUME [ "/home/$HOST_USER" ]
WORKDIR /home/$HOST_USER/

CMD [ "/bin/bash", "nix_setup.sh" ]
