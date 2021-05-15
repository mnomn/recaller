package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func routeTraffic(path string, body string) {

	var routed int
	if Config.Debug > 0 {
		fmt.Printf("Incomming path: %v\n", path)
	}
	for _, route := range Config.Routes {
		normalizeIn(&path)
		normalizeIn(&route.In)
		if strings.Index(path, route.In) == 0 {

			if len(route.Out) == 0 {
				fmt.Printf("In {} has no out specified")
				continue
			}

			if Config.Debug > 0 || route.Debug > 0 {
				fmt.Printf("Path %v configured\n", path)
			}

			newBody := TransformBody(body, route)

			opList, ok := oldPostsLists[route.In]
			fmt.Printf("OLD: %v, %v \n", opList, ok)

			if strings.HasPrefix(route.Out, "mqtt") {
				pubMqtt(string(newBody), route)
				routed += 1
				return
			}

			if strings.HasPrefix(route.Out, "http") {
				sendHttp(string(newBody), route)
				routed += 1
				return
			}
		}
	}
	if routed == 0 {
		if Config.Debug > 0 {
			fmt.Printf("URL %v not routed.\n", path)
		}
	}
}

func sendHttp(postString string, route Route) {
	// HTTP Post is default protocol
	fmt.Printf("Route %v to %v using http POST \n", route.In, route.Out[0:20])
	req, err := http.NewRequest("POST", route.Out, strings.NewReader(postString))

	if len(route.HeaderKey) > 0 && len(route.HeaderValue) > 0 {
		req.Header.Set(route.HeaderKey, route.HeaderValue)
	}
	req.Header.Set("Content-Type", "application/json")

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
		panic(err)
	}
	logCall(route.In, route.Out, resp.Status)
	fmt.Printf("Send http to %v\n", route.Out)
	defer resp.Body.Close()

}

func pubMqtt(postString string, route Route) {

	outUrl, _ := url.Parse(route.Out)
	out := route.Out
	if outUrl.Scheme == "mqtt" {
		// TODO: Handle ports.
		out = strings.Replace(out, "mqtt", "tcp", 1)
	}

	opts := MQTT.NewClientOptions().AddBroker(out)

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

	logCall(route.In, topic, "")
	fmt.Printf("Route to mqtt broker %v, topic %v\n", out, topic)

	cid := "ClientID"
	opts.SetClientID(cid)
	if route.RoootCaFile != "" || route.CertFile != "" || route.PrivateKeyFile != "" {
		if Config.Debug > 1 || route.Debug > 1 {
			fmt.Printf("Mqtt cert files%v\n%v\n%v\n", route.RoootCaFile, route.CertFile, route.PrivateKeyFile)
		}
		tlsConf, err := makeTlsConfig(route.RoootCaFile, route.CertFile, route.PrivateKeyFile)
		fmt.Printf("TLS CFG ERR: %v\n", err)
		opts.SetTLSConfig(tlsConf)
		opts.SetCleanSession(true)
	}

	fmt.Println("pubMqtt New client")

	c := MQTT.NewClient(opts)
	fmt.Println("pubMqtt New client DUN")

	if Config.Debug > 2 || route.Debug > 2 {
		fmt.Printf("Mqtt data: %v\n", postString)
	}

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("pubMqtt connect ERROR. %v\n", token.Error())
		os.Exit(1)
	} else {
		fmt.Println("pubMqtt New client PUB")
		token = c.Publish(topic, 0, false, postString)
		fmt.Println("pubMqtt New client WW")
		token.Wait()
		fmt.Println("pubMqtt New client Waited")
		if token.Error() != nil {
			fmt.Printf("posted token err: %v\n", token.Error())
		}
	}
	c.Disconnect(250)
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
