package main

import (
	"slices"
	"testing"
)

func TestConfigFolder(t *testing.T) {
	folder := "./../../configs"
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
	inA := "/testA"
	inB := "/testB"
	normalizeInPath((&inA))
	normalizeInPath((&inB))
	for _, route := range Config.Routes {
		if route.In == inA {
			foundA = true
		}
		if route.In == inB {
			foundB = true
		}
	}
	if !foundA || !foundB {
		t.Error("Did not read all config files\n")
	}
}

func TestTriggerActions(t *testing.T) {
	folder := "./../../test/testdata"
	readConfigFiles(&folder)

	const findTrigger = "atNoon"
	if slices.IndexFunc(Config.Triggers, func(t Trigger) bool { return t.Name == findTrigger }) < 0 {
		t.Errorf("Cannot find %v", findTrigger)
	}

	const findAction = "postOn"
	// Find atNoon trigger.
	if slices.IndexFunc(Config.Actions, func(a Action) bool { return a.Name == findAction }) < 0 {
		t.Errorf("Cannot find %v", findAction)
	}

	const findTriggerActions = "atNoonPostOn"
	ix := slices.IndexFunc(Config.TriggerActions, func(ta TriggerAction) bool { return ta.Name == findTriggerActions })

	if ix < 0 {
		t.Errorf("Cannot find %v", findTriggerActions)
	}

	ta := Config.TriggerActions[ix]
	if ta.Trigger != findTrigger {
		t.Errorf("Cannot find %v in %v", findTrigger, findTriggerActions)
	}

	if ta.Actions[0] != findAction {
		t.Errorf("Cannot find %v in %v", findAction, findTriggerActions)
	}
}
