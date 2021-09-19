package api

import (
	"bytes"
	"net/http"
)

// Do create request with parameter
func (f *Client) Do(urls string, body []byte, types string) (resp *http.Response, err error) {

	//create the final url
	urls, err = f.makeURL(urls)
	if err != nil {
		return nil, err
	}

	r := new(http.Request)

	if body != nil {
		r, err = http.NewRequest(types, urls, bytes.NewBuffer(body))
	} else {
		r, err = http.NewRequest(types, urls, nil)
	}

	if err != nil {
		return nil, err
	}

	for key, val := range f.Headers {
		r.Header.Set(key, val)
	}

	r = f.Parameter(r)
	return f.Client.Do(r)
}

// Post make Post request
func (f *Client) Post(path string, body interface{}) (resp *http.Response, err error) {

	bodys, err := f.Body(body)

	if err != nil {
		return nil, err
	}
	return f.Do(path, bodys, "POST")
}

// Get Create Get request
func (f *Client) Get(path string, body interface{}) (resp *http.Response, err error) {

	return f.Do(path, nil, "GET")
}

// Put create Put request
func (f *Client) Put(path string, body interface{}) (resp *http.Response, err error) {

	bodys, err := f.Body(body)

	if err != nil {
		return nil, err
	}
	return f.Do(path, bodys, "PUT")
}

// Delete create Delete request
func (f *Client) Delete(path string, body interface{}) (resp *http.Response, err error) {

	bodys, err := f.Body(body)

	if err != nil {
		return nil, err
	}

	return f.Do(path, bodys, "DELETE")
}
