# Route2cloud

A service that listens for http post and sends them on to another url/protocol. It is designed for small devices (IoT) that cannot handle some protocols/schemas or some types of security. It is possible to get incomming HTTP and send it on as HTTPS with certificates or user:passwords, or to an mqtt topic.

## Build
First do `go mod init github.com/mnomn/route2cloud`, 
then build using `go build` or use `make`

There is a Makefile target for the raspberry pi too.
It will cross compile and install on the raspberry pi, including generating and installing a systemd service.
Starting and enabling the service must be done manually.

## Install and Config
You must have a config folder '~/.route2cloude with conf files. See 'route2cloud.conf' and test_conf for some examples.

On linux you can generate a systemd file to starting and stoping using systemctl.
