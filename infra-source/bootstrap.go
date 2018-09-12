package main

import (
	"encoding/json"
	"os"

	"github.com/fatih/color"
)

// Function that read the config.json file and populates
//  the Config singleton to use further in the app during
//  runtime.
//
func loadAppConfig() {
	//sync.Once{}.Do()
	//current_dir, _ := os.Getwd()
	configFile, _ := os.Open("./config/app.json")
	loadStaticPages, _ := os.Open("./config/static.json")

	defer func() {
		configFile.Close()
		loadStaticPages.Close()
	}()

	configFileParser := json.NewDecoder(configFile)
	configFileParser.Decode(&Config)

	staticsFileParser := json.NewDecoder(loadStaticPages)
	staticsFileParser.Decode(&StaticPages)

	color.Green(" * Configurations Loaded SuccessFully ")
}
