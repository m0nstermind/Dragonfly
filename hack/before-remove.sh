#!/bin/bash

/usr/bin/systemctl stop dragonfly
/usr/bin/systemctl disable dragonfly

set -o nounset
set -o errexit
set -o pipefail

rm -f /usr/local/bin/dfget
rm -f /usr/local/bin/dfdaemon
