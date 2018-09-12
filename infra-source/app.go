package main

import (
	"github.com/asaskevich/govalidator"
	"github.com/fatih/color"
)

func main() {

	BloggerClient = nil
	//Loading the Config
	loadAppConfig()

	// Initalizing google auth
	gpLoad()

	// Application Loaded Succesfully
	// Gretting Message
	color.Yellow(" * %s  ", Config.AppName)

	// Enforcing the goValidator over the models (Structs)
	govalidator.SetFieldsRequiredByDefault(false)
	DefaultLogger = Log{}

	// Loading the Routers for the web and api on their
	// mentioned ports
	RegisterWebRoutes()
}
