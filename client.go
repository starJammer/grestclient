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

type Client interface {

	//Headers returns the default headers that will
	//be set with every request made with the client.
	//If none have been set, this should return an empty Header map.
	Headers() http.Header

	//SetHeaders sets the default headers that will be sent with
	//every request made with this client.
	//These headers can be overridden on a per-request-basis
	SetHeaders(http.Header)

	//Query returns the default query to use for all requests
	//An empty url.Values map should be returned if none has been set.
	Query() url.Values

	//SetQuery sets the default query to use for all requests
	SetQuery(url.Values)

	//SetBaseUrl sets the base url to use for all requests
	//If you want to use a different base url then you must
	//create a new client with the new base url. Call
	//CloneWithNewBaseUrl( url ) to get a clone of this
	//client's settings but with a new base url.
	//An error should be returned if the url is "unsupported",
	//whatever that may mean.
	//If you wish to use a username/password combination, set
	//the userinfo on the url. It will be used during requests.
	//To use one client with credentials and another without,
	//use CloneWithNewBaseUrl  and then change the base url.
	//Any query parameters added here should be ignored.
	//Clients should use the SetQuery method to set default
	//query parameters
	SetBaseUrl(*url.URL) error

	//BaseUrl returns the base url being used
	BaseUrl() *url.URL

	//Clones the client with everything the old client had
	//except for a new base url. The two clients should be
	//independent of each other and can be changed
	//without affecting each other.
	CloneWithNewBaseUrl(*url.URL) Client

	//GetHttpClient returns the current http.Client being used
	//If none has been set, this should return http.DefaultClient
	GetHttpClient() *http.Client

	//SetHttpClient sets the http.Client to use during requests
	//Use this to customize your http.Client as you wish. If you
	//don't set one, the default http.Client will be used.
	SetHttpClient(*http.Client)

	//SetMarshaler sets the marshal function to be used
	//to marshal the request bodies for requests
	//Doesn't have to mirror the Unmarshaler. Send plain text, get back json
	//Default is a json marshaler
	SetMarshaler(MarshalerFunc)
	//SetUnmarshaler sets the unmarshal function to be used
	//to unmarshal the response body for responses
	//Doesn't have to mirror the Marshaler. Send XML, get back json
	//Default is a json unmarshaler
	SetUnmarshaler(UnmarshalerFunc)

	//AddRequestMutator adds a mutator that the request will be
	//passed through before executing the request. All RequestMutators
	//are called AFTER the Marshaler is used.
	//RequestMutators should be called in the order they were added
	AddRequestMutators(...RequestMutator) Client

	//AddResponseMutator adds a mutator that the response will be
	//passed through after the server responds. All ResponseMutators
	//are called BEFORE the Unmarshaler is used.
	//ResponseMutators should be called in the order they were added
	AddResponseMutators(...ResponseMutator) Client

	//RemoveRequestMutator removes a request mutator
	SetRequestMutators(...RequestMutator) Client
	//RemoveResponseMutator removes a response mutator
	SetResponseMutators(...ResponseMutator) Client

	//Returns the RequestMutators
	RequestMutators() []RequestMutator
	//Returns the ResponseMutators
	ResponseMutators() []ResponseMutator

	//Get performs a get request with the base url plus the path appended to it. You can send query values and
	//supply a successResult that will be populated if the http response has a return code of 300.
	//errorResult is populated if the error code is 400 or more
	//Returns the raw http.Response and error similar to Do method of http.Client
	//The returned http.Response might be non-nil even though an error was also returned
	//depending on where the operation failed.
	Get(path string, query url.Values, successResult interface{}, errorResult interface{}) (*http.Response, error)

	//Post performs a post request with the base url plus the path appended to it. You can send query values and
	//supply a successResult that will be populated if the http response has a return code of 300.
	//errorResult is populated if the error code is 400 or more
	//With post you can also provide a post body.
	//Returns the raw http.Response and error similar to Do method of http.Client
	//The returned http.Response might be non-nil even though an error was also returned
	//depending on where the operation failed.
	Post(path string, query url.Values, postBody interface{}, successResult interface{}, errorResult interface{}) (*http.Response, error)

	//Put performs a put request with the base url plus the path appended to it. You can send query values and
	//supply a successResult that will be populated if the http response has a return code of 300.
	//errorResult is populated if the error code is 400 or more
	//With put you can also provide a put body.
	//Returns the raw http.Response and error similar to Do method of http.Client
	//The returned http.Response might be non-nil even though an error was also returned
	//depending on where the operation failed.
	Put(path string, query url.Values, putBody interface{}, successResult interface{}, errorResult interface{}) (*http.Response, error)

	//Patch performs a patch request with the base url plus the path appended to it. You can send query values and
	//supply a successResult that will be populated if the http response has a return code of 300.
	//errorResult is populated if the error code is 400 or more
	//With patch you can also provide a patch body.
	//Returns the raw http.Response and error similar to Do method of http.Client
	//The returned http.Response might be non-nil even though an error was also returned
	//depending on where the operation failed.
	Patch(path string, query url.Values, patchBody interface{}, successResult interface{}, errorResult interface{}) (*http.Response, error)

	//Head performs a head request with the base url plus the path appended to it.
	//supply a successResult that will be populated if the http response has a return code of 300.
	//errorResult is populated if the error code is 400 or more
	//Returns the raw http.Response and error similar to Do method of http.Client
	//The returned http.Response might be non-nil even though an error was also returned
	//depending on where the operation failed.
	Head(path string, successResult interface{}, errorResultg interface{}) (*http.Response, error)

	//Option performs an option request with the base url plus the path appended to it.
	//supply a successResult that will be populated if the http response has a return code of 300.
	//errorResult is populated if the error code is 400 or more
	//Returns the raw http.Response and error similar to Do method of http.Client
	//The returned http.Response might be non-nil even though an error was also returned
	//depending on where the operation failed.
	Options(path string, successResult interface{}, errorResult interface{}) (*http.Response, error)

	//Delete performs an delete request with the base url plus the path appended to it.
	//supply a successResult that will be populated if the http response has a return code of 300.
	//errorResult is populated if the error code is 400 or more
	//Returns the raw http.Response and error similar to Do method of http.Client
	//The returned http.Response might be non-nil even though an error was also returned
	//depending on where the operation failed.
	Delete(path string, query url.Values, successResult interface{}, errorResult interface{}) (*http.Response, error)
}

