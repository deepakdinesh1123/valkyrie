#!/usr/bin/env bash

set -e

source ~/.profile
chmod +x ~/odin/exec.sh

while [ ! -f ~/status.txt ]; do sleep 1; done

cd ~/odin
exec ./exec.sh