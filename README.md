# Recaller

Connect small devices (IoT) in a local network to servers which requires protocols and security not supported by the device.

This service listens to insecure http POST or PUT and re-sends messages to another host. The resend can use a different schema: http, https or mqtt. It is also possible to add security like username and password, http headers or certificates.

The incoming request body must be in JSON format. By default the body is resent unmodified in the outgoing request. It is posible to reformat the outgoing body with a template and overwrite Content-Type.

## Build

Bild with standard go tools.

Command line: `go build cmd/recaller`. To cross compile, add parameters: `env GOOS=linux GOARCH=arm GOARM=5 go build cmd/recaller -o recaller `.

It is also possible to use Makefile to build and install. For example `make install_pi`, which is a raspbery pi specific build and install target.

## Install

Copy the binary recaller to the target system and run it.

## Configuration

Configuration is defined in one or many files located in the config directory. Set config dir with "-d": `recaller -d /my/conf/dir`. All files with .conf will be red and they can be in toml or json format. See examples in configs directory or below.

### Top level configuration

`port`: Set which port incomming calls shall use. Default 8222.

`username` and `password`: Login needed to call this service, using "Basic Authentication". If omitted, no login is needed.

### Routes configuration

The "routes" is a list of rules for how to resend requests.

## Examples

```toml
# File: "recaller.conf"

username="user1"
password="password1"

[[routes]]
in = "/test1"
out = "https://acme.org/measurements"
headers = ["ApiKey:SecretXYZ!"]

[[routes]]
in = "/test2"
out = "mqtt://localhost"
topic = "testdata"
username = "mqttUser"
password = "pass123"

[[routes]]
in = "/test3"
out = "http://influx.myserver.com/api/v2/write?orgID=1111122222&bucket=bucket1"
headers = ["Content-Type:text/plain; charset=utf-8", "Authorization: Token abc123abc123abc123"]
bodyTemplate = "sensor_values,sensor_id={{.sensor}} temperature={{values.T}}"
```

With this config file, all incoming http requests must use basic authentication with user1:password1 and use default port 8222.

### Example 1

A post to `http://user1:password1@<ip>:8222/test1` will be re-posted with an extra header to `https://acme.org/measurements` with the same body.

### Example 2

A post to `http://user1:password1@<ip>:8222/test2` will be re-sent as mqtt to localhost. Mqtt login is mqttUser:pass123 and the topic will be "testdata".

### Example 3

A post to `http://user1:password1@<ip>:8222/test3` will be re-sent as plain text, not json. The bodyTemplate is used to create the outgoing body. For example incoming json `{"sensor":"S4","values":{"T":23.4,"unit":"C"}}` will be converted to text `sensor_values,sensor_id=S4 temperature=23.4"`.
If the incoming json does not fit the template, no message will be sent.
