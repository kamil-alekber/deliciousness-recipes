package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/kamil-alekber/deliciousness-recipes/internal/models/tokens"
	"github.com/kamil-alekber/deliciousness-recipes/internal/models/users"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}
	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

var conf = &oauth2.Config{
	Endpoint:     google.Endpoint,
	ClientID:     cfg.googleClientID,
	ClientSecret: cfg.googleClientSecret,
	RedirectURL:  "http://localhost:8080/login/google/redirect",
	Scopes: []string{

		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
}

// should use dynamic state however google caches the value need to account for that later
func (app *application) generateStateOauthCookie(w http.ResponseWriter) string {
	var exp = time.Now().Add(20 * time.Minute)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	http.SetCookie(w, &http.Cookie{
		Name:    "oauth_state",
		Value:   state,
		Expires: exp,
		Path:    "/",
	})

	return state
}

func (app *application) loginGoogle(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, conf.AuthCodeURL("state"), http.StatusFound)
}

func (app *application) loginGoogleRedirect(w http.ResponseWriter, r *http.Request) {
	token, err := conf.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		fmt.Printf("error exchanging token: %s", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := conf.Client(context.Background(), token)
	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		fmt.Printf("error getting user info from google auth: %s", err.Error())
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading response body: %s", err.Error())
		return
	}

	var profile struct {
		ID         string `json:"id"`
		Email      string `json:"email"`
		Name       string `json:"name"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
		Picture    string `json:"picture"`
	}

	if err := json.Unmarshal(body, &profile); err != nil {
		fmt.Printf("error unmarshalling JSON: %s", err.Error())
		return
	}

	_, err = app.users.CreateUser(context.Background(), users.CreateUserParams{
		ID:         profile.ID,
		Email:      profile.Email,
		Name:       profile.Name,
		GivenName:  profile.GivenName,
		FamilyName: profile.FamilyName,
		Picture:    profile.Picture,
	})

	if err != nil {
		fmt.Printf("error creating user: %s", err.Error())
		return
	}

	_, err = app.tokens.CreateToken(context.Background(), tokens.CreateTokenParams{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
		ExpiresIn:    token.ExpiresIn,
		Vendor:       "google",
		UserID:       profile.ID,
	})

	if err != nil {
		fmt.Printf("error creating token: %s", err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "userid",
		Value:    profile.ID,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 365),
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
