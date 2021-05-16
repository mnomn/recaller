# Route2cloud

A service that listens for http POST or PUT and re-sends mesasges to another url. The resend can append security and use a different schema: http, https or mqtt. It is also possible to add username and password, http headers or certificates. It is designed for small devices (IoT) that cannot handle some protocols/schemas or some types of security.

## Build
First do `go mod init github.com/mnomn/route2cloud`, 
then build using `go build` or use `make`

There is a Makefile target for the raspberry pi too.
It will cross compile and install on the raspberry pi, including generating and installing a systemd service.
Starting and enabling the service must be done manually.

## Install and Config
You must have a folder with conf files (default '~/.route2cloude'). See test_conf for some examples.

On linux you can generate a systemd file to starting and stoping using systemctl.
