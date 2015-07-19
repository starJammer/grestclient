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
	Headers() http.Header

	//SetHeaders sets the default headers that will be sent with
	//every request made with this client.
	SetHeaders(http.Header)

	//Query returns the default query to use for all requests
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

	//BaseUrl returns the base url being used. This implementation
	//allows you to change the base url here directly but other
	//implementations might give you a clone so changing it won't affect
	//the client. In those cases, use SetBaseUrl to change the url.
	BaseUrl() *url.URL

	//Clones the client with everything the old client had
	//Ideally, clones should be independent of the original and can be changed
	//without affecting the original and vice versa.
	//This implementation, however, shared the http.Client among clones.
	//All other 'things' like headers, base url, query, marshalers are separate and
	//can be adjusted without affecting the original/clones.
	Clone() Client

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

type ReadLener interface {
	io.Reader
	Len() int
}

//MarshalerFunc takes something and converts it into a
//ReadLener that can be used for the request body
type MarshalerFunc func(v interface{}) (ReadLener, error)

//UnmarshalerFunc takes the response body and converts it into
//something you can use.
type UnmarshalerFunc func(b io.ReadCloser, v interface{}) error

//ByteSliceToReadCloser takes a byte slice and converts it to an
//ReadLener that can be used as a request/resonse body
func ByteSliceToReadLener(b []byte) (ReadLener, error) {
	if b == nil {
		return nil, errors.New("ReadCloserFromByteSlice received a nil byte slice.")
	}

	buf := bytes.NewBuffer(b)
	return buf, nil
}

//StringToReadCloser takes a string and converts it to an
//ReadLener that can be used as a request/resonse body
func StringToReadLener(s string) ReadLener {
	buf := bytes.NewBufferString(s)
	return buf
}

//JsonMarshalerFunc can be used to marshal request bodies
//into json
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
//from json
func JsonUnmarshalerFunc(body io.ReadCloser, v interface{}) error {

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
func StringUnmarshalerFunc(body io.ReadCloser, v interface{}) error {
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
func SetupForJson(c Client) {
	c.SetMarshaler(JsonMarshalerFunc)
	c.SetUnmarshaler(JsonUnmarshalerFunc)
	c.AddRequestMutators(JsonContentTypeMutator)
	c.AddRequestMutators(JsonAcceptMutator)
}
