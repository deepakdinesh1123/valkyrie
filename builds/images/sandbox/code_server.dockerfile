ARG BASE_IMAGE=odin_sandbox:0.0.1-ubuntu

FROM ${BASE_IMAGE}

RUN nix profile install nixpkgs#code-server

COPY hack/sandbox/start.sh /home/valnix/start.sh
COPY configs/sandbox/code-server.yaml /home/valnix/.config/code-server/config.yaml
COPY configs/sandbox/code-server-settings.json /home/valnix/.local/share/code-server/User/settings.json

ENTRYPOINT ["/bin/sh", "/home/valnix/start.sh"]
