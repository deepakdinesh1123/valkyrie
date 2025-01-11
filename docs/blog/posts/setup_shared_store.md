---
date: 2025-01-11
authors:
  - deepak
categories:
  - Tech
tags:
  - Nix
  - Podman
slug: setup_shared_store
description:
  - Setting up shared nix store between containers using podman and overlay
title: Setting Up Shared Nix Store
---

# Setting Up Shared Nix Store

**You can read about why and how we started using shared nix store [here](shared_nix_store.md)**

## Requirements:
- Podman
- Nix (Single user installation)

<!-- more -->

## Setup

- Create a base Containerfile

```docker
FROM alpine:3.20

RUN apk update && \
    apk add --no-cache shadow xz curl git vim && \
    addgroup -g 2048 -S valnix && \ # setup your group
    adduser -u 2048 -G valnix -s /bin/sh -D valnix # setup your user

RUN mkdir /etc/nix && echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf # add your nix config to nix.conf
USER valnix

WORKDIR /home/valnix/
COPY nix_setup.sh /home/valnix/nix_setup.sh

CMD [ "/bin/sh", "nix_setup.sh" ]
```

- Add this `nix_setup.sh` file in the same directory

```shell
#!/usr/bin/env sh

set -e

mkdir -p ~/.local/state/nix/profiles
ln -s $NIX_CHANNELS_ENVIRONMENT ~/.local/state/nix/profiles/channels
ln -s ~/.local/state/nix/profiles/channels ~/.local/state/nix/profiles/channels-1-link
ln -s $NIX_USER_ENVIRONMENT ~/.local/state/nix/profiles/profiles-1-link
ln -s ~/.local/state/nix/profiles/profiles-1-link ~/.local/state/nix/profiles/profile

mkdir ~/.nix-defexpr
ln -s ~/.local/state/nix/profiles/channels ~/.nix-defexpr/channels

ln -s $NIX_USER_ENVIRONMENT ~/.nix-profile

echo PATH="$PATH:~/.nix-profile/bin" >> ~/.profile
echo https://nixos.org/channels/nixos-24.11 nixpkgs >> ~/.nix-channels
. ~/.profile
```

- Run this command to build the image
```shell
podman build -t shared_nix:latest .
```

- You can now execute nix commands inside container and access all the packages that were present in your local nix store in the container
