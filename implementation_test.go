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
	reader := StringToReadLener("thing")

	b, err := ioutil.ReadAll(reader)

	if err != nil {
		t.Fatal(err)
	}

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
	b, _ := ioutil.ReadAll(r)
	err = StringUnmarshalerFunc(b, &result)

	if err != nil {
		t.Fatal(err)
	}

	if result != "cosa" {
		t.Fatal("Could not properly unmarshal text: ", result)
	}

	//passing by value doesn't work
	err = StringUnmarshalerFunc(b, result)
	if err == nil {
		t.Fatal("Expected error about passing by value instead of reference.")
	}
}

func TestNoServer(t *testing.T) {

	base, err := url.Parse("http://127.0.0.1:9999")
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Get(&Params{Path: "get"})
	if err == nil {
		t.Fatal("Expected some error for there being no server.")
	}
	if res != nil {
		t.Fatal("Expected response to be nil.")
	}
}

func TestGetMethod(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			t.Fatal("Expected GET but got: ", req.Method)
		}

		if req.URL.Path != "/get" {
			t.Fatal("Expected path to be get but got: ", req.URL.Path)
		}
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Get(&Params{Path: "get"})
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

	res, err := client.Post(&Params{Path: "post"})
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

	res, err := client.Put(&Params{Path: "put"})
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

	res, err := client.Patch(&Params{Path: "patch"})
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

	res, err := client.Head(&Params{Path: "head"})
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

	res, err := client.Options(&Params{Path: "options"})
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

	res, err := client.Delete(&Params{Path: "delete"})
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Didn't get a response back.")
	}
}

func TestDefaultHeaderQueryPassedIntoGetRequest(t *testing.T) {
	var requestCount = 0
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

		requestCount++

		if requestCount == 1 {
			if oneOff, ok := req.Header["One-Off"]; !ok || oneOff[0] != "one" {
				t.Fatal("Expected a one off header in the first request.")
			}
		}

		if requestCount == 2 {
			if _, ok := req.Header["One-Off"]; ok {
				t.Fatal("Did NOT expect the one off header in the second request.")
			}
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

	res, err := client.Get(&Params{
		Path:    "get",
		Headers: http.Header{"One-Off": []string{"one"}},
	})

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Error("Expected a 200 code from the server but got :", res.StatusCode)
	}

	res, err = client.Get(&Params{
		Path: "get",
	})

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Error("Expected a 200 code from the server but got :", res.StatusCode)
	}

}

func TestQueryInCallOverridesDefaults(t *testing.T) {
	firstRequest := true
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if firstRequest {
			if req.URL.Query().Get("testquery") != "test-override" {
				t.Fatal("Did not get override query: ", req.URL.Query().Get("testquery"))
			}
		} else {
			if req.URL.Query().Get("testquery") != "test" {
				t.Fatal("Did not get default query on follow up request: ", req.URL.Query().Get("testquery"))
			}
		}
		firstRequest = false
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}

	client.Query().Add("testquery", "test")

	client.Get(&Params{
		Path:  "get",
		Query: url.Values{"testquery": []string{"test-override"}},
	})
	client.Get(&Params{Path: "get"})

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
	res, err := client.Post(&Params{
		Path:         "post",
		Body:         "hello",
		UnmarshalMap: UnmarshalMap{200: &success},
	})

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
	res, err := client.Post(&Params{
		Path: "post",
		Body: "hello",
		UnmarshalMap: UnmarshalMap{
			200:                 []interface{}{&success},
			http.StatusNotFound: &errResult,
		},
	})

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

	res, err := client.Post(&Params{Path: "post"})

	if err != nil {
		t.Fatal(err)
	}
	if res.Header.Get("X-Test-Response") != "test" {
		t.Fatal("Response mutator failed.")
	}
}

func TestJsonMarshaledBody(t *testing.T) {
	type test struct {
		Name string `json:"name"`
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

	_, err = client.Post(&Params{
		Path: "post",
		Body: body,
		UnmarshalMap: UnmarshalMap{
			200:                 &success,
			http.StatusNotFound: errResult,
		},
	})

	if err != nil {
		t.Fatal(err)
	}
	if success.Name != "test-result" {
		t.Fatal("Did not receive expected result.")
	}
}

func TestCloneClient(t *testing.T) {
	originalRequest := 0
	cloneRequest := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		originalRequest++
		if originalRequest > 1 {
			t.Fatal("Got more than one request at original server when we only expected one here.")
		}
		if req.Header.Get("X-Which") != "original" {
			t.Fatal("Got unexpected header from original: ", req.Header.Get("X-Which"))
		}
		if req.URL.Query().Get("query") != "original" {
			t.Fatal("Got unexpected query from original: ", req.URL.Query().Get("X-Which"))
		}
	}))
	defer server.Close()
	cloneServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cloneRequest++
		if cloneRequest > 1 {
			t.Fatal("Got more than one request at clone server when we only expected one here.")
		}
		if req.Header.Get("X-Which") != "clone" {
			t.Fatal("Got unexpected header from clone: ", req.Header.Get("X-Which"))
		}
		if req.URL.Query().Get("query") != "clone" {
			t.Fatal("Got unexpected query from clone: ", req.URL.Query().Get("X-Which"))
		}
	}))
	defer cloneServer.Close()

	base, err := url.Parse(server.URL)
	cloneBase, err := url.Parse(cloneServer.URL)
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}

	client.Headers().Set("X-Which", "original")
	client.Query().Set("query", "original")

	clone := client.Clone()
	clone.SetBaseUrl(cloneBase)
	clone.Headers().Set("X-Which", "clone")
	clone.Query().Set("query", "clone")

	client.Get(&Params{Path: ""})
	clone.Get(&Params{Path: ""})
}

func TestPathIsConcatenated(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/base/subpath" {
			t.Fatal("Proper subpathing is not occurring: ", req.URL.Path)
		}
	}))
	defer server.Close()

	base, err := url.Parse(server.URL + "/base")
	client, err := New(base)

	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Post(&Params{Path: "/subpath"})

	if err != nil {
		t.Fatal(err)
	}
}
