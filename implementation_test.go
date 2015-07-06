package grestclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestMeetsInterface(t *testing.T) {
	base, err := url.Parse("http://example.com")
	c, err := New(base)

	if err != nil {
		t.Fatal(err)
	}

	if c.BaseUrl().String() != base.String() {
		t.Fatal("Base url for client and base passed in don't match.")
	}
}

func TestStringToReadCloser(t *testing.T) {
	reader, err := StringToReadCloser("thing")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(reader)

	if string(b) != "thing" {
		t.Fatal("ReadCloser not created properly from string.")
	}
}

func TestTextMarshalers(t *testing.T) {
	r, err := StringMarshalerFunc("cosa")
	if err != nil {
		t.Fatal(err)
	}
	var result string
	err = StringUnmarshalerFunc(r, &result)

	if err != nil {
		t.Fatal(err)
	}

	if result != "cosa" {
		t.Fatal("Could not properly unmarshal text: ", result)
	}

	//passing by value doesn't work
	err = StringUnmarshalerFunc(r, result)
	if err == nil {
		t.Fatal("Expected error about passing by value instead of reference.")
	}
}

func TestGetMethod(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			t.Fatal("Expected GET but got: ", req.Method)
		}
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Get("get", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Didn't get a response back.")
	}
}

func TestPostMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Fatal("Expected POST but got: ", req.Method)
		}
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Post("post", nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Didn't get a response back.")
	}
}

func TestPutMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			t.Fatal("Expected PUT but got: ", req.Method)
		}
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Put("put", nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Didn't get a response back.")
	}
}

func TestPatchMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "PATCH" {
			t.Fatal("Expected PATCH but got: ", req.Method)
		}
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Patch("patch", nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Didn't get a response back.")
	}
}

func TestHeadMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "HEAD" {
			t.Fatal("Expected HEAD but got: ", req.Method)
		}
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Head("head", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Didn't get a response back.")
	}
}

func TestOptionsMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "OPTIONS" {
			t.Fatal("Expected OPTIONS but got: ", req.Method)
		}
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Options("options", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Didn't get a response back.")
	}
}

func TestDeleteMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			t.Fatal("Expected DELETE but got: ", req.Method)
		}
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Delete("delete", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Didn't get a response back.")
	}
}

func TestDefaultHeaderQueryPassedIntoGetRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			t.Fatal("Expected GET method but got : ", req.Method)
		}

		if req.URL.Path != "/get" {
			t.Fatal("Did not get correct path: ", req.URL.Path)
		}

		if v := req.Header.Get("Test-Header"); v != "test" {
			t.Fatal("Did not receive expected test header.")
		}
		if v := req.Header["Multi-Value-Header"]; len(v) != 2 {
			t.Fatal("Multi-Value-Header was not 2 values.")
		}
		query := req.URL.Query()

		if query.Get("testquery") != "test" {
			t.Fatal("Did not receive expected test query: ", query.Encode())
		}

		if len(query["multiquery"]) != 2 {
			t.Fatal("multiquery was not 2 values.", query.Encode())
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}

	client.Headers().Add("Test-Header", "test")
	client.Headers()["Multi-Value-Header"] = []string{"test", "test2"}

	client.Query().Add("testquery", "test")
	client.Query()["multiquery"] = []string{"test", "test2"}

	res, err := client.Get("get", nil, nil, nil)

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Error("Expected a 200 code from the server but got :", res.StatusCode)
	}

}

func TestStringMarshaledBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Fatal("Got an error when reading text body: ", err)
		}

		if string(b) != "hello" {
			t.Fatal("Incorrect body sent: ", string(b))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("world"))
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}
	var success string
	res, err := client.Post("post", nil, "hello", &success, nil)

	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Didn't get a response back.")
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("Unexpected status code: ", res.StatusCode)
	}
	if success != "world" {
		t.Fatal("Unmarshaling didn't work: ", success)
	}
}

func TestErrorResultUnmarshaledOnError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("world"))
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}
	var success string
	var errResult string
	res, err := client.Post("post", nil, "hello", &success, &errResult)

	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Didn't get a response back.")
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatal("Unexpected status code: ", res.StatusCode)
	}
	if success != "" {
		t.Fatal("Success should not have a value on error.")
	}
	if errResult != "world" {
		t.Fatal("Error not populated on error.")
	}
}

func TestDumbRequestResponseMutators(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-Test") != "test" {
			t.Fatal("Request mutator failed.")
		}
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}

	client.AddRequestMutators(func(r *http.Request) error {
		r.Header.Add("X-Test", "test")
		return nil
	})

	client.AddResponseMutators(func(r *http.Response) error {
		r.Header.Add("X-Test-Response", "test")
		return nil
	})

	res, err := client.Post("post", nil, nil, nil, nil)

	if err != nil {
		t.Fatal(err)
	}
	if res.Header.Get("X-Test-Response") != "test" {
		t.Fatal("Response mutator failed.")
	}
}

func TestJsonMarshaledBody(t *testing.T) {
	type test struct {
		Name string `json:name`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if req.Header.Get("Content-Type") != "application/json" {
			t.Fatal("Expected that we'd get a Content-Type of application/json but got: ", req.Header.Get("Content-Type"))
		}

		b, _ := ioutil.ReadAll(req.Body)
		var tester test
		err := json.Unmarshal(b, &tester)

		if err != nil {
			t.Fatal(err)
		}

		if tester.Name != "test" {
			t.Fatal("Didn't get expected value.")
		}
		tester.Name = "test-result"
		b, _ = json.Marshal(tester)

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}

	SetupForJson(client)

	var body, success, errResult test
	body.Name = "test"

	_, err = client.Post("post", nil, body, &success, errResult)

	if err != nil {
		t.Fatal(err)
	}
	if success.Name != "test-result" {
		t.Fatal("Did not receive expected result.")
	}
}