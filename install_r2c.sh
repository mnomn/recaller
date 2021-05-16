#!/bin/bash

echo "Install on linux/raspberry pi"

# Installation base path
INSTALL_PATH=/usr/local

# Target name
PROG_NAME=route2cloud

SYSTEMD_PATH=/etc/systemd/system/

###########################
# Move binary install dir #
mv ${PROG_NAME} ${INSTALL_PATH}/bin/${PROG_NAME}

[ -f ${SYSTEMD_PATH}/${PROG_NAME}.service ] || {
    echo "Install systemd file ${SYSTEMD_PATH}/${PROG_NAME}.service"

    echo "
[Unit]
Description=route2cloud
After=network.target

[Service]
ExecStart=${INSTALL_PATH}/bin/${PROG_NAME} -c ${INSTALL_PATH}/etc/${PROG_NAME}

[Install]
WantedBy=multi-user.target
" > ${SYSTEMD_PATH}/${PROG_NAME}.service
}

# Make sure there is a config dir
mkdir -p ${INSTALL_PATH}/etc/${PROG_NAME}

echo "${PROG_NAME} installed"
