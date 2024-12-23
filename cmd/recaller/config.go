package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	toml "github.com/komkom/toml"
)

type RootConfig struct {
	Address        string          `json:"address"`
	Username       string          `json:"username"`
	Password       string          `json:"password"`
	Debug          int             `json:"debug"`
	Routes         []Route         `json:"routes"`
	Triggers       []Trigger       `json:"trigger"`
	Actions        []Action        `json:"action"`
	TriggerActions []TriggerAction `json:"trigger-actions"`
}

type Route struct {
	FileName       string   `json:"-"`
	In             string   `json:"in"`
	Out            string   `json:"out"`
	Topic          string   `json:"topic"`
	Username       string   `json:"username"`
	Password       string   `json:"password"`
	Method         string   `json:"method"`
	Headers        []string `json:"headers"`
	BodyTemplate   string   `json:"bodyTemplate"`
	PrivateKeyFile string   `json:"privateKeyFile"`
	CertFile       string   `json:"certFile"`
	RoootCaFile    string   `json:"rootCaFile"`
	Debug          int      `json:"debug"`
}

type Trigger struct {
	FileName string `json:"-"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Time     string `json:"time"`
}

type Action struct {
	FileName string `json:"-"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Method   string `json:"method"`
	Url      string `json:"url"`
}

type TriggerAction struct {
	Name    string   `json:"name"`
	Trigger string   `json:"trigger"`
	Actions []string `json:"actions"`
}

// Read one or many config fies and store here
var Config RootConfig

var MeasureTime bool

// Remove leading "/"
func normalizeInPath(in *string) {
	if len(*in) > 0 && (*in)[0] == '/' {
		s2 := (*in)[1:]
		*in = s2
	}
	*in = strings.ToLower(*in)
}

func readConfigFIle(fullName string) (RootConfig, error) {
	var thisConfig RootConfig
	var err error
	reader, err := os.Open(fullName)
	if err != nil {
		err = fmt.Errorf("failed to read config file %v", fullName)
		return thisConfig, err
	}

	// First try to read as TOML ...
	dec := json.NewDecoder(toml.New(reader))
	tomlErr := dec.Decode(&thisConfig)

	if tomlErr != nil {
		// ... then try read as JSON
		reader2, _ := os.Open(fullName)
		if err := json.NewDecoder(reader2).Decode(&thisConfig); err != nil {
			err = fmt.Errorf("failed to parse %v", fullName)
			return thisConfig, err
		}
	}

	fmt.Println("Read " + fullName)
	return thisConfig, nil
}

func readConfigFiles(confDir *string) (err error) {
	err = nil

	files, err := os.ReadDir(*confDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	confFiles := 0
	for _, file := range files {
		var thisConfig RootConfig

		// Read all files ending with ".conf"
		shortName := file.Name()
		if !strings.HasSuffix(shortName, ".conf") {
			continue
		}
		fullName := *confDir + "/" + shortName

		thisConfig, err := readConfigFIle(fullName)
		if err != nil {
			fmt.Printf("Failed to read config file %v err: %v", fullName, err)
			continue
		}

		confFiles++
		updateGlobalValues(thisConfig)
	}

	if confFiles == 0 {
		fmt.Println("There are no config files!")
	}

	return //err
}

func updateGlobalValues(configFromFile RootConfig) {
	// Only add if values are not set
	if configFromFile.Address != "" {
		Config.Address = configFromFile.Address
	}
	if configFromFile.Username != "" {
		Config.Username = configFromFile.Username
	}
	if configFromFile.Password != "" {
		Config.Password = configFromFile.Password
	}
	if configFromFile.Debug > Config.Debug {
		Config.Debug = configFromFile.Debug
	}

	if configFromFile.Routes == nil {
		return
	}

	for _, route := range configFromFile.Routes {
		rawIn := route.In
		normalizeInPath(&route.In)
		if route.In == "" {
			fmt.Printf("Route \"in\" missing, too short or invalid \"%v\"\n", rawIn)
			continue
		}

		Config.Routes = append(Config.Routes, route)
		fmt.Printf("  route %v -> %v\n", rawIn, route.Out)
	}

	Config.Triggers = append(Config.Triggers, configFromFile.Triggers...)
	Config.Actions = append(Config.Actions, configFromFile.Actions...)
	Config.TriggerActions = append(Config.TriggerActions, configFromFile.TriggerActions...)
}

func readConfig() {
	thisDir, _ := os.Getwd()
	confDir := flag.String("c", thisDir, "Configuration directory, containing *.conf files. Default: Current directory.")
	measureTime := flag.Bool("t", false, "Measure time, Default: false.")
	flag.Parse()
	MeasureTime = *measureTime
	fmt.Printf("Configuration directory: %v, MeasureTime: %v\n", *confDir, MeasureTime)
	readConfigFiles(confDir)
}
