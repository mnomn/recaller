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

func readConfigFiles(confFlag *string) (err error) {
	err = nil
	if len(*confFlag) < 1 {
		*confFlag, _ = os.Getwd()
		fmt.Printf("Using default config dir: %s\n", *confFlag)
	}

	files, err := ioutil.ReadDir(*confFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range files {
		var thisConfig RootConfig

		// Read all files ending with ".conf"
		shortName := file.Name()
		if !strings.HasSuffix(shortName, ".conf") {
			continue
		}
		fullName := *confFlag + "/" + shortName

		fileBytes, err := ioutil.ReadFile(fullName)
		if err != nil {
			fmt.Printf("Failed to read config file %v ", fullName)
			continue
		}

		// First try to read as TOML ...
		dec := json.NewDecoder(toml.New(bytes.NewBuffer(fileBytes)))
		tomlErr := dec.Decode(&thisConfig)

		if tomlErr != nil {
			// ... then try read as JSON
			if err := json.Unmarshal(fileBytes, &thisConfig); err != nil {
				fmt.Printf("Failed to parse %v\n", fullName)
				continue
			}
		}

		fmt.Println("Read " + fullName)

		updateGlobalValues(thisConfig)
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
}

func readConfig() {
	confFlag := flag.String("c", "", "Configuration directory, containing *.conf files. Default: Current directory.")
	measureTime := flag.Bool("t", false, "Measure time, Default: false.")
	flag.Parse()
	MeasureTime = *measureTime
	fmt.Printf("Configuration directory: %v, MeasureTime: %v\n", *confFlag, MeasureTime)
	readConfigFiles(confFlag)
}
