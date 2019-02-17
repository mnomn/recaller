# Route2cloud

A service that listens to http post and sends them to another url. It is designed for small devices (IoT) that cannot handle some protocols or the some type of security. It is possible to get incomming HTTP and send it on as HTTPS with certificates or user:passwords, or to an mqtt topic.


## Build
First do `go mod init github.com/mnomn/route2cloud`, 
then build using `go build` or use the make

There is a Makefile target for the raspberry pi too.
It will cross compile and install on the raspberry pi, including generating and installing a systemd service.
Starting and enabling the service must be done manually.

You must have a config file in your home directoryfor the program to be useful. See 'route2cloud.json' for an example.