package main

// Authentication of user
// all authentication methods
// are defined here .......

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	stringer "strings"
	"time"

	"github.com/fatih/color"
	"golang.org/x/oauth2"
	blogger "google.golang.org/api/blogger/v3"
)

// Loads users  consenet screen
// on platform preferd web browser
func authentication() {
	// Authentication against google Oauth2 ///
	url := "https://accounts.google.com/o/oauth2/v2/auth?scope=email%20profile%20https://www.googleapis.com/auth/blogger&response_type=code&redirect_uri=http://localhost:20196/callback&client_id=545090673153-8be24hfucftsa0a5mkoa58aqn416qlcq.apps.googleusercontent.com"
	openbrowser(url)
}

// Fetches tokens for authenticated user
// stores them in database for future use
// and prepares http client
func fetchAndStoreTokens(w http.ResponseWriter, r *http.Request) {

	// Preparing Oauth2 request url //
	data := url.Values{}
	data.Set("code", r.FormValue("code"))
	data.Add("client_id", "545090673153-8be24hfucftsa0a5mkoa58aqn416qlcq.apps.googleusercontent.com")
	data.Add("client_secret", "81lDoXKUW8OU_RkE-PMx1nis")
	data.Add("redirect_uri", "http://localhost:20196/callback")
	data.Add("grant_type", "authorization_code")

	// making request for user tokens //
	postLoginAuth, httpErr := http.Post("https://www.googleapis.com/oauth2/v4/token", "application/x-www-form-urlencoded",
		stringer.NewReader(data.Encode()))
	if httpErr != nil {
		color.Red("Http Request failed  %s ", httpErr.Error())
		return
	}

	// Reading tokens from response
	resp, respReadErr := ioutil.ReadAll(postLoginAuth.Body)
	if respReadErr != nil {
		color.Red("failed to read response %s ", respReadErr.Error())
		return
	}

	// Unloading tokens from response bod
	responseMap := make(map[string]interface{}, 6)
	unloadingErr := json.Unmarshal(resp, &responseMap)
	if unloadingErr != nil {
		color.Red("Failed to unload response data %s ", unloadingErr.Error())
		return
	}

	// Storing current value of user session
	// in database for later on user
	UserSessionInfoToStore := BloggerUserSession{
		Id:          "0001",
		AccessToken: responseMap["access_token"].(string),
		ExpeireIn:   responseMap["expires_in"].(float64),
		AssigendAt:  time.Now(),
	}

	if responseMap["refresh_token"] != nil {
		UserSessionInfoToStore.RefreshToken = responseMap["refresh_token"].(string)
	}

	TokenToUse := &oauth2.Token{
		AccessToken: responseMap["access_token"].(string),
		TokenType:   responseMap["token_type"].(string),
	}
	if responseMap["refresh_token"] != nil {
		TokenToUse.RefreshToken = responseMap["refresh_token"].(string)
	}

	prepareClient(TokenToUse)

	Db.Save(&UserSessionInfoToStore)
}

// Using provided token
// Prepares bloggerClient singlton ...
func prepareClient(TokenToUse *oauth2.Token) {
	ctx := context.Background()
	a := &oauth2.Config{}

	client := a.Client(ctx, TokenToUse)
	BloggerClient, BloggerClientErr = blogger.New(client)
	if BloggerClientErr != nil {
		fmt.Println(BloggerClientErr.Error())
		return
	}
}
