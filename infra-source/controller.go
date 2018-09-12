package main

import (
	"net/http"
)

// Index page
func indexPage(w http.ResponseWriter, r *http.Request) {
	UserSession.Get(r, "mvc-user-session")
	View(w, r, nil, "index.html")
}

// About page
func aboutPage(w http.ResponseWriter, r *http.Request) {
	View(w, r, nil, "about.html")
}

// FAQ page
func faqPage(w http.ResponseWriter, r *http.Request) {
	View(w, r, nil, "faq.html")
}
