#!/usr/bin/env sh

set -e
. ~/.profile
chmod +x ~/valkyrie/exec.sh
cd ~/valkyrie
exec ./exec.sh
