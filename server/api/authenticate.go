package api

import (
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
)

// SetCode set Code string
func (c *Client) SetCode(code string) {
	c.Code = code
}

//SetState set state string
func (c *Client) SetState(state string) {
	c.State = state
}

// GetURL return the url to call for authenticate
func (c *Client) GetURL() (url string) {
	return c.Oauth.AuthCodeURL(c.State, c.AuthOpt...)

}

// Authentification ask autorisation api
func (c *Client) Authenticate(ctx context.Context) (err error) {

	var contexts context.Context
	if ctx == nil {
		contexts = context.Background()
	} else {
		contexts = ctx
	}

	if c.Types == 0 {
		c.Token, err = c.Oauth.PasswordCredentialsToken(context.Background(), c.Oauth.ClientID, c.Oauth.ClientSecret)
		if err != nil {
			return err
		}

		c.Client = c.Oauth.Client(contexts, c.Token)
	}
	if c.Types == 1 {
		c.Token, err = c.Oauth.Exchange(contexts, c.Code, c.AuthExange...)
		if err != nil {
			return err
		}
		c.Client = c.Oauth.Client(contexts, c.Token)
	}
	if c.Types == 2 {
		if c.Token, err = c.Legged.Token(contexts); err != nil {
			return err
		}
		c.Client = c.Legged.Client(contexts)
	}
	if c.Types == 3 {
		if c.Connector == nil {
			return errors.New("Connector is nil")
		}
		c.Connector.Connect(c)
	}
	return
}

// Body return []byte array for body send
func (c *Client) Body(payload interface{}) (body []byte, err error) {
	if payload != nil {
		body, err = json.Marshal(payload)
	}
	return
}

func (c *Client) Refresh(ctx context.Context) {
	var contexts context.Context
	if ctx == nil {
		contexts = context.Background()
	} else {
		contexts = ctx
	}
	token := new(oauth2.Token)
	token.AccessToken = c.Token.AccessToken
	c.Token = token
	c.Client = c.Oauth.Client(contexts, c.Token)
}
