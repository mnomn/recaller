package main

import "testing"

func TestConfigFolder(t *testing.T) {
	folder := "./test_conf"
	readConfigFiles(&folder)
	if routes == nil || len(routes) == 0 {
		t.Errorf("Cannot handle config folder\n")
	}
	if len(main_username) < 1 {
		t.Error("Cannot read from config folder\n")
	}

	// Routes contain config from two files
	foundA := false
	foundB := false
	for _, v := range routes {
		if v["in"] == "/testA" {
			foundA = true
		}
		if v["in"] == "/testB" {
			foundB = true
		}
	}
	if !foundA || !foundB {
		t.Error("Did not read all config files\n")
	}
}
