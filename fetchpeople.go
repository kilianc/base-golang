package main

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

type FetchPeople struct {
	FastHTTPDo func(req *fasthttp.Request, resp *fasthttp.Response) error
}

type FetchPeopleOptions struct {
	FastHTTPDo func(req *fasthttp.Request, resp *fasthttp.Response) error
}

func (f *FetchPeople) Fetch(url string, opts ...FetchPeopleOptions) ([]Person, error) {
	// default to fasthttp.Do if not set, this is useful for testing
	if (f.FastHTTPDo) == nil {
		f.FastHTTPDo = fasthttp.Do
	}

	if len(opts) > 0 && opts[0].FastHTTPDo != nil {
		f.FastHTTPDo = opts[0].FastHTTPDo
	}

	if url == "" {
		return nil, fmt.Errorf("url is empty")
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod("GET")

	if err := f.FastHTTPDo(req, resp); err != nil {
		return nil, err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var people []Person
	if err := json.Unmarshal(resp.Body(), &people); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %s", err)
	}

	return people, nil
}

type AsyncFetchPeopleResult struct {
	People []Person
	Err    error
}

func (ar AsyncFetchPeopleResult) Unwrap() ([]Person, error) {
	return ar.People, ar.Err
}

func AsyncFetchPeople(url string) <-chan AsyncFetchPeopleResult {
	done := make(chan AsyncFetchPeopleResult)

	go func() {
		data, err := (&FetchPeople{}).Fetch(url)
		done <- AsyncFetchPeopleResult{data, err}
	}()

	return done
}
