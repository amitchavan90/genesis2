package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/api/auth/google_auth/callback",
		ClientID:     "731932322217-fiu2t6peo6p87r95daua1814nta9ue09.apps.googleusercontent.com",
		ClientSecret: "vCuqGbPAeuK6LE_dAO5fanhz",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = "pseudo-random"
)

type GoogleProfile struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// handleGoogle logs a user in
func (c *AuthController) handleGoogleMain() func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var html = `<html><body><a href="/api/auth/google_auth/login">Google Login</a><body></html>`
		fmt.Fprint(w, html)
	}

	return fn
}

// handleGoogle logs a user in
func (c *AuthController) handleGoogleLogin() func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		url := googleOauthConfig.AuthCodeURL(oauthStateString)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}

	// TODO verify user and generate cookie

	return fn
}

// handleGoogle logs a user in
func (c *AuthController) handleGoogleCallback() func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.FormValue("code"))
		content, err := c.getGoogleUserInfo(r.FormValue("state"), r.FormValue("code"))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Fprintf(w, "Content: %s\n", content)
	}

	return fn
}

// getGoogleUserInfo gets user profile
func (c *AuthController) getGoogleUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	fmt.Printf("The code is: %v \n", code)

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	var userProfile GoogleProfile
	if err := json.Unmarshal(contents, &userProfile); err != nil {
		return nil, fmt.Errorf("Error when trying to unmarshal contents data: %s", err)
	}

	fmt.Printf("User Email: %v \n", userProfile.Email)

	return contents, nil
}

// facebookAuth logs a user in
func (c *AuthController) facebookAuth() func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Facebook Auth API"))
	}

	return fn
}
