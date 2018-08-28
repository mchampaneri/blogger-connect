package main

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

/* *0
 *  Registering the routes for the web
 */
func RegisterWebRoutes() {

	mainrouter := mux.NewRouter()

	/////////////////////////////////////////////////////////////////////////////////////////////
	////			 Opening the public directory for the open assets                        ////
	//// 			 Don't Change / Remove this line unless you know  		                 ////
	//// 			 what are you doing .					    		                     ////
	/////////////////////////////////////////////////////////////////////////////////////////////

	mainrouter.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	mainrouter.PathPrefix("/storage/").Handler(http.StripPrefix("/storage/", http.FileServer(http.Dir("./storage"))))

	/////////////////////////////////////////////////////////////////////////////////////////////
	////			 Define you mainrouter here							                         ////
	/////////////////////////////////////////////////////////////////////////////////////////////

	mainrouter.HandleFunc("/favicon.ico", faviconHandler)

	// Loading The Dynamic Routes
	dynamicRoutes(mainrouter)

	/////////////////////////////////////////////////////////////////////////////////////////////
	////			Opening Port 8085 (web) 					           ////
	/////////////////////////////////////////////////////////////////////////////////////////////

	mainrouter.NotFoundHandler = http.HandlerFunc(notFoundHandle)

	http.Handle("/", mainrouter)

	// logged_router := handlers.LoggingHandler(os.Stdout, mainrouter)
	if Config.Env == "dev" {
		http.ListenAndServe(Config.Port, handlers.CompressHandler(csrf.Protect([]byte("El0a6L8uqv"), csrf.Secure(false))(mainrouter)))
	} else if Config.Env == "prod" {
		http.ListenAndServe(Config.Port, handlers.CompressHandler(csrf.Protect([]byte("El0a6L8uqv"), csrf.Secure(true))(mainrouter)))
	}
}

func notFoundHandle(w http.ResponseWriter, r *http.Request) {
	View(w, r, nil, "404.html")
}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./favicon.ico")
}
