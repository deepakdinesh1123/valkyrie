ARG BASE_IMAGE=odin_sandbox:ubuntu

FROM ${BASE_IMAGE}

# Copy the required files
COPY hack/sandbox/start.sh /home/valnix/start.sh
COPY configs/sandbox/code-server.yaml /home/valnix/.config/code-server/config.yaml

RUN nix-env -iA nixpkgs.code-server
RUN mkdir -p /home/valnix/.config/code-server

ENTRYPOINT ["/home/valnix/start.sh"]