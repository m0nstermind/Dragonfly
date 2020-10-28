#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

ln -s /opt/dragonfly/df-client/dfget /usr/local/bin/dfget
ln -s /opt/dragonfly/df-client/dfdaemon /usr/local/bin/dfdaemon

# overwrite sysconfig from the package to sysconfig on the machine
mv -f /etc/sysconfig/dragonfly.rpmsave /etc/sysconfig/dragonfly >/dev/null 2>&1 || :

/usr/bin/systemctl daemon-reload
/usr/bin/systemctl enable dragonfly
