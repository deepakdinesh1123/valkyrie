FROM ubuntu:latest

RUN mkdir -p /var/run/sshd
RUN apt update && \
    apt install -y sudo adduser xz-utils curl openssh-server && \
    groupadd -o -g 1024 -r valnix && \
    adduser --uid 1024 --gid 1024 --disabled-password --gecos "" valnix && \
    echo 'valnix ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers
RUN curl -L https://nixos.org/nix/install -o install_nix.sh
RUN chmod +x install_nix.sh
RUN mkdir /nix
RUN chown -R valnix /nix

RUN echo "PasswordAuthentication no" >> /etc/ssh/sshd_config
RUN echo "PermitRootLogin no" >> /etc/ssh/sshd_config
RUN echo "AllowUsers valnix" >> /etc/ssh/sshd_config
RUN echo "PubkeyAuthentication yes" >> /etc/ssh/sshd_config
RUN echo "SyslogFacility AUTH" >> /etc/ssh/sshd_config


USER valnix
RUN sh install_nix.sh --no-daemon
ENV PATH="{$PATH}:/home/valnix/.nix-profile/bin:/nix/var/nix/profiles/default/bin"

COPY scripts/nix_env_init.sh .
CMD [ "/bin/sh", "nix_env_init.sh" ]
