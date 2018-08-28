package main

import (
	"fmt"
	"net/http"

	"github.com/dghubble/gologin/google"
	oauth2Login "github.com/dghubble/gologin/oauth2"
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"
	blogger "google.golang.org/api/blogger/v3"
)

type SocialController struct{}

var Social SocialController

var GpConf *oauth2.Config

var Client *http.Client

var BloggerClient *blogger.Service

var BloggerClientErr error

func init() {

	GpConf = &oauth2.Config{
		ClientID:     "545090673153-8be24hfucftsa0a5mkoa58aqn416qlcq.apps.googleusercontent.com",
		ClientSecret: "81lDoXKUW8OU_RkE-PMx1nis",
		RedirectURL:  "http://localhost:8081/gp/callback",
		Scopes:       []string{"profile", "email", blogger.BloggerScope},
		Endpoint:     googleOAuth2.Endpoint,
	}
}

func (SocialController) GPissueSession() http.Handler {

	fn := func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		TokenToUse, _ := oauth2Login.TokenFromContext(ctx)

		gpuser, err := google.UserFromContext(ctx)

		usersession, errS := UserSession.Get(req, "mvc-user-session")
		if errS != nil {
			fmt.Println(errS.Error(), " during accessing the session")
		}
		usersession.Values["AccessToken"] = TokenToUse.AccessToken
		usersession.Values["RefreshToken"] = TokenToUse.RefreshToken
		usersession.Values["Id"] = gpuser.Id
		usersession.Save(req, w)

		if err != nil {
			fmt.Println("Error at issuing the token", err.Error())
			return
		}

		a := &oauth2.Config{}
		client := a.Client(ctx, TokenToUse)
		BloggerClient, BloggerClientErr = blogger.New(client)
		if BloggerClientErr != nil {
			fmt.Println(BloggerClientErr.Error())
			return
		}

		http.Redirect(w, req, "/blogs", http.StatusMovedPermanently)

	}
	return http.HandlerFunc(fn)
}
