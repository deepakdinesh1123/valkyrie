#!/bin/sh

if ! [ -f /devpi/.nodeinfo ]; then
	devpi-init -c "$DEVPI_CONFIG"
fi

exec "$@"
