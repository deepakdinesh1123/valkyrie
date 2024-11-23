FROM codercom/code-server  

USER root
RUN apt update && \
    apt install -y adduser xz-utils curl && \
    groupadd -o -g 1024 -r valnix && \
    adduser --uid 1024 --gid 1024 --disabled-password --gecos "" valnix && \
    echo 'valnix ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers
RUN mkdir /nix
RUN chown -R valnix /nix
RUN mkdir -p /home/valnix/.config/process-compose && chown -R  valnix /nix
COPY start.sh /home/valnix/start.sh
RUN chmod +x /home/valnix/start.sh
RUN mkdir /etc/nix && echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf

USER valnix
WORKDIR /home/valnix
RUN curl -L https://nixos.org/nix/install -o install_nix.sh
RUN chmod +x install_nix.sh
RUN sh install_nix.sh --no-daemon
ENV PATH="$PATH:/home/valnix/.nix-profile/bin"

# RUN rm /home/valnix/.config/code-server/config.yaml
# COPY config.yaml  /home/valnix/.config/code-server/config.yaml
# COPY flake.nix /home/valnix/test/

ENTRYPOINT [ "/home/valnix/start.sh" ]