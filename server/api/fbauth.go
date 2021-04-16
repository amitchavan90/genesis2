package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var (
	fbOauthConfig = &oauth2.Config{
		ClientID:     "YOUR_CLIENT_ID",
		ClientSecret: "YOUR_CLIENT_SECRET",
		RedirectURL:  "YOUR_REDIRECT_URL_CALLBACK",
		Scopes:       []string{"public_profile"},
		Endpoint:     facebook.Endpoint,
	}
	fbOauthStateString = "pseudo-random"
)

const htmlIndex = `<html><body>
Logged in with <a href="/api/auth/facebook_auth/login">facebook</a>
</body></html>
`

func (c *AuthController) handleFacebookMain() func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlIndex))
	}

	return fn
}

func (c *AuthController) handleFacebookLogin() func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		Url, err := url.Parse(fbOauthConfig.Endpoint.AuthURL)
		if err != nil {
			log.Fatal("Parse: ", err)
		}
		parameters := url.Values{}
		parameters.Add("client_id", fbOauthConfig.ClientID)
		parameters.Add("scope", strings.Join(fbOauthConfig.Scopes, " "))
		parameters.Add("redirect_uri", fbOauthConfig.RedirectURL)
		parameters.Add("response_type", "code")
		parameters.Add("state", fbOauthStateString)
		Url.RawQuery = parameters.Encode()
		url := Url.String()
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}

	return fn
}

func (c *AuthController) handleFacebookCallback() func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")
		if state != fbOauthStateString {
			fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", fbOauthStateString, state)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		code := r.FormValue("code")

		token, err := fbOauthConfig.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Printf("fbOauthConfig.Exchange() failed with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		resp, err := http.Get("https://graph.facebook.com/me?access_token=" +
			url.QueryEscape(token.AccessToken))
		if err != nil {
			fmt.Printf("Get: %s\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		defer resp.Body.Close()

		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("ReadAll: %s\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		log.Printf("parseResponseBody: %s\n", string(response))

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	return fn
}
