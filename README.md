# Route2cloud

A service that listens for http POST or PUT and re-sends mesasges to another url. The resend can append security and use a different schema: http, https or mqtt. It is also possible to add username and password, http headers or certificates. It is designed for small devices (IoT) that cannot handle some protocols/schemas or some types of security.

## Build

- Clone the git and set up the go compiler (golang.com)
- Init the go module: `go mod init github.com/mnomn/route2cloud`, 
- Build: `go build`

### Cross compile

It is also possible to build for another target, like raspberry pi:  
`env GOOS=linux GOARCH=arm GOARM=5 go build -o route2cloud`

## Install and Configure

You must have a folder with conf files (default is current folder). See test_conf for some examples.

For linux and raspberry pi

- Copy route2cloud and install_r2c.sh to target computer.
- Run `sudo bash install_r2c.sh`
- Add configuration file(s) in "/usr/local/etc/route2cloud/"
