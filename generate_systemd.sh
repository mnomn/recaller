#!/bin/bash

UU=$USER
if [ $1 ]
then
    UU=$1
fi

filename=route2cloud@$UU.service

echo "Generating $filename."


echo "[Unit]
Description=route2cloud
After=network.target

[Service]
ExecStart=/usr/local/bin/route2cloud

[Install]
WantedBy=multi-user.target

" > $filename

