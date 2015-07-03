package grestclient

import (
	"errors"
	"net/http"
	"net/url"
)

//cloneUrl will do a deep clone of the url
//you pass in
func cloneUrl(u *url.URL) *url.URL {
	var userInfo *url.Userinfo
	if u.User != nil {
		if p, ok := u.User.Password(); ok {
			userInfo = url.UserPassword(u.User.Username(), p)
		} else {
			userInfo = url.User(u.User.Username())
		}
	}

	return &url.URL{
		Scheme:   u.Scheme,
		Opaque:   u.Opaque,
		User:     userInfo,
		Host:     u.Host,
		Path:     u.Path,
		RawQuery: u.RawQuery,
		Fragment: u.Fragment,
	}
}

type client struct {
	base        *url.URL
	reqMutators []RequestMutator
	resMutators []ResponseMutator
	headers     http.Header
	query       url.Values
	client      *http.Client
	marshaler   MarshalerFunc
	unmarshaler UnmarshalerFunc
}

func (c *client) Headers() http.Header {
	if c.headers == nil {
		c.headers = make(http.Header)
	}
	return c.headers
}

func (c *client) SetHeaders(h http.Header) {
	c.headers = h
}

func (c *client) Query() url.Values {
	if c.query == nil {
		c.query = make(url.Values)
	}
	return c.query
}

func (c *client) SetQuery(q url.Values) {
	c.query = q
}

func (c *client) AddRequestMutators(rm ...RequestMutator) Client {
	c.reqMutators = append(c.reqMutators, rm...)
	return c
}

func (c *client) AddResponseMutators(rm ...ResponseMutator) Client {
	c.resMutators = append(c.resMutators, rm...)
	return c
}

func (c *client) SetRequestMutators(rm ...ResponseMutator) Client {
	c.resMutators = rm
	return c
}

func (c *client) SetResponseMutators(rm ...ResponseMutator) Client {
	c.resMutators = rm
	return c
}

func (c *client) RequestMutators() []RequestMutator {
	return c.reqMutators
}

func (c *client) ResponseMutators() []ResponseMutator {
	return c.resMutators
}

func (c *client) Get(path string, query url.Values, successResult interface{}, errorResult interface{}) (*http.Response, error) {

}

func (c *client) Post(path string, query url.Values, postBody interface{}, successResult interface{}, errorResult interface{}) (*http.Response, error) {

}

func (c *client) Patch(path string, query url.Values, patchBody interface{}, successResult interface{}, errorResult interface{}) (*http.Response, error) {

}

func (c *client) Head(path string, successResult interface{}, errorResultg interface{}) (*http.Response, error) {

}

func (c *client) Options(path string, successResult interface{}, errorResult interface{}) (*http.Response, error) {

}

func (c *client) Delete(path string, successResult interface{}, errorResult interface{}) (*http.Response, error) {

}

func (c *client) do(r *http.Request) (*http.Response, error) {
	var err error
	if c.RequestMutators() != nil {
		for _, m := range c.RequestMutators() {
			err = m(r)
			if err != nil {
				return nil, err
			}
		}
	}
	var response *http.Response
	client := c.GetHttpClient()

	response, err = client.Do(r)
	if err != nil {
		return nil, err
	}

	if c.ResponseMutators() != nil {
		var err error
		for _, m := range c.ResponseMutators() {
			err = m(response)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, nil
}

func (c *client) SetBaseUrl(u *url.URL) error {
	if u == nil {
		return errors.New("Please specify a non nil url.")
	}
	c.base = u
	return nil
}

func (c *client) CloneWithNewBaseUrl(base *url.URL) Client {
	cc := &client{}
	cc.base = base
	//TODO: complete clone method
	return cc
}

func (c *client) GetHttpClient() *http.Client {
	if c.client == nil {
		c.client = &http.Client{}
	}
	return c.client
}

func (c *client) SetHttpClient(h *http.Client) {
	c.client = h
}

func (c *client) SetMarshaler(f MarshalerFunc) {
	c.marshaler = f
}

func (c *client) SetUnmarshaler(f UnmarshalerFunc) {
	c.unmarshaler = f
}

func New(base *url.URL) (Client, error) {

	if base == nil {
		return nil, errors.New("Please specify a non nil url.")
	}
	c := &client{}
	c.headers = make(http.Header)
	c.query = make(url.Values)
	c.reqMutators = make([]RequestMutator, 0)
	c.resMutators = make([]ResponseMutator, 0)

	return c, nil
}
