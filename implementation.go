package grestclient

import (
	"errors"
	"io/ioutil"
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

func (c *client) SetRequestMutators(rm ...RequestMutator) Client {
	c.reqMutators = rm
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
	r, err := c.prepareRequest("GET", path, query, nil)
	if err != nil {
		return nil, err
	}
	return c.do(r, successResult, errorResult)
}

func (c *client) Post(path string, query url.Values, postBody interface{}, successResult interface{}, errorResult interface{}) (*http.Response, error) {
	r, err := c.prepareRequest("POST", path, query, postBody)
	if err != nil {
		return nil, err
	}
	return c.do(r, successResult, errorResult)
}

func (c *client) Put(path string, query url.Values, putBody interface{}, successResult interface{}, errorResult interface{}) (*http.Response, error) {
	r, err := c.prepareRequest("PUT", path, query, putBody)
	if err != nil {
		return nil, err
	}
	return c.do(r, successResult, errorResult)
}

func (c *client) Patch(path string, query url.Values, patchBody interface{}, successResult interface{}, errorResult interface{}) (*http.Response, error) {
	r, err := c.prepareRequest("PATCH", path, query, patchBody)
	if err != nil {
		return nil, err
	}
	return c.do(r, successResult, errorResult)
}

func (c *client) Head(path string, successResult interface{}, errorResult interface{}) (*http.Response, error) {
	r, err := c.prepareRequest("HEAD", path, nil, nil)
	if err != nil {
		return nil, err
	}
	return c.do(r, successResult, errorResult)
}

func (c *client) Options(path string, successResult interface{}, errorResult interface{}) (*http.Response, error) {
	r, err := c.prepareRequest("OPTIONS", path, nil, nil)
	if err != nil {
		return nil, err
	}
	return c.do(r, successResult, errorResult)
}

func (c *client) Delete(path string, query url.Values, successResult interface{}, errorResult interface{}) (*http.Response, error) {
	r, err := c.prepareRequest("DELETE", path, query, nil)
	if err != nil {
		return nil, err
	}
	return c.do(r, successResult, errorResult)
}

func (c *client) do(r *http.Request, successResult interface{}, errorResult interface{}) (*http.Response, error) {
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
	defer response.Body.Close()

	if c.ResponseMutators() != nil {
		var err error
		for _, m := range c.ResponseMutators() {
			err = m(response)
			if err != nil {
				return response, err
			}
		}
	}

	if c.unmarshaler == nil {
		c.unmarshaler = StringUnmarshalerFunc
	}

	//make sure there is a body, or that there might be a body (when it is -1)
	if response.ContentLength > 0 || response.ContentLength == -1 {

		//unmarshal it depending on StatusCode
		if response.StatusCode < 300 && successResult != nil {
			//success
			err = c.unmarshaler(response.Body, successResult)
		} else if response.StatusCode < 600 && errorResult != nil {
			//error
			err = c.unmarshaler(response.Body, errorResult)
		}
	}

	if err != nil {
		//we have the response so return it even though unmarshaling might've
		//produced an error
		return response, err
	}

	return response, nil
}

func (c *client) prepareRequest(
	method string,
	path string,
	query url.Values,
	body interface{}) (*http.Request, error) {

	var err error
	reqUrl := cloneUrl(c.base)
	reqUrl.Path += path

	//set headers
	headers := setupHeaders(c.headers)
	//create query
	query = setupQuery(c.query, query)

	if c.marshaler == nil {
		c.marshaler = StringMarshalerFunc
	}

	r, err := http.NewRequest(method, reqUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	r.Header = headers
	r.URL.RawQuery = query.Encode()

	var readLener ReadLener
	if body != nil {

		readLener, err = c.marshaler(body)

		if err != nil {
			return nil, err
		}
		r.ContentLength = int64(readLener.Len())
		r.Body = ioutil.NopCloser(readLener)
	}

	return r, nil
}

func setupHeaders(headers ...http.Header) http.Header {
	finalheaders := make(http.Header)
	for _, current := range headers {
		for i, v := range current {
			finalheaders[i] = v
		}
	}

	return finalheaders
}

func setupQuery(queries ...url.Values) url.Values {
	finalquery := make(url.Values)
	for _, current := range queries {
		for i, v := range current {
			finalquery[i] = v
		}
	}

	return finalquery
}

func (c *client) SetBaseUrl(u *url.URL) error {
	if u == nil {
		return errors.New("Please specify a non nil url.")
	}
	u.RawQuery = ""
	c.base = u
	return nil
}

func (c *client) BaseUrl() *url.URL {
	return c.base
}

func (c *client) Clone() Client {
	cc := &client{}
	cc.base = cloneUrl(c.base)

	cc.reqMutators = make([]RequestMutator, len(c.reqMutators))
	copy(cc.reqMutators, c.reqMutators)

	cc.resMutators = make([]ResponseMutator, len(c.resMutators))
	copy(cc.resMutators, c.resMutators)

	cc.headers = headerCopy(c.headers)
	cc.query = queryCopy(c.query)
	cc.client = c.client
	cc.marshaler = c.marshaler
	cc.unmarshaler = c.unmarshaler

	return cc
}

func headerCopy(h http.Header) http.Header {
	if h == nil {
		return nil
	}
	c := make(http.Header, len(h))
	for i, v := range h {
		c[i] = v
	}
	return c
}

func queryCopy(q url.Values) url.Values {
	if q == nil {
		return nil
	}
	c := make(url.Values, len(q))
	for i, v := range q {
		c[i] = v
	}
	return c
}

func (c *client) GetHttpClient() *http.Client {
	if c.client == nil {
		c.client = http.DefaultClient
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

//New creates a new grestclient with the base url set
//to the passed in paramater.
func New(base *url.URL) (Client, error) {

	if base == nil {
		return nil, errors.New("Please specify a non nil url.")
	}
	c := &client{}
	c.base = base

	return c, nil
}
