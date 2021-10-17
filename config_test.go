package main

import (
	"testing"
)

func TestConfigFolder(t *testing.T) {
	folder := "./configuration_files"
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

		if route.Headers != nil {
			for _, hh := range route.Headers {
				t.Logf("H %v", hh)
			}
		} else {
			t.Log("No headers")
		}
	}
	if !foundA || !foundB {
		t.Error("Did not read all config files\n")
	}
}
