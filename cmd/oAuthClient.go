package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	oauth2gh "golang.org/x/oauth2/github"
)

var (
	// You must register the app at https://github.com/settings/developers
	// Set callback to http://127.0.0.1:9080/github/callback
	// Set ClientId and ClientSecret in ENV
	oauthConf = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID_ENV"),
		ClientSecret: os.Getenv("CLIENT_SEC_ENV"),
		Scopes:       []string{"user:email", "repo"},
		Endpoint:     oauth2gh.Endpoint,
	}

	// random string for oauth2 API calls to protect against CSRF
	oauthStateString    = "initialDefaultStr"
	oauthStateStringLen = 16
)

func authInit() string {
	// feed state for every new login
	oauthStateString = randString(oauthStateStringLen)
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)

	return url
}

func getAuthClient(respWriter http.ResponseWriter, req *http.Request) (*github.Client, context.Context) {
	state := req.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(respWriter, req, "/", http.StatusTemporaryRedirect)
		return nil, nil
	}

	ctx := context.Background()
	athCode := req.FormValue("code")
	token, err := oauthConf.Exchange(ctx, athCode)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(respWriter, req, "/", http.StatusTemporaryRedirect)
		return nil, nil
	}

	oauthClient := oauthConf.Client(ctx, token)
	client := github.NewClient(oauthClient)

	return client, ctx
}
