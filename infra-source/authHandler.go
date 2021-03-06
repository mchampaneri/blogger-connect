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

func gpLoad() {

	GpConf = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"profile", "email", blogger.BloggerScope},
		Endpoint:     googleOAuth2.Endpoint,
	}

	GpConf.RedirectURL = fmt.Sprintf("%s/gp/callback", Config.AppURL)
	fmt.Println("App url is ", GpConf.RedirectURL)
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
		usersession.Values["Valid"] = true
		usersession.Values["AccessToken"] = TokenToUse.AccessToken
		usersession.Values["RefreshToken"] = TokenToUse.RefreshToken
		usersession.Values["Id"] = gpuser.Id
		usersession.Values["Email"] = gpuser.Email
		usersession.Values["Name"] = gpuser.Name
		usersession.Values["Avatar"] = gpuser.Picture
		usersession.Save(req, w)

		if err != nil {
			fmt.Println("Error at issuing the token", err.Error())
			return
		}

		fmt.Println("Your access token is ", TokenToUse.AccessToken)

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
