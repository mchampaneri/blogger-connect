package main

import (
	"net/http"

	"github.com/asdine/storm"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	sciter "github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/window"
	blogger "google.golang.org/api/blogger/v3"
)

// We may need to use window on other
// places as well to modify dom
// so it has to be GLOBAL
var AppWindow *window.Window
var AppWindowErr error

// Root elements is prime element to
// manipulate dom [ inserting / deleting / updating dom elements]
// It needed to initalize once and can be used
// as singlton
var RootElement *sciter.Element
var RootSelectorErr error

// WebServer and routes are going to be
// gobal as it need to be used anywhere/ every where
var WebRouter mux.Router

// Database need to be access through
// appliaction flow
// so we are making handle of db global
var Db *storm.DB
var DbErr error

// UserSession is store which containes
// user tokens, It is required on every
// api call, so making it global
var CurrentSession BloggerUserSession

// Preparing http and
// blogger clients for communication
var Client *http.Client
var BloggerClient *blogger.Service
var BloggerClientErr error

// initlizing gui for application
func init() {

	BloggerClient = nil
	// loading and verifying database condition
	initDb()
	// preparing sciter for exection
	initSciter()

}

func main() {

	// Preparing BloggerUserSession
	// If AccessToken is expiried already and
	// we dont' have any refresh token then
	// user has to authenticate him/her self via
	// goolge user consent

	// if !initUserSession() {
	// color.Yellow("Authenticatign user")
	authentication()
	// }

	// Loop back opens 20196 " blogger - client port "
	// For api and Oauth calls
	go loopback()

	color.Green("Blogger Desktop client initlized")
	// Show winodw
	for BloggerClient == nil {
		// Just wait for Bloggger Client to get initlaized
	}
	AppWindow.Show()
	// Making window runnign
	AppWindow.Run()

}
