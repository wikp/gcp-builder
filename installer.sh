#!/bin/bash
export INSTALLER_BUILDER_BIN=`pwd`/bin
export INSTALLER_BUILDER_VERSION="0.1.31"
export INSTALLER_BUILDER_NAME="gcp-builder"
export FETCH_VERSION="0.1.1"

if [[ "$OSTYPE" == "linux-gnu" ]]; then
    export SYSTEM_TYPE="linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    export SYSTEM_TYPE="darwin"
else
   echo "Linux/MacOS is only supported"
   exit 2
fi

echo "Detected OS: $SYSTEM_TYPE"

export FETCH_DOWNLOAD_URL="https://github.com/gruntwork-io/fetch/releases/download/v${FETCH_VERSION}/fetch_${SYSTEM_TYPE}_amd64"

[ -d ${INSTALLER_BUILDER_BIN} ] || (echo "Creating ${INSTALLER_BUILDER_BIN} directory" && mkdir ${INSTALLER_BUILDER_BIN})

echo "Installing fetch $FETCH_DOWNLOAD_URL..."

(curl --progress-bar -L ${FETCH_DOWNLOAD_URL} -o ${INSTALLER_BUILDER_BIN}/fetch && \
    chmod +x ${INSTALLER_BUILDER_BIN}/fetch) || (echo "Fetch installation failed" && exit 2)

echo "Installing ${INSTALLER_BUILDER_NAME} version ${INSTALLER_BUILDER_VERSION}..."

(${INSTALLER_BUILDER_BIN}/fetch --repo="https://github.com/wendigo/${INSTALLER_BUILDER_NAME}" \
    --tag="v$INSTALLER_BUILDER_VERSION" \
    --release-asset="${INSTALLER_BUILDER_NAME}_${INSTALLER_BUILDER_VERSION}_${SYSTEM_TYPE}_x86_64.tar.gz" ${INSTALLER_BUILDER_BIN} && \
    cd ${INSTALLER_BUILDER_BIN} && tar -zxf "${INSTALLER_BUILDER_NAME}_${INSTALLER_BUILDER_VERSION}_${SYSTEM_TYPE}_x86_64.tar.gz" && \
    chmod +x ${INSTALLER_BUILDER_BIN}/${INSTALLER_BUILDER_NAME}) || (echo "${INSTALLER_BUILDER_NAME} installation failed" && exit 2)

echo "All tools installed, ready to rumble"

export PATH=$PATH:${INSTALLER_BUILDER_BIN}



