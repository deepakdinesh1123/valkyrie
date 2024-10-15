#!/usr/bin/env sh

set -e
. ~/.profile
chmod +x ~/odin/exec.sh

while [ ! -f ~/status.txt ]; do sleep 1; done

cd ~/odin
exec ./exec.sh