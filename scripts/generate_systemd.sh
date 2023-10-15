#!/bin/bash

UU=$USER
if [ $1 ]
then
    UU=$1
fi

filename=recaller@$UU.service

echo "Generating $filename."

echo "[Unit]
Description=recaller
After=network.target

[Service]
ExecStart=/usr/local/bin/recaller

[Install]
WantedBy=multi-user.target

" > $filename

