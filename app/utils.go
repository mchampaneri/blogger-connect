package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		color.Red("Something went wrong %s", err.Error())
	}

}

func loopback() {
	router := mux.NewRouter()
	// fetchAndStoreTokens fetches tokens
	// updates user token inside Database
	// and reinitalizes user token with
	// updated value
	// It also updates BloggerClient with
	// latest token values
	router.HandleFunc("/callback", fetchAndStoreTokens)

	http.ListenAndServe(":20196", router)
}
