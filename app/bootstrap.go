package main

import (
	"github.com/asdine/storm"
	"github.com/fatih/color"
	sciter "github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/window"
)

func initSciter() {
	// Defining size of rect we want for
	// client application
	rect := sciter.NewRect(0, 0, 800, 600)

	// Creatin window for client
	AppWindow, AppWindowErr = window.New(sciter.SW_MAIN|sciter.SW_CONTROLS|sciter.SW_RESIZEABLE|sciter.SW_ENABLE_DEBUG, rect)
	if AppWindowErr != nil {
		color.Red("Failed to create window for blogger client")
		return
	}

	// Fetching resource
	AppWindow.SetResourceArchive(resources)
	// Registering Callbacks
	AppWindow.DefineFunction("NavTo", NavTo)

	// Load Index Page [ main screen ]
	AppWindow.LoadFile("this://app/html/index.html")

	AppWindow.SetTitle("Bloggger Client")

}

func initDb() {

	// Opening blogger local database
	Db, DbErr = storm.Open("blogger.db")
	if DbErr != nil {
		color.Red("Failed to load database %s ", DbErr.Error())
		return
	}
}

func initUserSession() bool {
	// At this stage database will be already
	// initalied

	// reading last stored information about
	// user session
	sessionFetchErr := Db.One("Id", "0001", &CurrentSession)
	if sessionFetchErr != nil {
		color.Red("Can not fetch previous session %s", sessionFetchErr.Error())
		return false
	}

	return true

}
