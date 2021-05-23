package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	toml "github.com/komkom/toml"
)

type RootConfig struct {
	Address  string  `json:"address"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Debug    int     `json:"debug"`
	Routes   []Route `json:"routes"`
}

type Route struct {
	FileName       string `json:"-"`
	In             string `json:"in"`
	Out            string `json:"out"`
	Topic          string `json:"topic"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Method         string `json:"method"`
	Header         string `json:"header"`
	RegexpFind     string `json:"regexpFind"`
	RegexpReplace  string `json:"regexpReplace"`
	PrivateKeyFile string `json:"privateKeyFile"`
	CertFile       string `json:"certFile"`
	RoootCaFile    string `json:"rootCaFile"`
	Debug          int    `json:"debug"`
}

// Read one or many config fies and store here
var Config RootConfig

func readConfigFiles(confFlag *string) (err error) {
	err = nil
	if len(*confFlag) < 1 {
		*confFlag, _ = os.Getwd()
	}

	files, err := ioutil.ReadDir(*confFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range files {
		var thisConfig RootConfig

		shortName := file.Name()
		if !strings.HasSuffix(shortName, ".conf") {
			continue
		}
		fullName := *confFlag + "/" + shortName

		// Load both toml and json.
		// Try toml first
		fileBytes, err := ioutil.ReadFile(fullName)
		if err != nil {
			fmt.Printf("Failed to read config file %v ", fullName)
			continue
		}

		dec := json.NewDecoder(toml.New(bytes.NewBuffer(fileBytes)))
		tomlErr := dec.Decode(&thisConfig)

		if tomlErr != nil {
			if err := json.Unmarshal(fileBytes, &thisConfig); err != nil {
				fmt.Printf("Failed to parse %v\n", fullName)
				continue
			}
		}

		fmt.Println("Read config " + fullName)

		updateGlobalValues(thisConfig)
	}

	if Config.Debug > 0 {
		fmt.Println("Routes:")
		for _, r := range Config.Routes {
			fmt.Printf("  %v -> %v\n", r.In, r.Out)
		}
	}

	return //err
}

func updateGlobalValues(configFromFile RootConfig) {
	// Only overwrite if values are set
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
		readConfig()
	}

	for _, route := range configFromFile.Routes {
		Config.Routes = append(Config.Routes, route)
		if len(route.Header) > 0 {
			separator := strings.Index(route.Header, ":")
			if separator < 0 {
				fmt.Printf("Weader \"%v\" does not contain \":\"\n", route.Header)
			}
		}
	}
}

func readConfig() {
	confFlag := flag.String("c", "", "Configuration directory, containing *.conf files. Default: Current directory.")
	flag.Parse()
	readConfigFiles(confFlag)
}
