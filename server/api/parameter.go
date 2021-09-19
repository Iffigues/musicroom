package api

import (
	"net/http"
	"net/url"
	"path"
)

// makeURL create the url request to call
func (c *Client) makeURL(paths string) (urls string, err error) {
	u, err := url.Parse(c.host)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, paths)
	return u.String(), err
}

// SetParameter add parametre
func (c *Client) SetParameter(a, b string) {
	if c.Param == nil {
		c.Param = make(map[string]string)
	}
	c.Param[a] = b
}

// UnsetParameter unset one parameter
func (c *Client) UnsetParameter(b string) {
	if c.Param != nil {
		if _, ok := c.Param[b]; ok {
			delete(c.Param, b)
		}
	}
}

// UnsetAllParam delete all parameter
func (c *Client) UnsetAllParam() {
	if c.Param == nil {
		for key := range c.Param {
			delete(c.Param, key)
		}
	}
}

// Parameter add parameter in url example: "foo.html?tag=12"
func (c *Client) Parameter(u *http.Request) (r *http.Request) {

	if c.Param != nil {

		q := u.URL.Query()
		for key, val := range c.Param {
			q.Add(key, val)
		}

		u.URL.RawQuery = q.Encode()
	}

	return u
}

// SetHeader add value to the header
func (c *Client) SetHeader(a, b string) {
	if c.Headers == nil {
		c.Param = make(map[string]string)
	}
	c.Param[a] = b
}

// UnsetHeader delete on element of the header
func (c *Client) UnsetHeader(b string) {
	if c.Headers != nil {
		if _, ok := c.Headers[b]; ok {
			delete(c.Headers, b)
		}
	}
}

// UnsetAllHeader destroy header
func (c *Client) UnsetAllHeader() {
	if c.Headers == nil {
		for key := range c.Headers {
			delete(c.Headers, key)
		}
	}
}
