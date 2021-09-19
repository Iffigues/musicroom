package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2/linkedin"
)

func SendJSON(ar interface{}, w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json;")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ar); err != nil {
	}
}

func serve(c *Client, t *testing.T) (ts *httptest.Server) {
	ts = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//http.Redirect(w, r, c.GetURL(), 301)
		if r.Method == "POST" {
			c.Types = 1
			fmt.Fprintln(w, c.GetURL())
		}
		if r.Method == "GET" {
			ctx := context.Background()
			c.Code = r.FormValue("code")
			c.Types = 1
			err := c.Authenticate(ctx)
			fmt.Fprintln(w, err)
			ee, err := c.Get("/v2/people/id:_jmbSTu8rX", nil)
			if err != nil {
				t.Fatal(err)
			}
			bodyBytes, err := ioutil.ReadAll(ee.Body)
			bodyString := string(bodyBytes)
			fmt.Fprintln(w, bodyString, err)
		}
	}))
	return
}

type Body struct {
	Data string `json:"body"`
}

func lauth() {
	var t *Oauth
	t.AuthURL = "https://www.linkedin.com/oauth/v2/authorization"
	t.TokenURL = "https://www.linkedin.com/oauth/v2/accessToken"
	fmt.Println(t)
}

/*
func TestApi(t *testing.T) {
	ap := &Config{
		Host: "https://jsonplaceholder.typicode.com",
	}
	b, err := json.Marshal(ap)

	if err != nil {
		t.Fatal(err)
	}
	logger := log.New(false)
	configuration, err := config.Load("./tests/config.yml", logger)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewConfig(logger, configuration, b)
	if err != nil {
		t.Fatal(err)
	}
	c := g.NewClient()
	a, err := c.Get("/posts/1", nil, nil)
	t.Log(a, err)
	c.request.Headers = map[string]string{
		"Content-type": "application/json; charset=UTF-8",
	}
	a, err = c.Post("/posts", Body{Data: "send"}, nil)
	t.Log(a, err)
}
*/

func TestLinkedin(t *testing.T) {
	ap := &Config{
		Host:  "https://api.linkedin.com/",
		Oauth: Oauth{},
		Headers: map[string]string{
			"grant_type":   "client_credentials",
			"Content-Type": "application/x-www-form-urlencoded",
		},
	}
	ap.Oauth.ClientID = "86ckga2ieqzs62"
	ap.Oauth.Scopes = []string{
		"r_emailaddress",
		"r_liteprofile",
		"w_member_social",
	}
	ap.Oauth.ClientSecret = "SRkUbejCxCLpRHUd"
	ap.Oauth.AuthParam = map[string]string{
		"grant_type": "client_credentials",
	}
	ap.Oauth.AuthExange = map[string]string{
		"redirect_uri": "http://1b6cab60.ngrok.io",
	}
	ap.Oauth.RedirectURL = "http://1b6cab60.ngrok.io"
	ap.Oauth.TokenURL = linkedin.Endpoint.TokenURL
	ap.Oauth.AuthURL = linkedin.Endpoint.AuthURL
	c, err := ap.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	ts := serve(c, t)
	l, err := net.Listen("tcp", "gopiko.fr:8080")
	if err != nil {
		log.Fatal(err)
	}
	ts.Listener = l
	t.Log("e")
	// Start the server.
	ts.Start()
	defer ts.Close()
	for {
	}

}

/*
func TestInsee(t *testing.T) {
	ap := &Config{
		Host:  "https://api.insee.fr/",
		Oauth: Oauth{},
		Headers: map[string]string{
			"grant_type": "client_credentials",
		},
	}
	ap.Oauth.ClientID = "220RaLiYS2CAG668ZS1UZCnAvDsa"
	ap.Oauth.ClientSecret = "LTlyv3uAossGgzNMmJ9JCZuKERIa"
	ap.Oauth.AuthExange = map[string]string{
		"grant_type": "client_credentials",
	}

	ap.Oauth.TokenURL = "https://api.insee.fr/token"
	ap.Oauth.AuthURL = "https://api.insee.fr/token"
	//ap.Oauth.RedirectURL = "http://localhost.com"
	c, err := ap.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	c.Types = 1
	err = c.Authenticate(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	res, err := c.Get("/entreprises/sirene/V3/siren/821838075", nil)
	if res.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		bodyString := string(bodyBytes)
		t.Log(bodyString)
	}
	t.Fatal(res, err)
}
*/
