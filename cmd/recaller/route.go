package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func logResults(start time.Time, route Route, postString string, err error) {
	if err != nil {
		fmt.Printf("Failed to route to %s -> %s: %s\n", route.In, route.Out, err.Error())
		return
	}

	if Config.Debug > 2 || route.Debug > 2 {
		fmt.Printf("Post data: %v\n", postString)
	}

	elapsed := time.Since(start)
	fmt.Printf("Route to %s -> %s: %s\n", route.In, route.Out, elapsed)

}

func routeTraffic(path string, body string) {
	var routed int
	normalizeInPath(&path) // routes "In" are normalized during startup/config
	for _, route := range Config.Routes {
		if strings.HasPrefix(path, route.In) {
			start := time.Now()

			if len(route.Out) == 0 {
				fmt.Printf("Route In=%s has no out specified", route.In)
				continue
			}

			transformedBody, err := TransformBody(body, route)
			if transformedBody == "" || err != nil {
				fmt.Printf("Failed to transform body %v\n", body)
				continue
			}

			var sendError error

			if strings.HasPrefix(route.Out, "mqtt") {
				sendError = sendMqtt(string(transformedBody), route)
				routed += 1
			} else if strings.HasPrefix(route.Out, "http") {
				sendError = sendHttp(string(transformedBody), route)
				routed += 1
			} else {
				fmt.Printf("Unknown out schema %s. Use http or mqtt.", route.Out)
			}

			logResults(start, route, transformedBody, sendError)
		}
	}

	if routed == 0 {
		if Config.Debug > 0 {
			fmt.Printf("URL %v not routed.\n", path)
		}
	}
}
func sendHttp(postString string, route Route) error {
	// HTTP Post is default protocol

	method := "POST"
	if route.Method != "" {
		method = route.Method
	}
	req, err := http.NewRequest(method, route.Out, strings.NewReader(postString))

	contentTyprSet := false
	for _, header := range route.Headers {
		if separator := strings.Index(header, ":"); separator > 0 {
			req.Header.Set(header[:separator], header[separator+1:])
			if strings.Contains(header, "application/json") {
				contentTyprSet = true
			}
		}
	}

	if !contentTyprSet {
		req.Header.Set("Content-Type", "application/json")
	}

	var client *http.Client
	if route.RoootCaFile != "" && route.CertFile != "" && route.PrivateKeyFile != "" {
		if Config.Debug > 1 || route.Debug > 1 {
			fmt.Println("Found root_ca")
		}
		tc, _ := makeTlsConfig(route.RoootCaFile, route.CertFile, route.PrivateKeyFile)

		client = &http.Client{
			Transport: &http.Transport{TLSClientConfig: tc},
		}
	} else {
		client = &http.Client{}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func sendMqtt(postString string, route Route) error {

	serverString, err := convertToMqttServerString(route.Out)

	if err != nil {
		return err
	}

	opts := MQTT.NewClientOptions().AddBroker(serverString)

	// Default: same topic as input path
	topic := route.In

	// Optional parameters
	if len(route.Topic) > 0 {
		topic = route.Topic
	}
	if len(route.Username) > 0 && len(route.Password) > 0 {
		opts.SetUsername(route.Username)
		opts.SetPassword(route.Password)
	}

	cid := "ClientID"
	opts.SetClientID(cid)
	if route.RoootCaFile != "" || route.CertFile != "" || route.PrivateKeyFile != "" {
		if Config.Debug > 1 || route.Debug > 1 {
			fmt.Printf("Mqtt cert files%v\n%v\n%v\n", route.RoootCaFile, route.CertFile, route.PrivateKeyFile)
		}
		tlsConf, err := makeTlsConfig(route.RoootCaFile, route.CertFile, route.PrivateKeyFile)
		if err != nil {
			return err
		}
		opts.SetTLSConfig(tlsConf)
		opts.SetCleanSession(true)
	}

	c := MQTT.NewClient(opts)
	if c == nil {
		return fmt.Errorf("Cannot create MQTT client")
	}

	token := c.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	defer c.Disconnect(250)
	token = c.Publish(topic, 0, false, postString)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}

	return nil
}

// From https://github.com/manamanmana/aws-mqtt-chat-example
func makeTlsConfig(cafile, cert, key string) (*tls.Config, error) {
	var TLSConfig *tls.Config = &tls.Config{InsecureSkipVerify: false}

	var certPool *x509.CertPool
	var err error
	var tlsCert tls.Certificate
	if cafile != "" {
		certPool, err = getCertPool(cafile)
		if err != nil {
			return nil, err
		}
		TLSConfig.RootCAs = certPool
	}
	if cert != "" {
		certPool, err = getCertPool(cert)
		if err != nil {
			return nil, err
		}
		TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
		TLSConfig.ClientCAs = certPool
	}
	if key != "" {
		if cert == "" {
			return nil, fmt.Errorf("key specified but cert is not specified")
		}
		tlsCert, err = tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}
		TLSConfig.Certificates = []tls.Certificate{tlsCert}
	}
	return TLSConfig, nil
}

func getCertPool(pemPath string) (*x509.CertPool, error) {
	var certs *x509.CertPool = x509.NewCertPool()
	var pemData []byte
	var err error

	pemData, err = ioutil.ReadFile(pemPath)
	if err != nil {
		return nil, err
	}
	certs.AppendCertsFromPEM(pemData)
	return certs, nil
}

func convertToMqttServerString(configuredOut string) (string, error) {
	parsedUrl, err := url.Parse(configuredOut)

	if err != nil {
		return "", err
	}

	if parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		return "", fmt.Errorf("Invalid Out url: %v\n", configuredOut)
	}

	// paho does not understand mqtt or mqtts schema
	// convert "mqtt://hostname" to tcp://hostname:1883
	if strings.EqualFold(parsedUrl.Scheme, "mqtt") {
		port := ""
		if strings.Index(parsedUrl.Host, ":") < 0 {
			port = ":1883"
		}
		return fmt.Sprintf("tcp://%v%v", parsedUrl.Host, port), nil
	}

	if strings.EqualFold(parsedUrl.Scheme, "mqtts") {
		port := ""
		if strings.Index(parsedUrl.Host, ":") < 0 {
			port = ":8883"
		}
		return fmt.Sprintf("tcp://%v%v", parsedUrl.Host, port), nil
	}

	return configuredOut, nil

}
