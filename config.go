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
		os.Exit(1);
	}

	for _, file := range files {
		name := file.Name()
		fmt.Println("Read config " + name)
		if !strings.HasSuffix(name, ".conf") {
			continue
		}

		var config map[string]interface{}

		//		fmt.Println("FILE: " +  file.Name())
		raw, err := ioutil.ReadFile(*confFlag + "/" + name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := json.Unmarshal(raw, &config); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		tmp, exist := config["address"]
		if exist {
			address = tmp.(string)
		}
		tmp, exist = config["main_username"]
		if exist {
			main_username = tmp.(string)
		}
		tmp, exist = config["main_password"]
		if exist {
			main_password = tmp.(string)
		}

		tmp, exist = config["routes"]
		if exist {
			r := tmp.([]interface{})
			for _, v := range r {
				routes = append(routes, v.(map[string]interface{}))
			}
		}
	}
	return //err
}

func readConfig() {
	// TODO: Support folder with many json files in.
	// Read input parameters
	confFlag := flag.String("conf", "", "Configuration directory, containing *.conf files. Default: ~/.route2cloud")
	flag.Parse()
	readConfigFiles(confFlag)
}
