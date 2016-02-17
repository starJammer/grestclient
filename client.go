package grestclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	rt "reflect"
)

//Client lets you maintain query and header across http requests
//It will let you set (un)marshalers to simplify dealing with
//json and string responses.
type Client struct {
	base        *url.URL
	reqMutators []RequestMutator
	resMutators []ResponseMutator
	headers     http.Header
	query       url.Values
	client      *http.Client
	marshaler   MarshalerFunc
	unmarshaler UnmarshalerFunc
}

//Request represents a single request that can be made
//with the Client's Do method.
type Request struct {
	Path    string
	Headers http.Header
	Query   url.Values
	//Body is ignored by GET, HEAD and other methods that don't typically
	//have a body
	Body         interface{}
	UnmarshalMap UnmarshalMap
}

//Headers returns the default headers that will
//be set with every request made with the client.
func (c *Client) Headers() http.Header {
	if c.headers == nil {
		c.headers = make(http.Header)
	}
	return c.headers
}

//SetHeaders sets the default headers that will be sent with
//every request made with this client.
func (c *Client) SetHeaders(h http.Header) {
	c.headers = h
}

//Query returns the default query to use for all requests
func (c *Client) Query() url.Values {
	if c.query == nil {
		c.query = make(url.Values)
	}
	return c.query
}

//SetQuery sets the default query to use for all requests
func (c *Client) SetQuery(q url.Values) {
	c.query = q
}

//SetBaseUrl sets the base url to use for all requests
//If you want to use a different temporarily it is best to
//create a new client with the new base url. Call
//Clone() to get a clone of this
//client's settings and then change the url on the clone.
//An error should be returned if the url is "unsupported" by
//the implementation. For example, "unix://tmp.soc".
//Any query parameters added here should be ignored.
//Clients should use the SetQuery method to set default
//query parameters
func (c *Client) SetBaseUrl(u *url.URL) error {
	if u == nil {
		return errors.New("Please specify a non nil url.")
	}
	u.RawQuery = ""
	c.base = u
	return nil
}

//BaseUrl returns the base url being used. This implementation
//allows you to change the base url here directly but other
//implementations might give you a clone so changing it won't affect
//the client. In those cases, use SetBaseUrl to change the url after
//making your changes.
func (c *Client) BaseUrl() *url.URL {
	return c.base
}

