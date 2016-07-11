/*
Example use:

	import (
		"net/url"
	)

	base, err := url.Parse("http://example.com")
	c, err := New(base)

	var success string
	var otherSuccess string
	var fail string
	c.Get(&Params{
		Path: "path/to/resource",
		Headers: http.Header{"One-Off": []string{"one"}},
		Query: url.Values{"testquery": []string{"test-override"}},
		UnmarshalMap: UnmarshalMap{
			200: &success,
			201: &otherSuccess,
			202: &otherSuccess,
			400: &fail,
			404: &fail,
		},
	})

	c.Post(&Params{
		Path: "path/to/resource",
		Headers: http.Header{"One-Off": []string{"one"}},
		Query: url.Values{"testquery": []string{"test-override"}},
		Body: "hello",
		UnmarshalMap: UnmarshalMap{
			200: &success,
			201: &otherSuccess,
			202: &otherSuccess,
			400: &fail,
			404: &fail,
		},
	})

*/
package grestclient
