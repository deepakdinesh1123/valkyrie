FROM nixos/nix:latest

RUN echo 'experimental-features = nix-command flakes' >> /etc/nix/nix.conf
RUN nix-env -iA cachix -f https://cachix.org/api/v1/install
RUN cachix use valkyrie
CMD ["bin/bash", "-c", "sleep infinity"]