package grestclient

import (
	"net/http"
	"net/url"
)

type Client interface {

	//Headers returns the default headers that will
	//be set with every request made with the client.
	//If none have been set, this should return an empty Header map
	//These headers can be overridden on a per-request-basis
	Headers() http.Header

	//SetHeaders sets the default headers that will be sent with
	//every request made with this client.
	//These headers can be overridden on a per-request-basis
	SetHeaders(http.Header)

	//Query returns the default query to use for all requests
	//An empty url.Values map should be returned if none has been set.
	Query() url.Values
	//SetQuery sets the default queries to use for all requests
	SetQuery(url.Values)

	//SetBaseUrl sets the base url to use for all requests
	//If you want to use a different base url then you must
	//create a new client with the new base url. Can call
	//CloneWithNewBaseUrl( url ) to get a clone of this
	//client's settings but with a new base url.
	//An error should be returned if the url is unsupported,
	//for example, using an inadequate scheme. Each implementation
	//can support different url types.
	SetBaseUrl(*url.URL) error

	CloneWithNewBaseUrl(*url.URL) Client

	//GetHttpClient returns the current http.Client being used
	//If none has been set, this should return http.DefaultClient
	GetHttpClient() *http.Client

	//SetHttpClient sets the http.Client to use during requests
	//Use this to customize your http.Client as you wish. If you
	//don't set one, the default http.Client will be used.
	SetHttpClient(*http.Client)

	//SetMarshaler sets the marshal function to be used
	//to marshal
	SetMarshaler(func(interface{}) ([]byte, error))
	//SetUnmarshaler sets the unmarshal function to be used
	//to unmarshal the body of the REST response
	SetUnmarshaler(func([]byte, interface{}) error)

	//AddRequestMutator adds a mutator that the request will be
	//passed through before executing the request. All RequestMutators
	//are called AFTER the Marshaler is used.
	//RequestMutators should be called in the order they were added
	AddRequestMutator(RequestMutator) Client

	//AddResponseMutator adds a mutator that the response will be
	//passed through after the server responds. All ResponseMutators
	//are called BEFORE the Unmarshaler is used.
	//ResponseMutators should be called in the order they were added
	AddResponseMutator(RequestMutator) Client

	//RemoveRequestMutator removes a request mutator
	RemoveRequestMutator(RequestMutator) Client
	//RemoveResponseMutator removes a response mutator
	RemoveResponseMutator(RequestMutator) Client

	//Returns the RequestMutators
	RequestMutator() []RequestMutator
	//Returns the ResponseMutators
	ResponseMutator() []ResponseMutator

	Get(path string, query url.Values, successResult interface{}, errorResult interface{}) (*http.Response, error)
	Post(path string, query url.Values, postBody interface{}, successResult interface{}, errorResult interface{}) (*http.Response, error)
	Patch(path string, query url.Values, patchBody interface{}, successResult interface{}, errorResult interface{}) (*http.Response, error)

	Head(path string, successResult interface{}, errorResultg interface{}) (*http.Response, error)
	Options(path string, successResult interface{}, errorResult interface{}) (*http.Response, error)
	Delete(path string, successResult interface{}, errorResult interface{}) (*http.Response, error)
}

type RequestMutator func(*http.Request) error
type ResponseMutator func(*http.Response) error

//clone will clone the url you pass in
//with a deep copy
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
