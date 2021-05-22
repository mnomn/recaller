# Route2cloud

A service that listens for http POST or PUT and re-sends mesasges to another url. The resend can append security and use a different schema: http, https or mqtt. It is also possible to add username and password, http headers or certificates. It is designed to let small devices (IoT) in a local network send its data to external servers which requires protocols and security not supported by the device.

## Example

```toml
username= "user1"
password="password1"

[[routes]]
in= "/test12"
out= "https://acme.org/measurements"
header= "ApiKey:SecretXYZ!"

[[Routes]]
in="/test11"
out="mqtt://localhost"
topic= "testdata"
username= "mqttUser"
password= "pass123"
```

Incomming http requests must use basic authentication with user1:password1 and use default port 8222.

A post to `http://touser1:password1@192.168.0.22:8222/test2` will be re-posted to with an extra header to https://acme.org/measurements with the same body.

A post to `http://touser1:password1@192.168.0.22:8222/test2` will be re-sent as mqtt to localhost. Mqtt login is mqttUser:pass123 and the topic will be "testdata".

## Build

- Clone the git and set up the go compiler (golang.com)
- Init the go module: `go mod init github.com/mnomn/route2cloud`, 
- Build: `go build`

### Cross compile

It is also possible to build for another target, like raspberry pi:  
`env GOOS=linux GOARCH=arm GOARM=5 go build -o route2cloud`

## Install and Configure

Build and copy binary to target machine.

For linux and raspberry pi

- Copy route2cloud and install_r2c.sh to target computer.
- Run `sudo bash install_r2c.sh`
- Add configuration file(s)"

### Configuration format

Configuration is defined in one or many files. Files must be in config folder (default usr/local/etc/route2cloud/). Configuration files can be called anything, as long as they end in ".conf". Toml and json is supported. See examples in configuration_files directory.

#### Top level configuration

http and and username/password. Default 8222 without password. Only set this in one place/file.

#### Routes configuration

The "routes" is a list of rules for how to resend incomming requests.

