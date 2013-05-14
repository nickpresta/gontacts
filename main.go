package main

import (
	"fmt"
	"flag"
	"html/template"
	"net/http"
	"io"

	"code.google.com/p/goauth2/oauth"
)

var (
	clientId = flag.String("id", "", "Client ID")
	clientSecret = flag.String("secret", "", "Client Secret")
	apiURL = flag.String("api", "https://www.google.com/m8/feeds", "API URL")
	requestURL = flag.String("request", "https://www.google.com/m8/feeds/contacts/default/full/?max-results=10000&alt=json", "API request")
	code = flag.String("code", "", "Authorization Code")
	cachefile = flag.String("cache", "cache.json", "Token cache file")
	port = flag.Int("port", 8080, "Webserver port")
)

var config *oauth.Config
var templates = template.Must(template.ParseFiles("index.html"))

func main() {
	flag.Parse()

	// Set up a configuration.
	config = &oauth.Config{
		ClientId: *clientId,
		ClientSecret: *clientSecret,
		Scope: *apiURL,
		AuthURL: "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://accounts.google.com/o/oauth2/token",
		TokenCache: oauth.CacheFile(*cachefile),
		RedirectURL: "http://localhost:8080/oauth2callback",
	}

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/auth", handleAuth)
	http.HandleFunc("/oauth2callback", handleOAuth2Callback)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", "")
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	url := config.AuthCodeURL("")
	http.Redirect(w, r, url, http.StatusFound)
}

func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	t := &oauth.Transport{Config: config}
	t.Exchange(code)

	resp, _ := t.Client().Get(*requestURL)
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}
