FROM alpine:3.20

ARG HOST_UID
ARG HOST_GID
ARG HOST_USER
ARG HOST_GROUP

RUN apk update && \
    apk add --no-cache shadow xz curl && \
    addgroup -g $HOST_GID -S $HOST_GROUP && \
    adduser -u $HOST_UID -G $HOST_GROUP -s /bin/sh -D $HOST_USER

RUN mkdir /etc/nix && echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf
USER $HOST_USER
COPY hack/nix_setup.sh ~/nix_setup.sh
RUN mkdir ~/odin

CMD [ "bash", "~/nix_setup.sh" ]