//Clones the client with everything the old client had
//Ideally, clones should be independent of the original and can be changed
//without affecting the original and vice versa.
//This implementation, however, shared the http.Client among clones.
//All other 'things' like headers, base url, query, marshalers are separate and
//can be adjusted without affecting the original/clones.
func (c *Client) Clone() *Client {
	cc := &Client{}
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

//AddRequestMutator adds a mutator that the request will be
//passed through before executing the request. All RequestMutators
//are called AFTER the Marshaler is used.
//RequestMutators should be called in the order they were added
func (c *Client) AddRequestMutators(rm ...RequestMutator) *Client {
	c.reqMutators = append(c.reqMutators, rm...)
	return c
}

//AddResponseMutator adds a mutator that the response will be
//passed through after the server responds. All ResponseMutators
//are called BEFORE the Unmarshaler is used.
//ResponseMutators should be called in the order they were added
func (c *Client) AddResponseMutators(rm ...ResponseMutator) *Client {
	c.resMutators = append(c.resMutators, rm...)
	return c
}

//SetRequestMutators removes a request mutator
func (c *Client) SetRequestMutators(rm ...RequestMutator) *Client {
	c.reqMutators = rm
	return c
}

//SetResponseMutators removes a response mutator
func (c *Client) SetResponseMutators(rm ...ResponseMutator) *Client {
	c.resMutators = rm
	return c
}

//Returns the RequestMutators
func (c *Client) RequestMutators() []RequestMutator {
	return c.reqMutators
}

//Returns the ResponseMutators
func (c *Client) ResponseMutators() []ResponseMutator {
	return c.resMutators
}

//Get performs a get request with the base url plus the path appended to it.
//You can send query values, header values and
//supply a successResult that will be populated if the http response has a return code less than 400.
//errorResult is populated if the error code is 400 or more
//Returns the raw http.Response and error similar to Do method of http.Client
//The returned http.Response might be non-nil even though an error was also returned
//depending on where the operation failed.
func (c *Client) Get(req *Request) (*http.Response, error) {
	r, err := c.prepareRequest("GET", req.Path, req.Headers, req.Query, nil)
	if err != nil {
		return nil, err
	}
	return c.do(r, req.UnmarshalMap)
}

//Post performs a post request with the base url plus the path appended to it.
//You can send query values, header values and
//supply a successResult that will be populated if the http response has a return code less than 400.
//errorResult is populated if the error code is 400 or more
//With post you can also provide a post body.
//Returns the raw http.Response and error similar to Do method of http.Client
//The returned http.Response might be non-nil even though an error was also returned
//depending on where the operation failed.
func (c *Client) Post(req *Request) (*http.Response, error) {
	r, err := c.prepareRequest("POST", req.Path, req.Headers, req.Query, req.Body)
	if err != nil {
		return nil, err
	}
	return c.do(r, req.UnmarshalMap)
}

//Put performs a put request with the base url plus the path appended to it.
//You can send query values, header values and
//supply a successResult that will be populated if the http response has a return code less than 400.
//errorResult is populated if the error code is 400 or more
//With put you can also provide a put body.
//Returns the raw http.Response and error similar to Do method of http.Client
//The returned http.Response might be non-nil even though an error was also returned
//depending on where the operation failed.
func (c *Client) Put(path string, headers http.Header, query url.Values, putBody interface{}, unmarshalMap UnmarshalMap) (*http.Response, error) {
	r, err := c.prepareRequest("PUT", path, headers, query, putBody)
	if err != nil {
		return nil, err
	}
	return c.do(r, unmarshalMap)
}

//Patch performs a patch request with the base url plus the path appended to it.
//You can send query values, header values and
//supply a successResult that will be populated if the http response has a return code less than 400.
//errorResult is populated if the error code is 400 or more
//With patch you can also provide a patch body.
//Returns the raw http.Response and error similar to Do method of http.Client
//The returned http.Response might be non-nil even though an error was also returned
//depending on where the operation failed.
func (c *Client) Patch(path string, headers http.Header, query url.Values, patchBody interface{}, unmarshalMap UnmarshalMap) (*http.Response, error) {
	r, err := c.prepareRequest("PATCH", path, headers, query, patchBody)
	if err != nil {
		return nil, err
	}
	return c.do(r, unmarshalMap)
}

//Head performs a head request with the base url plus the path appended to it.
//You can send header values and supply a successResult that will be populated
//if the http response has a return code less than 400.
//errorResult is populated if the error code is 400 or more
//Returns the raw http.Response and error similar to Do method of http.Client
//The returned http.Response might be non-nil even though an error was also returned
//depending on where the operation failed.
func (c *Client) Head(path string, headers http.Header, query url.Values) (*http.Response, error) {
	r, err := c.prepareRequest("HEAD", path, headers, query, nil)
	if err != nil {
		return nil, err
	}
	return c.do(r, nil)
}

//Option performs an option request with the base url plus the path appended to it.
//You can send header values and supply a successResult that will be populated
//if the http response has a return code less than 400.
//errorResult is populated if the error code is 400 or more
//Returns the raw http.Response and error similar to Do method of http.Client
//The returned http.Response might be non-nil even though an error was also returned
//depending on where the operation failed.
func (c *Client) Options(path string, headers http.Header, query url.Values, optionsBody interface{}, unmarshalMap UnmarshalMap) (*http.Response, error) {
	r, err := c.prepareRequest("OPTIONS", path, headers, query, optionsBody)
	if err != nil {
		return nil, err
	}
	return c.do(r, unmarshalMap)
}

//Delete performs an delete request with the base url plus the path appended to it.
//You can send header values and supply a successResult that will be populated
//if the http response has a return code less than 400.
//errorResult is populated if the error code is 400 or more
//Returns the raw http.Response and error similar to Do method of http.Client
//The returned http.Response might be non-nil even though an error was also returned
//depending on where the operation failed.
func (c *Client) Delete(path string, headers http.Header, query url.Values, unmarshalMap UnmarshalMap) (*http.Response, error) {
	r, err := c.prepareRequest("DELETE", path, headers, query, nil)
	if err != nil {
		return nil, err
	}
	return c.do(r, unmarshalMap)
}

//UnmarshalMap represents a mapping from HTTP status
//codes to interfaces that a client should unmarshal
//to.
//For example, you call
//Get( path, headers, query, UnmarshalMap{ 200 : success, 201 : someothersuccess, 202 : success, 404 : uauthorizedPiece }
//If the http response is either a 202 or a 200 then the response body is unmarshaled into success.
//A 201 response unmarshals into someothersuccess, and a 404 unmarshals into unauthorizedPiece
type UnmarshalMap map[int]interface{}

//ReadLener is an io.Reader than can tell you the length of its content.
//Len can either be 0 for no bytes, -1 for an unknown number of bytes,
//or >0 for a specific number of bytes.
//
//I made ReadLener because ContentLength needs to be set
//to something in the http.Request. Your Len method can return 0 or -1 if you
//want but some APIs depend on the ContentLength being set correctly and accurately.
//bytes.Buffers is already a ReadLener, it has Read and Len methods
//so no worries there.
//I ran into issues with the ArangoDB REST API where it failed
//to answer requests properly if ContentLength was -1 or if it was
//not correct/accurate.
type ReadLener interface {
	io.Reader
	Len() int
}

//MarshalerFunc takes v, marshals it, and converts it into a
//ReadLener that can be used for the htt.Request.Body.
//bytes.Buffer is a ReadLener.
//Use StringToReadLener to convert strings to a ReadLener
//Use ByteSliceToReadLener to convert a []byte into a ReadLener
type MarshalerFunc func(v interface{}) (ReadLener, error)

//UnmarshalerFunc is passed a []byte from the response body
//and should unmarshal it into v.
type UnmarshalerFunc func(b []byte, v interface{}) error

//ByteSliceToReadCloser takes a byte slice and converts it to an
//ReadLener that can be used as a request/resonse body
func ByteSliceToReadLener(b []byte) (ReadLener, error) {
	if b == nil {
		return nil, errors.New("ByteSliceToReadLener received a nil byte slice.")
	}

	buf := bytes.NewBuffer(b)
	return buf, nil
}

//StringToReadLener takes a string and converts it to an
//ReadLener that can be used as a request/resonse body
func StringToReadLener(s string) ReadLener {
	buf := bytes.NewBufferString(s)
	return buf
}

//JsonMarshalerFunc can be used by the client to marshal
//request bodies into json.
func JsonMarshalerFunc(v interface{}) (ReadLener, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	body, err := ByteSliceToReadLener(b)
	if err != nil {
		return nil, err
	}
	return body, nil
}

//JsonUnmarshalerFunc can be used to unmarshal response bodies
//from json.
func JsonUnmarshalerFunc(b []byte, v interface{}) error {

	err := json.Unmarshal(b, v)

	if err != nil {
		return err
	}

	return nil
}

//StringMarshalerFunc can be used to marshal strings into a request.
func StringMarshalerFunc(v interface{}) (ReadLener, error) {
	switch t := v.(type) {
	case fmt.Stringer:
		return StringToReadLener(t.String()), nil
	case string:
		return StringToReadLener(t), nil
	}

	return nil, errors.New("Did not know how to use the body as text.")
}

//StringUnmarshalerFunc can be used to unmarshal strings from a response.
func StringUnmarshalerFunc(b []byte, v interface{}) error {

	switch v.(type) {
	case *string:
		elem := rt.ValueOf(v).Elem()

		if elem.CanSet() {
			elem.SetString(string(b))
			return nil
		}
	case string:
		return errors.New("You must pass the string by reference or pass a pointer. For example, ( &stringVar )")
	}

	return errors.New("Did not know how to unmarshal the text coming back.")
}

//JsonContentTypeMutator sets the Content-Type of the request to be
//application/json
func JsonContentTypeMutator(r *http.Request) error {
	r.Header.Add("Content-Type", "application/json")
	return nil
}

func JsonAcceptMutator(r *http.Request) error {
	r.Header.Add("Accept", "application/json")
	return nil
}

//RequestMutators are called before the request is made but after the marshaler function has been
//called.
type RequestMutator func(*http.Request) error
type ResponseMutator func(*http.Response) error

//SetupClientForJson is a convenience method that sets the
//marshaler and unmarshaler funcs on the client to be the
//Json funcs in this package. It also sets a request mutator
//to set the Content-Type to json.
func SetupForJson(c *Client) {
	c.SetMarshaler(JsonMarshalerFunc)
	c.SetUnmarshaler(JsonUnmarshalerFunc)
	c.AddRequestMutators(JsonContentTypeMutator)
	c.AddRequestMutators(JsonAcceptMutator)
}

func (c *Client) do(r *http.Request, unmarshalMap UnmarshalMap) (*http.Response, error) {
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
	client := c.GetHttpDoer()

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

	if unmarshalMap != nil {
		//make sure there is a body, or that there might be a body (when it is -1)
		if response.ContentLength > 0 || response.ContentLength == -1 {
			//unmarshal it depending on StatusCode
			if destination, ok := unmarshalMap[response.StatusCode]; ok && destination != nil {
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					return response, err
				}
				err = c.unmarshaler(body, destination)
			}
		}
	}

	if err != nil {
		//we have the http response so return it even though unmarshaling might've
		//produced an error
		return response, err
	}

	return response, nil
}

