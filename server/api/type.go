package api

import (
	"encoding/json"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"net/url"
	//config "gitlab.com/wiserskills/v3/dataserver/config"
	//log "gitlab.com/wiserskills/v3/dataserver/log"
	"golang.org/x/oauth2"
)

// Oauth is the oauth value of client
type Oauth struct {
	AuthParam    map[string]string
	AuthExange   map[string]string
	RedirectURL  string
	ClientID     string
	QV           map[string]string
	AuthStyle    int
	Scopes       []string
	AuthURL      string
	TokenURL     string
	ClientSecret string
}

// Config represent the request
type Config struct {
	Types   int		  `json:"type"`
	Oauth   Oauth             `json:"oauth"`
	Host    string            `json:"url"`
	BaseURL string            `json:"baseurl"`
	Key     string            `json:"key"`
	Headers map[string]string `json:"headers"`
	Param   map[string]string `json:"param"`
	Body    []byte            `json:"body"`
	Code    string            `json:"code"`
	State   string            `json:"state"`
}

// Client is the request client
type Client struct {
	Oauth       *oauth2.Config
	Legged      *clientcredentials.Config
	QueryValues url.Values
	Types       int
	Code        string
	State       string
	AuthOpt     []oauth2.AuthCodeOption
	AuthExange  []oauth2.AuthCodeOption
	AuthStyle   int
	key         string
	host        string
	Connector   Connector
	Headers     map[string]string
	Param       map[string]string
	Token       *oauth2.Token
	Host        string
	Client      *http.Client
}

// New Config create APIConfig obj
func NewConfig(data []byte) (a *Config, err error) {
	a = &Config{}
	err = json.Unmarshal(data, a)
	return
}

// NewClient create new client obj
func (a *Config) NewClient() (c *Client, err error) {

	if _, err := url.ParseRequestURI(a.Host); err != nil {
		return nil, err
	}

	if a.Oauth.TokenURL != "" {
		if _, err := url.ParseRequestURI(a.Oauth.TokenURL); err != nil {
			return nil, err
		}
	}

	if a.Oauth.AuthURL != "" {
		if _, err := url.ParseRequestURI(a.Oauth.AuthURL); err != nil {
			return nil, err
		}
	}

	if _, err := url.ParseRequestURI(a.Host); err != nil {
		return nil, err
	}

	n := &Client{
		Types:     a.Types,
		host:      a.Host,
		AuthStyle: a.Oauth.AuthStyle,
		Headers:   a.Headers,
		Param:     a.Param,
		Code:      a.Code,
		State:     a.State,
		Oauth: &oauth2.Config{
			ClientID:     a.Oauth.ClientID,
			ClientSecret: a.Oauth.ClientSecret,
			RedirectURL:  a.Oauth.RedirectURL,
			Scopes:       a.Oauth.Scopes,
			Endpoint: oauth2.Endpoint{
				TokenURL:  a.Oauth.TokenURL,
				AuthURL:   a.Oauth.AuthURL,
				AuthStyle: oauth2.AuthStyle(a.Oauth.AuthStyle),
			},
		},
	}
	if a.Oauth.QV != nil && len(a.Oauth.QV) > 0 {
		params := url.Values{}
		for k, v := range a.Oauth.QV {
			params.Add(k, v)
		}
		n.QueryValues = params
	}
	n.Legged = &clientcredentials.Config{
		ClientID:       n.Oauth.ClientID,
		ClientSecret:   n.Oauth.ClientSecret,
		TokenURL:       n.Oauth.Endpoint.TokenURL,
		Scopes:         n.Oauth.Scopes,
		EndpointParams: n.QueryValues,
	}
	for key, val := range a.Oauth.AuthParam {
		n.AuthOpt = append(n.AuthOpt, oauth2.SetAuthURLParam(key, val))
	}

	for key, val := range a.Oauth.AuthExange {
		n.AuthExange = append(n.AuthExange, oauth2.SetAuthURLParam(key, val))
	}

	return n, nil
}

// AddConnector add interface to client obj
func (c *Client) AddConnector(connect Connector) {
	c.Connector = connect
}
