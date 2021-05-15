package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

const defaultDir = ".route2cloud"

type RootConfig struct {
	Address  string  `json: "address"`
	Username string  `json: "username"`
	Password string  `json: "password"`
	Debug    int     `json: "debug"`
	Routes   []Route `json: "routes"`
}

type Route struct {
	FileName       string `json: "-"`
	In             string `json: "in"`
	Out            string `json: "out"`
	Topic          string `json: "topic"`
	Username       string `json: "username"`
	Password       string `json: "password"`
	HeaderKey      string `json: "headerKey"` // TODO: Combine to header key and value
	HeaderValue    string `json: "headerValue"`
	RegexpFind     string `json: "regexpFind"`
	RegexpReplace  string `json: "regexpReplace"`
	PrivateKeyFile string `json: "privateKeyFile"`
	CertFile       string `json: "certFile"`
	RoootCaFile    string `json: "rootCaFile"`
	Debug          int    `json: "debug"`
}

// Read one or many config fies and store here
var Config RootConfig

func readConfigFiles(confFlag *string) (err error) {
	err = nil
	if len(*confFlag) < 1 {
		u, _ := user.Current()
		*confFlag = u.HomeDir + "/" + defaultDir
	}
	fmt.Println("Read config files in " + *confFlag + "/")

	files, err := ioutil.ReadDir(*confFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range files {
		nameName := file.Name()
		if !strings.HasSuffix(nameName, ".conf") {
			continue
		}
		fmt.Println("Read config " + nameName)

		var config RootConfig

		raw, err := ioutil.ReadFile(*confFlag + "/" + nameName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := json.Unmarshal(raw, &config); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		readGlobalValues(config)

		if config.Routes == nil {
			continue
		}

		for _, route := range config.Routes {
			route.FileName = nameName
			Config.Routes = append(Config.Routes, route)
		}
	}

	if Config.Debug > 0 {
		fmt.Println("Routes:")
		for _, r := range Config.Routes {
			fmt.Printf("  %v -> %v\n", r.In, r.Out)
		}
	}

	return //err
}

func readGlobalValues(configFromFile RootConfig) {
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
}

func readConfig() {
	confFlag := flag.String("conf", "", "Configuration directory, containing *.conf files. Default: ~/.route2cloud")
	flag.Parse()
	fmt.Println("Generate main config")

	readConfigFiles(confFlag)
}