//MarshalerFunc takes something and converts it into a
//io.ReadCloser that can be used for the request body
type MarshalerFunc func(v interface{}) (io.ReadCloser, error)

//UnmarshalerFunc takes the response body and converts it into
//something you can use.
type UnmarshalerFunc func(body io.ReadCloser, v interface{}) error

func ByteSliceToReadCloser(b []byte) (io.ReadCloser, error) {
	if b == nil {
		return nil, errors.New("ReadCloserFromByteSlice received a nil byte slice.")
	}

	buf := bytes.NewBuffer(b)
	return ioutil.NopCloser(buf), nil
}

func StringToReadCloser(s string) (io.ReadCloser, error) {
	buf := bytes.NewBufferString(s)
	return ioutil.NopCloser(buf), nil
}

func JsonMarshalerFunc(v interface{}) (io.ReadCloser, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	body, err := ByteSliceToReadCloser(b)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func JsonUnmarshalerFunc(body io.ReadCloser, v interface{}) error {
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, v)

	if err != nil {
		return err
	}

	return nil
}

func StringMarshalerFunc(v interface{}) (io.ReadCloser, error) {
	switch t := v.(type) {
	case fmt.Stringer:
		return StringToReadCloser(t.String())
	case string:
		return StringToReadCloser(t)
	}

	return nil, errors.New("Did not know how to use the body as text.")
}

func StringUnmarshalerFunc(body io.ReadCloser, v interface{}) error {
	defer body.Close()

	if v == nil {
		return nil
	}

	b, err := ioutil.ReadAll(body)

	if err != nil {
		return err
	}

	switch v.(type) {
	case *string:
		elem := rt.ValueOf(v).Elem()

		if elem.CanSet() {
			elem.SetString(string(b))
			return nil
		}
	case string:
		return errors.New("You must pass the string by reference: ")
	}

	return errors.New("Did not know how to unmarshal the text coming back.")
}

//RequestMutators are called before the request is made but after the marshaler function has been
//called.
type RequestMutator func(*http.Request) error
type ResponseMutator func(*http.Response) error
