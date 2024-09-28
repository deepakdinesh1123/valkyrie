FROM ubuntu:22.04

RUN apt update && \
    apt install -y adduser xz-utils curl && \
    groupadd -o -g 1000 -r vagrant && \
    adduser --uid 1000 --gid 1000 --disabled-password --gecos "" vagrant
RUN mkdir /etc/nix && echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf
USER vagrant
COPY hack/nix_setup.sh /home/vagrant/nix_setup.sh
ENV PATH="$PATH:/nix/store/2nhrwv91g6ycpyxvhmvc0xs8p92wp4bk-nix-2.24.9/bin"

CMD [ "bash", "/home/vagrant/nix_setup.sh" ]
