FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

ARG HOST_UID
ARG HOST_GID
ARG HOST_USER
ARG HOST_GROUP

RUN apt update && \
    apt install -y adduser xz-utils curl

RUN groupadd -o -g $HOST_GID -r $HOST_GROUP && \
    adduser --uid $HOST_UID --gid $HOST_GID --disabled-password --gecos "" $HOST_USER

RUN mkdir /etc/nix && echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf
USER $HOST_USER
COPY hack/nix_setup.sh ~/nix_setup.sh
RUN mkdir ~/odin

CMD [ "bash", "~/nix_setup.sh" ]