func (c *Client) prepareRequest(
	method string,
	path string,
	headers http.Header,
	query url.Values,
	body interface{}) (*http.Request, error) {

	var err error
	reqUrl := cloneUrl(c.base)
	reqUrl.Path += path

	//set headers
	headers = setupHeaders(c.headers, headers)
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

//GetHttpClient returns the current http.Client being used
//If none has been set, this should return http.DefaultClient
func (c *Client) GetHttpDoer() HttpDoer {
	if c.client == nil {
		c.client = http.DefaultClient
	}
	return c.client
}

type HttpDoer interface {
	Do(r *http.Request) (*http.Response, error)
}

//SetHttpClient sets the http.Client to use during requests
//Use this to customize your http.Client as you wish. If you
//don't set one, the default http.Client will be used.
func (c *Client) SetHttpDoer(h *http.Client) {
	c.client = h
}

//SetMarshaler sets the marshal function to be used
//to marshal the request bodies for requests
//Doesn't have to mirror the Unmarshaler. Send plain text, get back json
//Default is a json marshaler
func (c *Client) SetMarshaler(f MarshalerFunc) {
	c.marshaler = f
}

//SetUnmarshaler sets the unmarshal function to be used
//to unmarshal the response body for responses
//Doesn't have to mirror the Marshaler. Send XML, get back json
//Default is a json unmarshaler
func (c *Client) SetUnmarshaler(f UnmarshalerFunc) {
	c.unmarshaler = f
}

//New creates a new grestclient with the base url set
//to the passed in paramater.
func New(base *url.URL) (*Client, error) {

	if base == nil {
		return nil, errors.New("Please specify a non nil url.")
	}
	c := &Client{}
	c.base = base

	return c, nil
}

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
