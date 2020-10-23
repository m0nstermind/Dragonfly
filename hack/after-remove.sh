#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

/usr/bin/systemctl daemon-reload
rmdir --ignore-fail-on-non-empty /etc/dragonfly