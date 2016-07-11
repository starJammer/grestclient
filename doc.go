/*
Example use:

	import (
		"net/url"
	)

	base, err := url.Parse("http://example.com")
	c, err := New(base)

	c.Get(&Params{
		Path: "path/to/resource",
		Headers:
		Query:
		Body:
		unmarshalerMap:
	})

*/
package grestclient
