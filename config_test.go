package main

import (
	"testing"
)

func TestConfigFolder(t *testing.T) {
	folder := "./test_conf"
	readConfigFiles(&folder)
	if Config.Routes == nil {
		t.Errorf("Cannot handle config folder (nil)\n")
	}
	if Config.Routes == nil || len(Config.Routes) == 0 {
		t.Errorf("Cannot handle config folder\n")
	}
	if len(Config.Username) < 1 {
		t.Error("Cannot read from config folder\n")
	}

	// Routes contain config from two files
	foundA := false
	foundB := false
	for _, route := range Config.Routes {
		if route.In == "/testA" {
			foundA = true
		}
		if route.In == "/testB" {
			foundB = true
		}
	}
	if !foundA || !foundB {
		t.Error("Did not read all config files\n")
	}
}
