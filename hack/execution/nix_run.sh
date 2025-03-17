#!/usr/bin/env sh

set -e
. ~/.profile
chmod +x ~/odin/exec.sh
cd ~/odin
exec ./exec.sh
