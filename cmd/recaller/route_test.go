package main

import (
	"strings"
	"testing"
)

func TestCreateServerString(t *testing.T) {

	var tcpStr = "tcp://abc.df:1883"

	var invalidStr = "test.com:8888"

	var mqttStr = "mqtt://abc.df"
	var mqttStrAfter = "tcp://abc.df:1883"

	var mqttsStr = "mqtts://abc.df"
	var mqttsStrAfter = "tcp://abc.df:8883"

	var mqttPortStr = "mqtt://abc.df:1876"
	var mqttPortStrAfter = "tcp://abc.df:1876"

	serverStr, err := convertToMqttServerString(tcpStr)
	if !strings.EqualFold(serverStr, tcpStr) {
		t.Errorf("Cannot use tcp url")
	}

	serverStr, err = convertToMqttServerString(invalidStr)
	if err == nil {
		t.Errorf("Expected error when using invalid out url")
	}

	///// MQTT /////
	serverStr, err = convertToMqttServerString(mqttStr)
	if !strings.EqualFold(serverStr, mqttStrAfter) {
		t.Errorf("Expected mqtt to be converted to tcp string")
	}

	///// MQTTS /////
	serverStr, err = convertToMqttServerString(mqttsStr)
	if !strings.EqualFold(serverStr, mqttsStrAfter) {
		t.Errorf("Expected mqtts and port to be transformed")
	}

	///// MQTT_PORT /////
	serverStr, err = convertToMqttServerString(mqttPortStr)
	if !strings.EqualFold(serverStr, mqttPortStrAfter) {
		t.Errorf("Expected mqtt to be converted to tcp server string")
	}

	///// TCP PORT /////
	serverStr, err = convertToMqttServerString(mqttPortStr)
	if !strings.EqualFold(serverStr, mqttPortStrAfter) {
		t.Errorf("Expected mqtt with port to be converted to tcp server string")
	}
}
