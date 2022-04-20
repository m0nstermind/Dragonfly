#!/bin/bash
# Copyright The Dragonfly Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o nounset
set -o errexit
set -o pipefail
set -x

export DOCKER_PLATFORM="--platform linux/amd64"
export INSTALL_HOME=/opt/dragonfly
export INSTALL_CLIENT_PATH=df-client
export INSTALL_SUPERNODE_PATH=df-supernode

# odkl rpm configuration
export RPM_CONFIG_HOME=/etc/dragonfly
export RPM_SYSCONFIG_PATH=/etc/sysconfig/dragonfly
export RPM_SYSTEMD_HOME=/usr/lib/systemd/system
export RPM_NAME=${RPM_NAME:-"df-client"}
# /odkl rpm configuration

export GO_SOURCE_EXCLUDES=( \
    "test" \
)

USE_DOCKER=${USE_DOCKER:-"0"}
BUILD_IMAGE=golang:1.13.15

if [[ "1" == "${USE_DOCKER}" ]]
then
	GOOS=$(docker run $DOCKER_PLATFORM --rm ${BUILD_IMAGE} go env GOOS | tr -dc '[:print:]' )
	GOARCH=$(docker run $DOCKER_PLATFORM --rm ${BUILD_IMAGE} go env GOARCH | tr -dc '[:print:]' )
else
	GOOS=$(go env GOOS)
	GOARCH=$(go env GOARCH)
fi
echo "GOOS=${GOOS}, GOPARCH=${GOARCH}"
export GOOS
export GOARCH
