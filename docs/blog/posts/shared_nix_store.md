---
date: 2025-01-10
authors:
  - deepak
categories:
  - Tech
tags:
  - Nix
  - Podman
slug: shared_nix_store
description:
  - Using Podman and overlays to share nix store between containers
title: Shared nix store
---

# Shared Nix Store

## Problem

When we started working on odin we realised that fetching packages from the official nix binary cache or
our own binary cache took a lot of time and CPU and for an execution engine it was simply not possible to justify the long setup time
but it was also not possible to just pre-install every dependency in the container since our primary goal was to support almost all
of the packages provided by nix and allow maximum customisability, pre-installing everything would mean that our container image size would be
more than 100GB at the very least. So we started expermenting on how we can share nix store between containers wihout compromising on anything.

<!-- more -->

## What we tried

Before getting to a solution that actually worked we tried a number of things that ultimately failed

### Volumes:

Our first thought was to use container volumes to share nix store between containers we created a container and volume mounted
our nix store expecting everything to work but we ran into a myriad of problems.

- Nix requires that the /nix directory is read writable by the user, which would mean that we cannot mount the volume as read only. We had to bind
mount it and pray that no one tries to delete the store
- There were also a lot of permission issues due to which we constantly received errors from nix

### Overlays:

Our next thought was to use an overlay so that changes made to the upper directory are not propogated to the actual nix store(lower directory),
At this point we were only working with docker and thought that read only mounting the nix store and creating an overlay would
solve all our problems but we soon ran into another long list of issues

- Docker does not allow creating an overlay inside of the container now and even if it did it needs root access to do so.
- We tried setting up the overlay outside of the container and then bind mounting it to the container but this was even more of a problem since the
upper directory was being written to the host and we ran into the same permission issue we were facing with volumes

### Local Overlay store

We found out that replit uses something known as local overlay store a feature provided by nix to easily share nix store between containers. We were
really excited and thanked the god for saving us but we ended up ditching this method as we could not set it up and we could not find enough docs
to solve the problems that we were having.


At this point we had almost given up and thought of abandoing the project altogether and work on something else but we decided to give it another
week before shutting it down, while searching for how to share folders between containers we found a stack overflow thread that saved our project

## Podman to the rescue

A StackOverflow thread mentioned that we can mount folders as overlay volumes in podman, which meant that we could mount a folder as an overlay volume
and anything written inside the container will not be written to the actual folder, rather it will be written to the upper directory. Even if a user
tries to delete all the files in the folder the actual files in the nix store would remain unharmed. This allowed us to share the nix store between podman containers without
worrying

**You can read about how to set up shared nix store using podman [here](setup_shared_store.md)**
