package main

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

type FetchPeopleOptions struct {
	FastHTTPDo func(req *fasthttp.Request, resp *fasthttp.Response) error
}

func FetchPeople(url string, opts ...FetchPeopleOptions) ([]Person, error) {
	if url == "" {
		return nil, fmt.Errorf("url is empty")
	}

	fastHTTPDo := fasthttp.Do

	if len(opts) > 0 && opts[0].FastHTTPDo != nil {
		fastHTTPDo = opts[0].FastHTTPDo
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod("GET")

	if err := fastHTTPDo(req, resp); err != nil {
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

type FetchPeopleResultAsync struct {
	People []Person
	Err    error
}

func (ar FetchPeopleResultAsync) Unwrap() ([]Person, error) {
	return ar.People, ar.Err
}

func FetchPeopleAsync(url string, options ...FetchPeopleOptions) <-chan FetchPeopleResultAsync {
	done := make(chan FetchPeopleResultAsync)

	go func() {
		data, err := FetchPeople(url, options...)
		done <- FetchPeopleResultAsync{data, err}
	}()

	return done
}
