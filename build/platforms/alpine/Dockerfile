FROM alpine:3.20

ARG HOST_UID
ARG HOST_GID
ARG HOST_USER
ARG HOST_GROUP

RUN apk update && \
    apk add --no-cache shadow xz curl git vim && \
    addgroup -g ${HOST_GID} -S ${HOST_GROUP} && \
    adduser -u ${HOST_UID} -G ${HOST_GROUP} -s /bin/sh -D ${HOST_USER}

RUN mkdir /etc/nix && echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf
USER $HOST_USER
RUN mkdir ~/odin && chown $HOST_USER:$HOST_USER ~/odin
RUN cd /tmp && git clone --depth 1 --branch 24.05 --single-branch https://github.com/NixOS/nixpkgs.git 

COPY hack/execution/* /home/$HOST_USER/

VOLUME [ "/home/$HOST_USER" ]
WORKDIR /home/$HOST_USER/

CMD [ "/bin/sh", "nix_setup.sh" ]