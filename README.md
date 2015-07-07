# grestclient

A REST client library written in go. Still under development. API can change at any time.

## Introduction

I started using [jmcvetta/napping][1] and was inspired to write my 
own go rest client based on that package. I haven't written much open
source code but this interested me enough to actually try and create
something someone else could use. I'll highlight some of
the differences between this package and the `napping` package in
a *Differences* section somewhere.

## Installation

`go get https://github.com/starJammer/grestclient`

## Description

This package is essentially a small/thin wrapper around http.Client. As
mentioned, I was using [jmcvetta/napping][1] but I wanted it to work slightly
differently. Instead of changing that package completely and doing a pull request,
I made my own. The main two differences are that the following:

1. `grestclient` lets you specify the marshaler to use to marshal the request body and 
the unmarshaler to use on the response body.
2. `grestclient` lets you specify request/response mutators to execute before the request
goes out and after the response comes back.

Besides that there are some small things that I made because I prefer it
this way.

1. When the client is created it stores a base url. So you don't have to
keep specifying it again on each request. Instead, each request only needs
the relative path to be specified.
2. It uses http.Request/http.Response directly instead of creating new types
to wrap them. Not better or worse than `napping` but I went with this out
of personal feeling.
3. No direct way of just doing a request yet. So there is no `session.Send`
equivalent. I've thought of putting it in for letting others have more
control over the request but I haven't had needed this use case. The 
RequestMutator can help you alter a request about to be executed if need.

## Usage

```go
package main

import ( 
    gr "github.com/starJammer/grestclient"
    "net/url"
    "net/http"
    "fmt"
)

func main(){
    u, err := url.Parse( "http://example.com" )
    if err != nil {
        fmt.Println( "Bad url: ", err )
    }

    //you'll get an error if you pass in nil, every client needs
    //a base url.
    c, err := gr.New( u )
    if err != nil {
        fmt.Println( "Bad url: ", err )
    }

    //By default bodies are marshaled/unmarshaled as text using
    //the StringMarshalerFunc/StringUnmarshalerFunc pair but you can
    //use the JsonMarshalerFunc/JsonUnmarshalerFunc pair too. Or you can
    //mix and match. Send text and expect json. Create your own
    //MarshalerFuncs/UnmarshalerFuncs.

    //these are the defaults so no need to do this really
    c.SetMarshaler( gr.StringUnmarshalerFunc )
    c.SetUnmarshaler( gr.StringUnmarshalerFunc )

    //Let's do json
    c.SetMarshaler( gr.JsonMarshalerFunc )
    c.SetUnmarshaler( gr.JsonUnmarshalerFunc )
    //makes sure request's Content-Type header has application/json in it.
    c.AddRequestMutators( gr.JsonContentTypeMutator )

    //I included a convenience function for the above 3 calls so you can do
    gr.SetupForJson( c ) //instead of calling them individually, unless you need
    //mixed (un)marshaling or something.

    var successResult, errorResult string
    c.Get( "path/to/resource", 
           url.Values{}, //you can also pass nil here. these are used in the query portion of the url
           &successResult, //you can pass nil here if you want. 
                           //This is where any responses with a code less than 300 get unmarshaled to
           &errorResult ) //You can pass nil if you want.
                          //This is where any responses with a code greater than 400 get unmarshaled to

    var postBody string = "latino"

    c.Post( "path/to/resource", 
           url.Values{}, //you can also pass nil here. these are used in the query portion of the url
           &postBody, //the post body can be nil if you want. You can pass by reference or by value.
                      //Using pointers or passing by reference is usually preferenced unless it's simple type
           &successResult, //you can pass nil here if you want. 
                           //This is where any responses with a code less than 300 get unmarshaled to
           &errorResult ) //You can pass nil if you want.
                          //This is where any responses with a code greater than 400 get unmarshaled to

    c.Post( "", //use blank url to perform an operation on just the base url itself
           url.Values{}, //you can also pass nil here. these are used in the query portion of the url
           &postBody, //the post body can be nil if you want. You can pass by reference or by value.
                      //Using pointers or passing by reference is usually preferenced unless it's simple type
           &successResult, //you can pass nil here if you want. 
                           //This is where any responses with a code less than 300 get unmarshaled to
           &errorResult ) //You can pass nil if you want.
                          //This is where any responses with a code greater than 400 get unmarshaled to


    //That's it for the most part. You can use the other methods too:

    //c.Put - similar to post but puts instead of posts

    //c.Patch - similar to post but patches instead of posts

    //c.Head - similar to get but there is no query parameter
    //on the method even though you can set default queries on the client itself.
    //Does HEAD need query params? I can add them

    //c.Options - similar to HEAD but does options method instead

    //c.Delete - similar to get. Lets you use query values but no body is included.

    //Default headers and queries can be set on the client so that they are
    //issued with EVERY request.

    c.Headers().Add( "X-Test", "test" )
    c.Query().Add( "sync", "sync" )

    //You can override a default query if you call a method that lets you 
    //specify a query.
    c.Get( "/path", url.Values{ "sync": []string{ "not-sync" } }
    //The sync query parameter will now be not-sync instead of sync.

    //No way to override headers for specific methods yet unless you want
    //to use a RequestMutator to override headers on some requests. 

    c.AddRequestMutators( func( r *http.Request ) error {
        if r.Method == "POST" {
            r.Header.Set( "X-Test", "not-test" )
        }
    })

    //However, you could also just clone the client and set different headers
    //on the clone
    clonedClient := c.Clone()
    clonedClient.Headers().Set( "X-Test", "not-test" )

}
```

## Contributing

Fork it and make a pull request. I don't plan on doing any versioning
for the package really other than the branches. I hope that my design 
is good enough that the interface won't change much if at all.

[1]: https://github.com/jmcvetta/napping
