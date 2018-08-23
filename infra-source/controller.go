package main

import (
	"net/http"
)

func indexPage(w http.ResponseWriter, r *http.Request) {
	UserSession.Get(r, "mvc-user-session")
	View(w, r, nil, "index.html")
}
