#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

USE_DOCKER=${USE_DOCKER:-"0"}
VERSION=${VERSION:-"0.0.$(date +%s)"}
RPM_VERSION=${RPM_VERSION:-$VERSION}
RPM_RELEASE=${RPM_RELEASE:-"1"}

curDir=$(cd "$(dirname "$0")" && pwd)
cd "${curDir}" || return
BUILD_SOURCE_HOME=$(cd ".." && pwd)

. ./env.sh

BUILD_PATH=bin/${GOOS}_${GOARCH}
DFDAEMON_BINARY_NAME=dfdaemon
DFGET_BINARY_NAME=dfget

ODKL_DFD_PATH=odkl/${DFDAEMON_BINARY_NAME}

main() {
    cd "${BUILD_SOURCE_HOME}" || return
    # Maybe we should use a variable to set the directory for release,
    # however using a variable after `rm -rf` seems risky.
    mkdir -p release
    rm -rf release/*

    if [ "1" == "${USE_DOCKER}" ]
    then
        echo "Begin to package with docker."
        FPM="docker run --rm -v $(pwd):$(pwd) -w $(pwd) inoc603/fpm:alpine"
    else
        echo "Begin to package in local environment."
        FPM="fpm"
    fi

    case "${1-}" in
        rpm )
            build_rpm
            ;;
        deb )
            build_deb
            ;;
        * )
            build_rpm
            build_deb
            ;;
    esac
}

# TODO: Add description
DFCLIENT_DESCRIPTION="df-client"
# TODO: Add maintainer
MAINTAINER="dragonflyoss"

build_rpm() {
    ${FPM} -s dir -t rpm -f -p release --rpm-os=linux \
        --description "${DFCLIENT_DESCRIPTION}" \
        --maintainer "${MAINTAINER}" \
        --after-install ./hack/after-install.sh \
        --before-remove ./hack/before-remove.sh \
        --after-remove ./hack/after-remove.sh \
        -n ${RPM_NAME} -v "${RPM_VERSION}" \
        --iteration ${RPM_RELEASE} \
        --rpm-dist el7 \
        --config-files ${RPM_CONFIG_HOME}/dfdaemon.yml \
        --config-files ${RPM_SYSCONFIG_PATH} \
	"${BUILD_PATH}/${DFGET_BINARY_NAME}=${INSTALL_HOME}/${INSTALL_CLIENT_PATH}/${DFGET_BINARY_NAME}" \
	"${ODKL_DFD_PATH}/conf/dfdaemon.yml=${RPM_CONFIG_HOME}/dfdaemon.yml" \
	"${ODKL_DFD_PATH}/conf/dfdaemon-log.yml=${RPM_CONFIG_HOME}/dfdaemon-log.yml" \
	"${ODKL_DFD_PATH}/conf/dfget.yml=${RPM_CONFIG_HOME}/dfget.yml" \
	"${ODKL_DFD_PATH}/conf/dfget-log.yml=${RPM_CONFIG_HOME}/dfget-log.yml" \
	"${ODKL_DFD_PATH}/conf/sysconfig=${RPM_SYSCONFIG_PATH}" \
	"${ODKL_DFD_PATH}/systemd/dragonfly.service=${RPM_SYSTEMD_HOME}/dragonfly.service" \
	"${BUILD_PATH}/${DFDAEMON_BINARY_NAME}=${INSTALL_HOME}/${INSTALL_CLIENT_PATH}/${DFDAEMON_BINARY_NAME}"
}

build_deb() {
    ${FPM} -s dir -t deb -f -p release \
        --description "${DFCLIENT_DESCRIPTION}" \
        --maintainer "${MAINTAINER}" \
        --after-install ./hack/after-install.sh \
        --before-remove ./hack/before-remove.sh \
        -n df-client -v "${VERSION}" \
	"${BUILD_PATH}/${DFGET_BINARY_NAME}=${INSTALL_HOME}/${INSTALL_CLIENT_PATH}/${DFGET_BINARY_NAME}" \
	"${BUILD_PATH}/${DFDAEMON_BINARY_NAME}=${INSTALL_HOME}/${INSTALL_CLIENT_PATH}/${DFDAEMON_BINARY_NAME}"
}

main "$@"
