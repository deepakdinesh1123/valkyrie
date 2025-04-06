FROM python:3.12

RUN pip install devpi-server devpi-web

ENV DEVPI_CONFIG=/etc/devpi-config.yml

COPY hack/devpi/docker-entrypoint.sh /docker-entrypoint.sh
COPY hack/devpi/start-devpi-server.sh /start-devpi-server.sh
COPY configs/devpi-config.yml $DEVPI_CONFIG

VOLUME /devpi
WORKDIR /devpi

ENTRYPOINT ["sh", "/docker-entrypoint.sh"]
CMD ["sh", "/start-devpi-server.sh"]
