#!/bin/bash

echo "Install on linux/raspberry pi"

ARG1=$1

# Installation base path
INSTALL_PATH=/usr/local

# Target name
PROG_NAME=recaller

SYSTEMD_PATH=/etc/systemd/system/

###########################
# Move binary install dir #
mv ${PROG_NAME} ${INSTALL_PATH}/bin/${PROG_NAME}

[ -f ${SYSTEMD_PATH}/${PROG_NAME}.service ] || {
    echo "Install systemd file ${SYSTEMD_PATH}/${PROG_NAME}.service"

    echo "
[Unit]
Description=recaller
After=network.target

[Service]
ExecStart=${INSTALL_PATH}/bin/${PROG_NAME} -t -c ${INSTALL_PATH}/etc/${PROG_NAME}

[Install]
WantedBy=multi-user.target
" > ${SYSTEMD_PATH}/${PROG_NAME}.service

    systemctl daemon-reload
}

# Make sure there is a config dir
mkdir -p ${INSTALL_PATH}/etc/${PROG_NAME}

[ "$ARG1" != "-r" ] || {
    echo "Restart service after install"
    systemctl restart ${PROG_NAME}.service
}

echo "${PROG_NAME} installed"
